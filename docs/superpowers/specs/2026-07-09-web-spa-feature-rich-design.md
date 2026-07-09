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
  green (`#39ff14`-family) as the single accent for links, buttons, borders,
  and focus rings; muted gray-green (`#7a9a85`-ish) for secondary/disabled
  text. No second accent color.

## Components & data flow

- `PlayersPage` — thin wrapper rendering the existing `PlayerManager`
  unchanged (list + create-player form), restyled with Tailwind classes.
- `HistoryPage` (new) — fetches `GET /history` on mount via
  `api/history.ts::listHistory()`. Renders a table of interactions (player
  IDs, moves, played-at). Optional player filter: a `<select>` populated from
  `GET /players`; choosing a player refetches via
  `listHistoryForPlayer(id)` → `GET /history/{id}`.
- `PlayPage` (new) — fetches `GET /players` on mount for two `<select>`
  dropdowns (Player A / Player B). Move selection is two button pairs
  (Cooperate / Betray) per side, local state only. Submit calls
  `api/play.ts::playRound(req)` → `POST /play`; on success shows the returned
  scores inline; disables the submit button while the request is in flight.
- Each page fetches independently on mount — **no shared player context**.
  Simpler code; the roster is small enough that redundant fetches are cheap.
- New frontend types in `types.ts` mirroring the Go JSON wire shapes:
  - `Interaction` (camelCase: `id`, `playerA`, `playerB`, `playerAMove`,
    `playerBMove`, `playedAt`) matching `internal/types.Interaction`.
  - `PlayRequest`/`PlayResponse` (snake_case: `player_a`, `player_b`,
    `player_a_move`, `player_b_move`, `id`, `player_a_score`,
    `player_b_score`) matching `api.PlayRequest`/`api.PlayResponse`.
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

## Out of scope

- Shared/global player state across pages.
- A second accent color beyond phosphor green.
- Editing/deleting players or history entries.
- Real-time updates (polling/websockets) for history.
