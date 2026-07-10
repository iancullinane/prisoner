import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import App from "./App";
import * as playersApi from "./api/players";

describe("App", () => {
  it("redirects the default route to the players page", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);

    render(<App />);

    expect(
      await screen.findByRole("heading", { name: /players/i }),
    ).toBeInTheDocument();
  });

  it("renders the sidebar navigation", async () => {
    vi.spyOn(playersApi, "listPlayers").mockResolvedValue([]);

    render(<App />);

    expect(screen.getByRole("link", { name: /history/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /^play$/i })).toBeInTheDocument();
  });
});
