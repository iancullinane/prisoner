# Faux OS Window Frame — Design

2026-07-24

## Goal

Constrain the web app inside a centered, terminal-style faux OS window that scales
with the viewport, sitting on a subtly textured "desktop" background. Same left
nav, same pages — just framed.

## Palette changes (`web/src/theme.css`)

- `--color-hacker-green`: `#39ff14` → `#2bd94a` (deeper, less neon)
- `--color-hacker-ochre`: `#d4a017` (new; primary accent)
- `--color-hacker-desktop`: `#050505` (new; page background behind the window)
- `--color-hacker-orange` is removed; ochre takes over its accent role
  (Sidebar active/hover states switch from orange to ochre). Two accents max:
  green + ochre.
- Window body stays `--color-hacker-bg` (`#0a0a0a`); chrome stays black.

## Desktop background

`body` gets `--color-hacker-desktop` plus a CSS-only faint scanline/grid pattern
in dark green and a soft radial vignette. No images, no fake taskbar/props.

## WindowFrame component (`web/src/components/WindowFrame.tsx`)

New component wrapping the app content:

- Centered in the viewport with a `~3vmin` gutter (min ~1rem) on all sides.
- Fluid size, capped at `max-w-[1280px] max-h-[860px]`; fills available space
  below the cap.
- 1px ochre border, subtle green glow shadow, sharp corners (no border-radius).
- **Title bar**: black chrome, thin ochre bottom border. Left: title
  `PRISONER — ~/<route>` where the path segment comes from `useLocation`
  (e.g. `~/players`). Right: live clock `HH:MM:SS`, updated on a 1s interval,
  cleaned up on unmount. No decorative window controls.
- Children render below the title bar and fill the remaining frame height.

Props: `children` only. Route and clock are derived internally — keeps the
call site in `Layout` trivial.

## Layout changes (`web/src/components/Layout.tsx`)

- Renders `<WindowFrame>` around the existing `Sidebar` + `<main>`.
- Sidebar's `md:min-h-screen` becomes fill-frame-height (the frame, not the
  viewport, is now the height reference). Content area scrolls inside the
  frame (`overflow-y-auto` on `main`) rather than growing the page.

## Mobile (`< md`)

Gutters collapse to 0 (full-bleed window), frame border and title bar kept,
desktop texture effectively hidden. Sidebar keeps its existing horizontal
top-bar collapse.

## Testing (TDD)

- `WindowFrame.test.tsx` (new):
  - renders children
  - title shows `~/players` when route is `/players` (MemoryRouter)
  - clock renders and advances with fake timers; interval cleared on unmount
- `Sidebar.test.tsx` / `Layout` snapshots updated for class/token changes only.

## Out of scope

Draggable/resizable window, real window controls, multiple windows, fake
taskbar or desktop icons.
