# Web SPA: Players, History, Play — Design

Date: 2026-07-09
Status: Approved

## Goal

Turn `web/` from a single Players screen into a small multi-page SPA covering
Players, History, and Play, with a cyberpunk-terminal visual identity (black
background, neon-green accent, US Army display font) and a responsive
left-sidebar navigation.

## Architecture

- **Routing**: add `react-router-dom`. Routes: `/players` (default redirect
  from `/`), `/history`, `/play`.
- **API path collision with the deployed Ingress**: `sheeta-games/k8s/backend-ingress.yaml`
  currently routes the exact path prefixes `/players`, `/play`, and `/history`
  straight to the Go backend Service, bypassing the SPA. That collides
  head-on with the SPA routes above — a hard refresh or deep link to
  `/history` would hit the backend and get raw JSON instead of the React
  app, and `/play` has no `GET` handler at all (only `POST`), so it would
  404. This is resolved by versioning the API (next bullet) rather than by
  picking different frontend route names, so the API gets a stable,
  future-proof namespace instead of just dodging the immediate collision.
- **API versioning**: all business routes in `api/server.go` move under
  `/api/v1` (`/api/v1/players`, `/api/v1/players/{name}`,
  `/api/v1/history`, `/api/v1/history/{id}`, `/api/v1/play`). `/healthz`
  stays unprefixed — ALB target-group health checks hit the pod port
  directly via the `alb.ingress.kubernetes.io/healthcheck-path` annotation,
  never through an Ingress path rule, so it doesn't need to move. Frontend
  API modules (`api/players.ts`, `api/history.ts`, `api/play.ts`) prefix
  their fetch paths with `/api/v1`; `VITE_API_BASE_URL` semantics are
  unchanged (same-origin `""` in prod, `http://localhost:5001` in dev).
- **Infra follow-up (separate repo, `sheeta-games`)**: `backend-ingress.yaml`'s
  three path rules (`/players`, `/play`, `/history`) collapse into a single
  `path: /api`, `pathType: Prefix` rule → `prisoner-backend`. Ingress
  `Prefix` matching is segment-based, so `/api` also matches
  `/api/v1/players` etc. — this is a simplification, and any future
  endpoint added under `/api/v1/*` routes correctly with no further Ingress
  edits. The `/` catch-all → frontend is unchanged, and `/players`,
  `/history`, `/play` now fall through to the SPA. **Deploy sequencing**:
  the versioned backend/frontend images and the Ingress change must land
  together — deploying one without the other 404s API calls for the window
  in between. This Ingress change lives in `sheeta-games` and ships as its
  own PR in that repo, coordinated with the `prisoner` release.
- **Styling**: add Tailwind CSS (v4, `@tailwindcss/vite` plugin) — replaces
  ad-hoc CSS. New `tailwind.config` content globs over `src/**/*.tsx`.
- **Layout**: a `Layout` component renders a `Sidebar` (fixed left nav on
  desktop, collapses to a top bar below the `md` breakpoint) plus a `<main>`
  area rendering `<Outlet />`. `App.tsx` becomes `BrowserRouter` + `Layout` +
  `Routes`.
- **Fonts**: US Army Bold/Regular/Light `.ttf` files live in
  `web/src/assets/fonts/`, loaded via `@font-face` in a new global stylesheet
  (`src/theme.css`), registered as `font-army` in the Tailwind theme and used
  for headings/nav; body text falls back to a system/monospace stack for
  readability at small sizes.
- **Palette**: black backgrounds (`#0a0a0a` page, `#000` chrome), phosphor
  green (`#39ff14`-family) as the primary accent for links, buttons, borders,
  and focus rings; burnt orange (`#ff6a1a`-family) as a secondary accent for
  warnings, destructive actions (e.g. a future delete), and highlighting the
  active nav item; muted gray-green (`#7a9a85`-ish) for secondary/disabled
  text.

## Components & data flow

- `PlayersPage` — thin wrapper rendering the existing `PlayerManager`
  unchanged (list + create-player form), restyled with Tailwind classes.
- `HistoryPage` (new) — fetches `GET /api/v1/history` on mount via
  `api/history.ts::listHistory()`. Renders a table of interactions (player
  IDs, moves, played-at). Optional player filter: a `<select>` populated from
  `GET /api/v1/players`; choosing a player refetches via
  `listHistoryForPlayer(id)` → `GET /api/v1/history/{id}`.
- `PlayPage` (new) — fetches `GET /api/v1/players` on mount for two
  `<select>` dropdowns (Player A / Player B). Move selection is two button
  pairs (Cooperate / Betray) per side, local state only. Submit calls
  `api/play.ts::playRound(req)` → `POST /api/v1/play`; on success shows the
  returned scores inline; disables the submit button while the request is
  in flight.
- Each page fetches independently on mount — **no shared player context**.
  Simpler code; the roster is small enough that redundant fetches are cheap.
- New frontend types in `types.ts` mirroring the Go JSON wire shapes:
  - `Interaction` (camelCase: `id`, `playerA`, `playerB`, `playerAMove`,
    `playerBMove`, `playedAt`) matching `internal/types.Interaction`.
  - `PlayRequest`/`PlayResponse` (camelCase: `playerA`, `playerB`,
    `playerAMove`, `playerBMove`, `id`, `playerAScore`, `playerBScore`)
    matching `api.PlayRequest`/`api.PlayResponse`. These were previously
    snake_case on the wire; the Go structs were updated to camelCase to match
    `Interaction`'s convention (see `api/types.go`, `api/types_test.go`).
  - `Move` as a `"C" | "B"` string union.
  - `Result` as a `"T" | "R" | "P" | "S"` string union (Temptation, Reward,
    Punish, Sucker — matches `pkg/prisoner.Result`'s `MarshalJSON`).
- New API modules follow the existing `api/players.ts` pattern (fetch,
  throw `Error` with status on `!response.ok`, return parsed JSON):
  - `api/history.ts`: `listHistory()`, `listHistoryForPlayer(id: string)`.
  - `api/play.ts`: `playRound(req: PlayRequest): Promise<PlayResponse>`.

## Error handling

- Each page owns local `error: string | null` state, same pattern as
  `PlayerManager` today, rendered as `<p role="alert">`.
- No global error boundary — out of scope for this pass.
- `PlayPage` guards against double-submits by disabling its submit button for
  the duration of the in-flight request and re-enabling on both success and
  failure.
- Player-cannot-play-itself is already enforced server-side (400); the UI
  additionally disables submit client-side when Player A === Player B to
  avoid a round-trip for the obvious case.

## Testing

Test-first for all new code, matching existing `vitest` +
`@testing-library/react` setup (`players.test.ts`, `PlayerManager.test.tsx`
as reference patterns):

- `api/history.test.ts` — `listHistory`/`listHistoryForPlayer` success and
  non-ok-response cases.
- `api/play.test.ts` — `playRound` success and non-ok-response cases.
- `components/Sidebar.test.tsx` — renders nav links, marks the active route.
- `components/HistoryPage.test.tsx` — loads and renders interactions;
  surfaces fetch errors.
- `components/PlayPage.test.tsx` — loads players into dropdowns, submits a
  play, renders the score, disables submit while pending, surfaces errors.
- `api/server_test.go` (existing, backend) — update request paths from
  `/players`, `/history`, `/play` to `/api/v1/players`, `/api/v1/history`,
  `/api/v1/play`; `/healthz` test is unchanged.

## Out of scope

- Shared/global player state across pages.
- Editing/deleting players or history entries.
- Real-time updates (polling/websockets) for history.
