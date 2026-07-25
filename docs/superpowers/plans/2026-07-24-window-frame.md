# Faux OS Window Frame Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Constrain the web app inside a centered, terminal-style faux OS window (title bar with route path + clock) on a textured desktop background.

**Architecture:** New `WindowFrame` component owns the desktop centering, frame chrome, title bar, and clock. `Layout` wraps the existing `Sidebar` + `<main>` in it. Theme tokens shift to deeper green + ochre accent (orange removed).

**Tech Stack:** React 19, react-router-dom 7, Tailwind CSS 4 (`@theme` tokens), Vitest + Testing Library.

## Global Constraints

- Spec: `docs/superpowers/specs/2026-07-24-window-frame-design.md`
- All work under `web/`. Run commands from `web/` with `direnv exec . <cmd>` (Node comes from devenv, not PATH).
- Test command: `cd web && direnv exec . npx vitest run <file>` (full suite: `direnv exec . npm test`).
- Colors (exact): green `#2bd94a`, ochre `#d4a017`, desktop `#050505`, window body `#0a0a0a` (existing `--color-hacker-bg`), chrome `#000000` (existing `--color-hacker-chrome`).
- No decorative window controls. No new dependencies.
- Never push to main; current branch `icullinane/version2` is the default branch.

---

### Task 1: Theme tokens + desktop background + ochre accent swap

**Files:**
- Modify: `web/src/theme.css`
- Modify: `web/src/components/Sidebar.tsx:19-20`
- Modify: `web/src/pages/HistoryPage.tsx:50`
- Test: existing suites must stay green (no color assertions exist)

**Interfaces:**
- Consumes: nothing.
- Produces: Tailwind utilities `hacker-ochre`, `hacker-desktop` (from `--color-hacker-ochre`, `--color-hacker-desktop`); `hacker-orange` no longer exists — later tasks must use `hacker-ochre`.

- [ ] **Step 1: Update theme tokens and desktop background**

Replace the `@theme` block and `body` rule in `web/src/theme.css` with:

```css
@theme {
  --font-army: "US Army", ui-monospace, SFMono-Regular, monospace;
  --color-hacker-bg: #0a0a0a;
  --color-hacker-desktop: #050505;
  --color-hacker-chrome: #000000;
  --color-hacker-green: #2bd94a;
  --color-hacker-ochre: #d4a017;
  --color-hacker-fg-muted: #7a9a85;
}

body {
  background-color: var(--color-hacker-desktop);
  color: var(--color-hacker-green);
  /* faint CRT scanlines + soft center vignette, CSS-only */
  background-image:
    radial-gradient(ellipse at center, rgb(43 217 74 / 0.06), transparent 70%),
    repeating-linear-gradient(
      0deg,
      rgb(43 217 74 / 0.03) 0px,
      rgb(43 217 74 / 0.03) 1px,
      transparent 1px,
      transparent 3px
    );
}
```

(This removes `--color-hacker-orange` entirely.)

- [ ] **Step 2: Swap orange → ochre at both usage sites**

`web/src/components/Sidebar.tsx` — in the NavLink className, change:

```tsx
              isActive
                ? "border-l-2 border-hacker-ochre text-hacker-ochre"
                : "text-hacker-green hover:text-hacker-ochre"
```

`web/src/pages/HistoryPage.tsx:50` — change:

```tsx
          <tr className="text-sm uppercase text-hacker-ochre">
```

- [ ] **Step 3: Verify no orange remains and suite is green**

Run: `grep -rn "hacker-orange" web/src` — Expected: no matches.
Run: `cd web && direnv exec . npm test` — Expected: all tests PASS.

- [ ] **Step 4: Commit**

```bash
git add web/src/theme.css web/src/components/Sidebar.tsx web/src/pages/HistoryPage.tsx
git commit -m "feat(web): deeper green, ochre accent, desktop background tokens"
```

---

### Task 2: WindowFrame component (TDD)

**Files:**
- Create: `web/src/components/WindowFrame.tsx`
- Create: `web/src/components/WindowFrame.test.tsx`

**Interfaces:**
- Consumes: Tailwind utilities `hacker-ochre`, `hacker-bg`, `hacker-chrome` from Task 1; `useLocation` from react-router-dom (must render inside a Router).
- Produces: `export default function WindowFrame({ children }: { children: ReactNode })` — Task 3 wraps Layout content in it. Children are rendered inside a `flex min-h-0 flex-1 flex-col md:flex-row` container.

- [ ] **Step 1: Write the failing tests**

Create `web/src/components/WindowFrame.test.tsx`:

