import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import HistoryPage from "./HistoryPage";
import * as historyApi from "../api/history";
import * as playersApi from "../api/players";
import type { Move } from "../types";

const sevenHoursAgo = new Date(Date.now() - 7 * 60 * 60 * 1000).toISOString();

function interaction(
  playerAName: string,
  playerBName: string,
  playerAMove: Move,
  playerBMove: Move,
  playedAt = sevenHoursAgo,
) {
  return { playerAName, playerBName, playerAMove, playerBMove, playedAt };
}

// betrayal (0·5), mutual trust (3·3), mutual ruin (1·1)
const fixtures = [
  interaction("Alice", "Bob", "C", "B"),
  interaction("Alice", "Cara", "C", "C"),
  interaction("Bob", "Cara", "B", "B"),
];

describe("HistoryPage", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);
    vi.spyOn(historyApi, "listHistory").mockResolvedValue(fixtures);
  });

  it("renders rows with move chips, outcomes, and points", async () => {
    render(<HistoryPage />);

    // two table rows plus the "most trusting" summary cell
    expect(await screen.findAllByText("Alice")).toHaveLength(3);
    expect(screen.getAllByText("✓ C")).toHaveLength(3); // C-B, C-C
    expect(screen.getAllByText("✕ B")).toHaveLength(3); // C-B, B-B
    expect(screen.getByText("Betrayal")).toBeInTheDocument();
    expect(screen.getByText("Mutual trust")).toBeInTheDocument();
    expect(screen.getByText("Mutual ruin")).toBeInTheDocument();

    // Pts is the 6th column; the winning score is wrapped in <b>
    const betrayalRow = screen.getByText("Betrayal").closest("tr")!;
    expect(within(betrayalRow).getAllByRole("cell")[5]).toHaveTextContent(
      "0 · 5",
    );
    const trustRow = screen.getByText("Mutual trust").closest("tr")!;
    expect(within(trustRow).getAllByRole("cell")[5]).toHaveTextContent("3 · 3");
  });

  it("summarizes rounds, betrayal rate, and extremes in the strip", async () => {
    render(<HistoryPage />);

    const strip = await screen.findByRole("region", { name: /summary/i });
    expect(within(strip).getByText("Rounds")).toBeInTheDocument();
    expect(within(strip).getByText("3")).toBeInTheDocument();
    // 3 betrayals across 6 moves
    expect(within(strip).getByText("50%")).toBeInTheDocument();
    // Alice cooperates every round; Bob betrays twice
    expect(within(strip).getByText("Alice")).toBeInTheDocument();
    expect(within(strip).getByText("100% coop")).toBeInTheDocument();
    expect(within(strip).getByText("Bob")).toBeInTheDocument();
    expect(within(strip).getByText("2 betrayals")).toBeInTheDocument();
  });

  it("filters rows by outcome with the segmented bar", async () => {
    const user = userEvent.setup();
    render(<HistoryPage />);
    await screen.findByText("Mutual trust");

    const segbar = screen.getByRole("group", { name: /outcome/i });
    await user.click(within(segbar).getByRole("button", { name: "betrayals" }));

    expect(
      within(segbar).getByRole("button", { name: "betrayals" }),
    ).toHaveAttribute("aria-pressed", "true");
    expect(within(segbar).getByRole("button", { name: "all" })).toHaveAttribute(
      "aria-pressed",
      "false",
    );
    expect(screen.getByText("Betrayal")).toBeInTheDocument();
    expect(screen.queryByText("Mutual trust")).not.toBeInTheDocument();
    expect(screen.queryByText("Mutual ruin")).not.toBeInTheDocument();
  });

  it("shows an empty state when no rows match the filter", async () => {
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([
      interaction("Alice", "Cara", "C", "C"),
    ]);
    const user = userEvent.setup();
    render(<HistoryPage />);
    await screen.findByText("Mutual trust");

    const segbar = screen.getByRole("group", { name: /outcome/i });
    await user.click(within(segbar).getByRole("button", { name: "betrayals" }));

    expect(screen.getByText(/no rounds match/i)).toBeInTheDocument();
    expect(screen.getByText(/press \[3\] to play one/i)).toBeInTheDocument();
  });

  it("renders played-at as a relative time", async () => {
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([fixtures[0]]);

    render(<HistoryPage />);

    expect(await screen.findByText("7 hours ago")).toBeInTheDocument();
  });

  it("shows an error message when listHistory fails", async () => {
    vi.spyOn(historyApi, "listHistory").mockRejectedValue(
      new Error("network down"),
    );

    render(<HistoryPage />);

    expect(await screen.findByText(/network down/i)).toBeInTheDocument();
  });

  it("refetches via listHistoryForPlayer when a player filter is selected", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      { id: "p1", name: "Chris", CreatedAt: "2026-01-01T00:00:00Z" },
    ]);
    vi.spyOn(historyApi, "listHistoryForPlayer").mockResolvedValue(fixtures);

    const user = userEvent.setup();
    render(<HistoryPage />);

    await waitFor(() => expect(historyApi.listHistory).toHaveBeenCalledTimes(1));

    await user.selectOptions(
      screen.getByLabelText(/filter by player/i),
      screen.getByRole("option", { name: "Chris" }),
    );

    await waitFor(() =>
      expect(historyApi.listHistoryForPlayer).toHaveBeenCalledWith("p1"),
    );
  });
});
