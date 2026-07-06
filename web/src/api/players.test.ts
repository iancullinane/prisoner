import { describe, it, expect, vi, beforeEach } from "vitest";
import { listPlayers, createPlayer } from "./players";

describe("listPlayers", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("returns players parsed from the response body", async () => {
    const players = [{ id: "1", name: "Chris", CreatedAt: "2026-01-01T00:00:00Z" }];
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify(players), { status: 200 }),
    );

    const result = await listPlayers();

    expect(fetch).toHaveBeenCalledWith(expect.stringMatching(/\/players$/));
    expect(result).toEqual(players);
  });

  it("throws when the response is not ok", async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response("boom", { status: 500 }),
    );

    await expect(listPlayers()).rejects.toThrow();
  });
});

describe("createPlayer", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("posts to /players/{name} and returns the created player", async () => {
    const player = { id: "2", name: "Pepper", CreatedAt: "2026-01-01T00:00:00Z" };
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify(player), { status: 200 }),
    );

    const result = await createPlayer("Pepper");

    expect(fetch).toHaveBeenCalledWith(
      expect.stringMatching(/\/players\/Pepper$/),
      expect.objectContaining({ method: "POST" }),
    );
    expect(result).toEqual(player);
  });

  it("throws when the response is not ok", async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response("boom", { status: 500 }),
    );

    await expect(createPlayer("Pepper")).rejects.toThrow();
  });
});
