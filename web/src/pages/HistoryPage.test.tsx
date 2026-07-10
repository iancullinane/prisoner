import { describe, it, expect, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import HistoryPage from "./HistoryPage";
import * as historyApi from "../api/history";
import * as playersApi from "../api/players";

const interaction = {
  id: "i1",
  playerA: "p1",
  playerB: "p2",
  playerAMove: "C" as const,
  playerBMove: "B" as const,
  playedAt: "2026-01-01T00:00:00Z",
};

describe("HistoryPage", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("renders history rows returned by listHistory", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([interaction]);

    render(<HistoryPage />);

    expect(await screen.findByText("p1")).toBeInTheDocument();
    expect(screen.getByText("p2")).toBeInTheDocument();
  });

  it("shows an error message when listHistory fails", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);
    vi.spyOn(historyApi, "listHistory").mockRejectedValue(new Error("network down"));

    render(<HistoryPage />);

    expect(await screen.findByText(/network down/i)).toBeInTheDocument();
  });

  it("refetches via listHistoryForPlayer when a player filter is selected", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      { id: "p1", name: "Chris", CreatedAt: "2026-01-01T00:00:00Z" },
    ]);
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([interaction]);
    vi.spyOn(historyApi, "listHistoryForPlayer").mockResolvedValue([interaction]);

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
