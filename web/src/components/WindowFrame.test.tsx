import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, act } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import WindowFrame from "./WindowFrame";

vi.mock("../api/players", () => ({
  listPlayers: vi.fn().mockResolvedValue([
    { id: "p1", name: "ian", CreatedAt: "2026-01-01" },
    { id: "p2", name: "steve", CreatedAt: "2026-01-01" },
  ]),
}));

vi.mock("../api/history", () => ({
  listHistory: vi.fn().mockResolvedValue([
    {
      playerAName: "ian",
      playerBName: "steve",
      playerAMove: "C",
      playerBMove: "B",
      playedAt: "2026-08-20T12:00:00Z",
    },
    {
      playerAName: "steve",
      playerBName: "ian",
      playerAMove: "C",
      playerBMove: "C",
      playedAt: "2026-08-21T12:00:00Z",
    },
    {
      playerAName: "ian",
      playerBName: "steve",
      playerAMove: "B",
      playerBMove: "B",
      playedAt: "2026-08-22T12:00:00Z",
    },
  ]),
}));

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
  describe("with fake timers", () => {
    beforeEach(() => {
      vi.useFakeTimers({ now: new Date(2026, 6, 24, 10, 15, 30) });
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    async function renderAndSettle(path?: string) {
      const view = renderFrame(path);
      // flush the mocked api promises so effects settle inside act
      await act(async () => {});
      return view;
    }

    it("renders its children", async () => {
      await renderAndSettle();
      expect(screen.getByText("inner content")).toBeInTheDocument();
    });

    it("shows the title with the current route as a path", async () => {
      await renderAndSettle("/history");
      expect(screen.getByText("prisoner — ~/history")).toBeInTheDocument();
    });

    it("shows a clock that ticks every second", async () => {
      await renderAndSettle();
      expect(screen.getByText("10:15:30")).toBeInTheDocument();

      act(() => {
        vi.advanceTimersByTime(1000);
      });
      expect(screen.getByText("10:15:31")).toBeInTheDocument();
    });

    it("clears its interval on unmount", async () => {
      const { unmount } = await renderAndSettle();
      unmount();
      act(() => {
        vi.advanceTimersByTime(1000);
      });
      expect(vi.getTimerCount()).toBe(0);
    });
  });

  it("shows total rounds and players in the titlebar", async () => {
    // render on /history so the path text doesn't also match /players/
    renderFrame("/history");
    expect(await screen.findByText(/rounds/)).toHaveTextContent("rounds 3");
    expect(screen.getByText(/players/)).toHaveTextContent("players 2");
  });
});
