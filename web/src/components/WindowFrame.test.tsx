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
