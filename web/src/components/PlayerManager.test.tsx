import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import PlayerManager from "./PlayerManager";
import * as playersApi from "../api/players";

describe("PlayerManager", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("renders the players returned by listPlayers", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      { id: "1", name: "Chris", CreatedAt: "2026-01-01T00:00:00Z" },
    ]);

    render(<PlayerManager />);

    expect(await screen.findByText("Chris")).toBeInTheDocument();
  });

  it("shows an error message when listPlayers fails", async () => {
    vi.spyOn(playersApi, "listPlayers").mockRejectedValue(new Error("network down"));

    render(<PlayerManager />);

    expect(await screen.findByText(/network down/i)).toBeInTheDocument();
  });

  it("submits the form, creates a player, and adds it to the list", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);
    vi.spyOn(playersApi, "createPlayer").mockResolvedValue({
      id: "2",
      name: "Pepper",
      CreatedAt: "2026-01-01T00:00:00Z",
    });

    const user = userEvent.setup();
    render(<PlayerManager />);

    await waitFor(() => expect(playersApi.listPlayers).toHaveBeenCalledTimes(1));

    await user.type(screen.getByLabelText(/name/i), "Pepper");
    await user.click(screen.getByRole("button", { name: /add player/i }));

    expect(await screen.findByText("Pepper")).toBeInTheDocument();
    expect(playersApi.createPlayer).toHaveBeenCalledWith("Pepper");
  });
});
