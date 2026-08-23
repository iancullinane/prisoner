import type { Move, PrettyInteraction, Result } from "../types";

export const PAYOFF: Record<`${Move}${Move}`, [number, number]> = {
  CC: [3, 3],
  CB: [0, 5],
  BC: [5, 0],
  BB: [1, 1],
};

export function payoffFor(a: Move, b: Move): [number, number] {
  return PAYOFF[`${a}${b}`];
}

export type OutcomeKind = "trust" | "betray" | "ruin";

export const OUTCOME_LABELS: Record<OutcomeKind, string> = {
  trust: "Mutual trust",
  betray: "Betrayal",
  ruin: "Mutual ruin",
};

export function outcomeKind(a: Move, b: Move): OutcomeKind {
  if (a === b) return a === "C" ? "trust" : "ruin";
  return "betray";
}

const RESULT_POINTS: Record<Result, number> = { T: 5, R: 3, P: 1, S: 0 };

export function resultToPoints(r: Result): number {
  return RESULT_POINTS[r];
}

export interface PlayerStats {
  games: number;
  coop: number;
  betray: number;
  rate: number;
  avgPts: number;
  w: number;
  l: number;
  d: number;
}

// Interactions carry names, not ids, so stats join by name; players sharing
// a name would merge.
export function playerStats(
  name: string,
  history: PrettyInteraction[],
): PlayerStats {
  let games = 0;
  let coop = 0;
  let pts = 0;
  let w = 0;
  let l = 0;
  let d = 0;

  for (const r of history) {
    const isA = r.playerAName === name;
    if (!isA && r.playerBName !== name) continue;
    games++;
    const mine = isA ? r.playerAMove : r.playerBMove;
    if (mine === "C") coop++;
    const [pa, pb] = payoffFor(r.playerAMove, r.playerBMove);
    const [myPts, theirPts] = isA ? [pa, pb] : [pb, pa];
    pts += myPts;
    if (myPts > theirPts) w++;
    else if (myPts < theirPts) l++;
    else d++;
  }

  return {
    games,
    coop,
    betray: games - coop,
    rate: games ? coop / games : 0,
    avgPts: games ? pts / games : 0,
    w,
    l,
    d,
  };
}

export interface HistorySummary {
  rounds: number;
  betrayalRate: number;
  mostTrusting: { name: string; rate: number } | null;
  mostTreacherous: { name: string; betrayals: number } | null;
}

export function historySummary(history: PrettyInteraction[]): HistorySummary {
  if (!history.length) {
    return { rounds: 0, betrayalRate: 0, mostTrusting: null, mostTreacherous: null };
  }

  const betrayals = history.reduce(
    (n, r) =>
      n + (r.playerAMove === "B" ? 1 : 0) + (r.playerBMove === "B" ? 1 : 0),
    0,
  );

  const names = [
    ...new Set(history.flatMap((r) => [r.playerAName, r.playerBName])),
  ];
  const ranked = names.map((name) => ({ name, stats: playerStats(name, history) }));

  const best = ranked.reduce((a, b) => (b.stats.rate > a.stats.rate ? b : a));
  const worst = ranked.reduce((a, b) =>
    b.stats.betray > a.stats.betray ? b : a,
  );

  return {
    rounds: history.length,
    betrayalRate: betrayals / (history.length * 2),
    mostTrusting: { name: best.name, rate: best.stats.rate },
    mostTreacherous: { name: worst.name, betrayals: worst.stats.betray },
  };
}
