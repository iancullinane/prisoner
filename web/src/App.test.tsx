import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";
import * as playersApi from "./api/players";
import * as historyApi from "./api/history";

describe("App", () => {
  beforeEach(() => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);
    vi.spyOn(historyApi, "listHistory").mockResolvedValue([]);
  });

  it("redirects the default route to the players page", async () => {
    render(<App />);

    expect(
      await screen.findByRole("heading", { name: /players/i }),
    ).toBeInTheDocument();
  });

  it("renders the sidebar navigation", async () => {
    render(<App />);

    expect(screen.getByRole("link", { name: /history/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /play$/i })).toBeInTheDocument();
  });
});
