import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import PlayersPage from "./PlayersPage";
import * as playersApi from "../api/players";

describe("PlayersPage", () => {
  it("renders a Players heading and the player list", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([
      { id: "1", name: "Chris", CreatedAt: "2026-01-01T00:00:00Z" },
    ]);

    render(<PlayersPage />);

    expect(screen.getByRole("heading", { name: /players/i })).toBeInTheDocument();
    expect(await screen.findByText("Chris")).toBeInTheDocument();
  });
});
