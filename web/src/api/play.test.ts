import { describe, it, expect, vi, beforeEach } from "vitest";
import { playRound } from "./play";
import type { PlayRequest } from "../types";

const request: PlayRequest = {
  playerA: "p1",
  playerB: "p2",
  playerAMove: "C",
  playerBMove: "B",
};

describe("playRound", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("posts to /api/v1/play and returns the response body", async () => {
    const response = { id: "i1", playerAScore: "S", playerBScore: "T" };
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify(response), { status: 201 }),
    );

    const result = await playRound(request);

    expect(fetch).toHaveBeenCalledWith(
      expect.stringMatching(/\/api\/v1\/play$/),
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify(request),
      }),
    );
    expect(result).toEqual(response);
  });

  it("throws when the response is not ok", async () => {
    vi.mocked(fetch).mockResolvedValue(new Response("boom", { status: 400 }));

    await expect(playRound(request)).rejects.toThrow();
  });
});
