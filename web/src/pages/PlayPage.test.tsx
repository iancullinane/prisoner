import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import PlayPage from "./PlayPage";
import * as playersApi from "../api/players";
import * as playApi from "../api/play";
import type { PlayResponse } from "../types";

const players = [
  { id: "p1", name: "Chris", CreatedAt: "2026-01-01T00:00:00Z" },
  { id: "p2", name: "Pepper", CreatedAt: "2026-01-01T00:00:00Z" },
];

describe("PlayPage", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("submits a play and renders the returned score", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue(players);
    vi.spyOn(playApi, "playRound").mockResolvedValue({
      id: "i1",
      playerAScore: "S",
      playerBScore: "T",
    });

    const user = userEvent.setup();
    render(<PlayPage />);

    await waitFor(() => expect(playersApi.listPlayers).toHaveBeenCalledTimes(1));

    await user.selectOptions(screen.getByLabelText(/player a/i), "p1");
    await user.selectOptions(screen.getByLabelText(/player b/i), "p2");
    await user.click(screen.getByRole("button", { name: /play round/i }));

    // S/T letters are translated to points (S=0, T=5)
    const score = await screen.findByText(/score/i);
    expect(score).toHaveTextContent("Chris: 0 pts");
    expect(score).toHaveTextContent("Pepper: 5 pts");
    expect(playApi.playRound).toHaveBeenCalledWith({
      playerA: "p1",
      playerB: "p2",
      playerAMove: "C",
      playerBMove: "C",
    });
  });

  it("shows the payoff matrix reference", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue(players);

    render(<PlayPage />);

    expect(
      await screen.findByRole("table", { name: /payoff/i }),
    ).toBeInTheDocument();
    expect(screen.getByText("mutual trust")).toBeInTheDocument();
    expect(screen.getByText("mutual ruin")).toBeInTheDocument();
  });

  it("disables submit while a request is in flight", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue(players);
    let resolvePlay: (value: PlayResponse) => void;
    vi.spyOn(playApi, "playRound").mockReturnValue(
      new Promise((resolve) => {
        resolvePlay = resolve;
      }) as ReturnType<typeof playApi.playRound>,
    );

    const user = userEvent.setup();
    render(<PlayPage />);

    await waitFor(() => expect(playersApi.listPlayers).toHaveBeenCalledTimes(1));

    await user.selectOptions(screen.getByLabelText(/player a/i), "p1");
    await user.selectOptions(screen.getByLabelText(/player b/i), "p2");
    await user.click(screen.getByRole("button", { name: /play round/i }));

    expect(screen.getByRole("button", { name: /playing/i })).toBeDisabled();

    resolvePlay!({ id: "i1", playerAScore: "R", playerBScore: "R" });
    await waitFor(() =>
      expect(screen.getByRole("button", { name: /play round/i })).not.toBeDisabled(),
    );
  });

  it("shows an error message when playRound fails", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue(players);
    vi.spyOn(playApi, "playRound").mockRejectedValue(new Error("could not record interaction"));

    const user = userEvent.setup();
    render(<PlayPage />);

    await waitFor(() => expect(playersApi.listPlayers).toHaveBeenCalledTimes(1));

    await user.selectOptions(screen.getByLabelText(/player a/i), "p1");
    await user.selectOptions(screen.getByLabelText(/player b/i), "p2");
    await user.click(screen.getByRole("button", { name: /play round/i }));

    expect(await screen.findByText(/could not record interaction/i)).toBeInTheDocument();
  });

  it("disables submit when player A and player B are the same", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue(players);

    const user = userEvent.setup();
    render(<PlayPage />);

    await waitFor(() => expect(playersApi.listPlayers).toHaveBeenCalledTimes(1));

    await user.selectOptions(screen.getByLabelText(/player a/i), "p1");
    await user.selectOptions(screen.getByLabelText(/player b/i), "p1");

    expect(screen.getByRole("button", { name: /play round/i })).toBeDisabled();
  });
});