```tsx
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, act } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import WindowFrame from "./WindowFrame";

function renderFrame(path = "/players") {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <WindowFrame>
        <p>inner content</p>
      </WindowFrame>
    </MemoryRouter>,
  );
}

describe("WindowFrame", () => {
  beforeEach(() => {
    vi.useFakeTimers({ now: new Date(2026, 6, 24, 10, 15, 30) });
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("renders its children", () => {
    renderFrame();
    expect(screen.getByText("inner content")).toBeInTheDocument();
  });

  it("shows the title with the current route as a path", () => {
    renderFrame("/history");
    expect(screen.getByText("PRISONER — ~/history")).toBeInTheDocument();
  });

  it("shows a clock that ticks every second", () => {
    renderFrame();
    expect(screen.getByText("10:15:30")).toBeInTheDocument();

    act(() => {
      vi.advanceTimersByTime(1000);
    });
    expect(screen.getByText("10:15:31")).toBeInTheDocument();
  });

  it("clears its interval on unmount", () => {
    const { unmount } = renderFrame();
    unmount();
    act(() => {
      vi.advanceTimersByTime(1000);
    });
    // nothing to assert beyond "no act/setState warnings"; getTimerCount proves cleanup
    expect(vi.getTimerCount()).toBe(0);
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd web && direnv exec . npx vitest run src/components/WindowFrame.test.tsx`
Expected: FAIL — cannot resolve `./WindowFrame`.

- [ ] **Step 3: Implement WindowFrame**

Create `web/src/components/WindowFrame.tsx`:

```tsx
import { useEffect, useState, type ReactNode } from "react";
import { useLocation } from "react-router-dom";

function formatClock(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

export default function WindowFrame({ children }: { children: ReactNode }) {
  const { pathname } = useLocation();
  const [now, setNow] = useState(() => new Date());

  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  return (
    <div className="fixed inset-0 flex items-center justify-center md:p-[max(1rem,3vmin)]">
      <div className="flex h-full w-full max-h-[860px] max-w-[1280px] flex-col border border-hacker-ochre bg-hacker-bg shadow-[0_0_24px_rgb(43_217_74_/_0.15)]">
        <header className="flex items-center justify-between border-b border-hacker-ochre bg-hacker-chrome px-3 py-1.5 font-army text-sm uppercase tracking-wide text-hacker-ochre">
          <span>PRISONER — ~{pathname}</span>
          <time>{formatClock(now)}</time>
        </header>
        <div className="flex min-h-0 flex-1 flex-col md:flex-row">
          {children}
        </div>
      </div>
    </div>
  );
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd web && direnv exec . npx vitest run src/components/WindowFrame.test.tsx`
Expected: 4 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/components/WindowFrame.tsx web/src/components/WindowFrame.test.tsx
git commit -m "feat(web): WindowFrame terminal-window chrome with route title and clock"
```

---

### Task 3: Wrap Layout in WindowFrame

**Files:**
- Modify: `web/src/components/Layout.tsx`
- Modify: `web/src/components/Sidebar.tsx:11`
- Test: `web/src/App.test.tsx` (existing) + full suite

**Interfaces:**
- Consumes: `WindowFrame` from Task 2 (children get a `flex-col md:flex-row` container, so Layout no longer needs its own flex wrapper).
- Produces: final composed app; no new exports.

- [ ] **Step 1: Rewrite Layout to use WindowFrame**

Replace `web/src/components/Layout.tsx` with:

```tsx
import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import WindowFrame from "./WindowFrame";

export default function Layout() {
  return (
    <WindowFrame>
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-6">
        <Outlet />
      </main>
    </WindowFrame>
  );
}
```

- [ ] **Step 2: Let the sidebar fill the frame, not the viewport**

`web/src/components/Sidebar.tsx:11` — remove `md:min-h-screen` (the frame's flex row now sets the height):

```tsx
    <nav className="flex gap-4 border-b border-hacker-green/30 bg-hacker-chrome p-4 md:w-48 md:flex-col md:border-b-0 md:border-r">
```

- [ ] **Step 3: Run the full suite**

Run: `cd web && direnv exec . npm test`
Expected: all tests PASS (App/Layout tests exercise the new frame; title text appears once per page render).

- [ ] **Step 4: Visual sanity check**

Run: `cd web && direnv exec . npm run build`
Expected: build succeeds (catches Tailwind class typos and TS errors).

- [ ] **Step 5: Commit**

```bash
git add web/src/components/Layout.tsx web/src/components/Sidebar.tsx
git commit -m "feat(web): render app inside centered WindowFrame"
```
