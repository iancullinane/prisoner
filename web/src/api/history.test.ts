import { describe, it, expect, vi, beforeEach } from "vitest";
import { listHistory, listHistoryForPlayer } from "./history";

const sample = [
  {
    playerAName: "Alice",
    playerBName: "Bob",
    playerAMove: "C",
    playerBMove: "B",
    playedAt: "2026-01-01T00:00:00Z",
  },
];

describe("listHistory", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("returns history parsed from the response body", async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify(sample), { status: 200 }),
    );

    const result = await listHistory();

    expect(fetch).toHaveBeenCalledWith(expect.stringMatching(/\/api\/v1\/history$/));
    expect(result).toEqual(sample);
  });

  it("throws when the response is not ok", async () => {
    vi.mocked(fetch).mockResolvedValue(new Response("boom", { status: 500 }));

    await expect(listHistory()).rejects.toThrow();
  });
});

describe("listHistoryForPlayer", () => {
  beforeEach(() => {
    vi.stubGlobal("fetch", vi.fn());
  });

  it("returns history for a single player", async () => {
    vi.mocked(fetch).mockResolvedValue(
      new Response(JSON.stringify(sample), { status: 200 }),
    );

    const result = await listHistoryForPlayer("p1");

    expect(fetch).toHaveBeenCalledWith(
      expect.stringMatching(/\/api\/v1\/history\/p1$/),
    );
    expect(result).toEqual(sample);
  });

  it("throws when the response is not ok", async () => {
    vi.mocked(fetch).mockResolvedValue(new Response("boom", { status: 500 }));

    await expect(listHistoryForPlayer("p1")).rejects.toThrow();
  });
});
