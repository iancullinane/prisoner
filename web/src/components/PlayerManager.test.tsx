import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import PlayerManager from "./PlayerManager";
import * as playersApi from "../api/players";
import * as historyApi from "../api/history";
import type { Move } from "../types";

function player(id: string, name: string) {
  return { id, name, CreatedAt: "2026-01-01T00:00:00Z" };
}

function interaction(
  playerAName: string,
  playerBName: string,
  playerAMove: Move,
  playerBMove: Move,
) {
  return {
    playerAName,
    playerBName,
    playerAMove,
    playerBMove,
    playedAt: "2026-08-20T12:00:00Z",
  };
}

describe("PlayerManager", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([]);
  });

  it("renders a stats row per player derived from history", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      player("1", "Chris"),
      player("2", "Dana"),
    ]);
    // Chris cooperates twice (a 0-pt loss, then a 3-pt draw);
    // Dana betrays once (5-pt win) and cooperates once (draw)
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([
      interaction("Chris", "Dana", "C", "B"),
      interaction("Chris", "Dana", "C", "C"),
    ]);

    render(<PlayerManager />);

    const chrisRow = (await screen.findByText("Chris")).closest("tr")!;
    const cells = within(chrisRow).getAllByRole("cell");
    expect(cells[1]).toHaveTextContent("2"); // games
    expect(cells[2]).toHaveTextContent("100%"); // cooperation rate
    expect(cells[3]).toHaveTextContent("0"); // betrayals
    expect(cells[4]).toHaveTextContent("1.50"); // avg pts (0 + 3) / 2
    expect(cells[5]).toHaveTextContent("0W · 1L · 1D");

    const danaRow = screen.getByText("Dana").closest("tr")!;
    const danaCells = within(danaRow).getAllByRole("cell");
    expect(danaCells[2]).toHaveTextContent("50%");
    expect(danaCells[3]).toHaveTextContent("1");
    expect(danaCells[4]).toHaveTextContent("4.00"); // (5 + 3) / 2
    expect(danaCells[5]).toHaveTextContent("1W · 0L · 1D");
  });

  it("shows a no-rounds record for players without games", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      player("1", "Chris"),
    ]);

    render(<PlayerManager />);

    const row = (await screen.findByText("Chris")).closest("tr")!;
    expect(within(row).getByText("no rounds yet")).toBeInTheDocument();
  });

  it("has no delete controls", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      player("1", "Chris"),
    ]);

    render(<PlayerManager />);
    await screen.findByText("Chris");

    expect(
      screen.queryByRole("button", { name: /remove|delete/i }),
    ).not.toBeInTheDocument();
  });

  it("shows an error message when listPlayers fails", async () => {
    vi.spyOn(playersApi, "listPlayers").mockRejectedValue(
      new Error("network down"),
    );

    render(<PlayerManager />);

    expect(await screen.findByText(/network down/i)).toBeInTheDocument();
  });

  it("submits the prompt form, creates a player, and adds it to the table", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);
    vi.spyOn(playersApi, "createPlayer").mockResolvedValue(
      player("2", "Pepper"),
    );

    const user = userEvent.setup();
    render(<PlayerManager />);

    await waitFor(() => expect(playersApi.listPlayers).toHaveBeenCalledTimes(1));

    await user.type(screen.getByLabelText(/name/i), "Pepper");
    await user.click(screen.getByRole("button", { name: /add/i }));

    expect(await screen.findByText("Pepper")).toBeInTheDocument();
    expect(playersApi.createPlayer).toHaveBeenCalledWith("Pepper");
  });
});
