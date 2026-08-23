import { describe, expect, it } from "vitest";
import type { Move, PrettyInteraction, Result } from "../types";
import {
  OUTCOME_LABELS,
  PAYOFF,
  historySummary,
  outcomeKind,
  payoffFor,
  playerStats,
  resultToPoints,
} from "./stats";

function round(
  playerAName: string,
  playerBName: string,
  playerAMove: Move,
  playerBMove: Move,
  playedAt = "2026-08-20T12:00:00Z",
): PrettyInteraction {
  return { playerAName, playerBName, playerAMove, playerBMove, playedAt };
}

describe("payoffFor", () => {
  const cases: Array<{ a: Move; b: Move; want: [number, number] }> = [
    { a: "C", b: "C", want: [3, 3] },
    { a: "C", b: "B", want: [0, 5] },
    { a: "B", b: "C", want: [5, 0] },
    { a: "B", b: "B", want: [1, 1] },
  ];

  it.each(cases)("$a vs $b pays $want", ({ a, b, want }) => {
    expect(payoffFor(a, b)).toEqual(want);
    expect(PAYOFF[`${a}${b}`]).toEqual(want);
  });
});

describe("outcomeKind", () => {
  const cases: Array<{
    a: Move;
    b: Move;
    kind: "trust" | "betray" | "ruin";
    label: string;
  }> = [
    { a: "C", b: "C", kind: "trust", label: "Mutual trust" },
    { a: "C", b: "B", kind: "betray", label: "Betrayal" },
    { a: "B", b: "C", kind: "betray", label: "Betrayal" },
    { a: "B", b: "B", kind: "ruin", label: "Mutual ruin" },
  ];

  it.each(cases)("$a vs $b is $kind", ({ a, b, kind, label }) => {
    expect(outcomeKind(a, b)).toBe(kind);
    expect(OUTCOME_LABELS[kind]).toBe(label);
  });
});

describe("resultToPoints", () => {
  const cases: Array<{ r: Result; pts: number }> = [
    { r: "T", pts: 5 },
    { r: "R", pts: 3 },
    { r: "P", pts: 1 },
    { r: "S", pts: 0 },
  ];

  it.each(cases)("$r is worth $pts", ({ r, pts }) => {
    expect(resultToPoints(r)).toBe(pts);
  });
});

describe("playerStats", () => {
  const history: PrettyInteraction[] = [
    round("ian", "steve", "C", "B"), // ian 0 (loss), steve 5 (win)
    round("ian", "steve", "C", "C"), // draw, 3 each
    round("gus", "ian", "B", "B"), // draw, 1 each
  ];

  it("aggregates games, moves, points, and record by name", () => {
    expect(playerStats("ian", history)).toEqual({
      games: 3,
      coop: 2,
      betray: 1,
      rate: 2 / 3,
      avgPts: (0 + 3 + 1) / 3,
      w: 0,
      l: 1,
      d: 2,
    });
    expect(playerStats("steve", history)).toEqual({
      games: 2,
      coop: 1,
      betray: 1,
      rate: 1 / 2,
      avgPts: 4,
      w: 1,
      l: 0,
      d: 1,
    });
  });

  it("returns zeroes for a player with no rounds", () => {
    expect(playerStats("nobody", history)).toEqual({
      games: 0,
      coop: 0,
      betray: 0,
      rate: 0,
      avgPts: 0,
      w: 0,
      l: 0,
      d: 0,
    });
  });
});

describe("historySummary", () => {
  it("summarizes rounds, betrayal rate, and extremes", () => {
    const history = [
      round("ian", "steve", "C", "B"),
      round("ian", "gus", "C", "C"),
      round("gus", "steve", "B", "B"),
    ];
    // 6 moves, 3 betrayals
    const summary = historySummary(history);
    expect(summary.rounds).toBe(3);
    expect(summary.betrayalRate).toBeCloseTo(0.5);
    expect(summary.mostTrusting).toEqual({ name: "ian", rate: 1 });
    // steve: 2 betrayals; gus: 1 betrayal → steve is most treacherous
    expect(summary.mostTreacherous).toEqual({ name: "steve", betrayals: 2 });
  });

  it("handles an empty history", () => {
    expect(historySummary([])).toEqual({
      rounds: 0,
      betrayalRate: 0,
      mostTrusting: null,
      mostTreacherous: null,
    });
  });
});
