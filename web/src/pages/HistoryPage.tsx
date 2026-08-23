import { useEffect, useState } from "react";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { listHistory, listHistoryForPlayer } from "../api/history";
import { listPlayers } from "../api/players";
import MoveChip from "../components/MoveChip";
import SectionHeading from "../components/SectionHeading";
import {
  OUTCOME_LABELS,
  historySummary,
  outcomeKind,
  payoffFor,
  type OutcomeKind,
} from "../lib/stats";
import type { Player, PrettyInteraction } from "../types";

dayjs.extend(relativeTime);

type OutcomeFilter = "all" | OutcomeKind;

const OUTCOME_FILTERS: Array<{ value: OutcomeFilter; label: string }> = [
  { value: "all", label: "all" },
  { value: "trust", label: "mutual trust" },
  { value: "betray", label: "betrayals" },
  { value: "ruin", label: "mutual ruin" },
];

const OUTCOME_COLORS: Record<OutcomeKind, string> = {
  trust: "text-coop",
  betray: "text-betray",
  ruin: "text-frame",
};

const TH_CLASS =
  "sticky top-0 whitespace-nowrap border-b border-frame-dim bg-black px-3 py-2 text-left text-[10px] font-bold uppercase tracking-[.14em] text-frame";

function StripCell({
  label,
  warn = false,
  children,
}: {
  label: string;
  warn?: boolean;
  children: React.ReactNode;
}) {
  return (
    <div className="bg-panel-2 px-3 py-2">
      <div className="text-[10px] uppercase tracking-[.12em] text-dim">
        {label}
      </div>
      <div
        className={`mt-0.5 text-[17px] font-bold ${warn ? "text-betray" : "text-fg-bright"}`}
      >
        {children}
      </div>
    </div>
  );
}

export default function HistoryPage() {
  const [history, setHistory] = useState<PrettyInteraction[]>([]);
  const [players, setPlayers] = useState<Player[]>([]);
  const [selectedPlayerId, setSelectedPlayerId] = useState("");
  const [outcome, setOutcome] = useState<OutcomeFilter>("all");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    listPlayers()
      .then(setPlayers)
      .catch((err: Error) => setError(err.message));
  }, []);

  useEffect(() => {
    const request = selectedPlayerId
      ? listHistoryForPlayer(selectedPlayerId)
      : listHistory();

    request.then(setHistory).catch((err: Error) => setError(err.message));
  }, [selectedPlayerId]);

  const rows = history.filter(
    (r) =>
      outcome === "all" ||
      outcomeKind(r.playerAMove, r.playerBMove) === outcome,
  );
  const summary = historySummary(history);

  return (
    <section>
      <SectionHeading>History</SectionHeading>
      {error && <p role="alert" className="mb-3 text-betray">{error}</p>}

      <section
        aria-label="Summary"
        className="mb-4 grid grid-cols-[repeat(auto-fit,minmax(140px,1fr))] gap-px border border-rule bg-rule"
      >
        <StripCell label="Rounds">{summary.rounds}</StripCell>
        <StripCell label="Betrayal rate" warn>
          {summary.rounds ? `${Math.round(summary.betrayalRate * 100)}%` : "—"}
        </StripCell>
        <StripCell label="Most trusting">
          {summary.mostTrusting ? (
            <>
              {summary.mostTrusting.name}{" "}
              <small className="text-[11px] font-normal text-dim">
                {Math.round(summary.mostTrusting.rate * 100)}% coop
              </small>
            </>
          ) : (
            "—"
          )}
        </StripCell>
        <StripCell label="Most treacherous" warn>
          {summary.mostTreacherous ? (
            <>
              {summary.mostTreacherous.name}{" "}
              <small className="text-[11px] font-normal text-dim">
                {summary.mostTreacherous.betrayals} betrayals
              </small>
            </>
          ) : (
            "—"
          )}
        </StripCell>
      </section>

      <div className="mb-3 flex flex-wrap items-center gap-x-5 gap-y-2">
        <div>
          <label
            htmlFor="history-player-filter"
            className="mr-2 text-[10.5px] uppercase tracking-[.1em] text-dim"
          >
            Filter by player
          </label>
          <select
            id="history-player-filter"
            value={selectedPlayerId}
            onChange={(e) => setSelectedPlayerId(e.target.value)}
            className="min-w-[150px] border border-rule bg-black px-2 py-1 text-fg-bright focus:border-coop focus:outline-none focus:ring-1 focus:ring-coop"
          >
            <option value="">All players</option>
            {players.map((player) => (
              <option key={player.id} value={player.id}>
                {player.name}
              </option>
            ))}
          </select>
        </div>
        <div
          role="group"
          aria-label="Outcome"
          className="inline-flex border border-rule"
        >
          {OUTCOME_FILTERS.map((f) => (
            <button
              key={f.value}
              type="button"
              aria-pressed={outcome === f.value}
              onClick={() => setOutcome(f.value)}
              className={`border-r border-rule px-2.5 py-1 text-[11px] uppercase tracking-[.06em] last:border-r-0 ${
                outcome === f.value
                  ? "bg-frame font-bold text-black"
                  : "text-dim hover:text-fg"
              }`}
            >
              {f.label}
            </button>
          ))}
        </div>
      </div>

      <p className="mb-3 text-[11px] text-dimmer">
        <span className="mr-4">
          <b className="text-coop">C</b> cooperate · 3 pts each
        </span>
        <span className="mr-4">
          <b className="text-betray">B</b> betray · 5 / 0
        </span>
        <span>both betray · 1 each</span>
      </p>

      <div className="border border-rule">
        <table className="w-full border-collapse text-left text-[12.5px]">
          <thead>
            <tr>
              <th className={TH_CLASS}>Player A</th>
              <th className={TH_CLASS}>Move</th>
              <th className={TH_CLASS}>Player B</th>
              <th className={TH_CLASS}>Move</th>
              <th className={TH_CLASS}>Outcome</th>
              <th className={`${TH_CLASS} text-right`}>Pts</th>
              <th className={TH_CLASS}>Played</th>
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td colSpan={7} className="px-3 py-6 text-center text-dimmer">
                  no rounds match this filter —{" "}
                  <span className="text-dim">press [3] to play one</span>
                </td>
              </tr>
            ) : (
              rows.map((r, index) => {
                const kind = outcomeKind(r.playerAMove, r.playerBMove);
                const [pa, pb] = payoffFor(r.playerAMove, r.playerBMove);
                return (
                  <tr
                    key={index}
                    className="border-b border-rule last:border-b-0 even:bg-white/[.014] hover:bg-coop/5"
                  >
                    <td className="whitespace-nowrap px-3 py-1.5 text-fg-bright">
                      {r.playerAName}
                    </td>
                    <td className="whitespace-nowrap px-3 py-1.5">
                      <MoveChip move={r.playerAMove} />
                    </td>
                    <td className="whitespace-nowrap px-3 py-1.5 text-fg-bright">
                      {r.playerBName}
                    </td>
                    <td className="whitespace-nowrap px-3 py-1.5">
                      <MoveChip move={r.playerBMove} />
                    </td>
                    <td
                      className={`whitespace-nowrap px-3 py-1.5 text-[11px] font-bold tracking-[.08em] ${OUTCOME_COLORS[kind]}`}
                    >
                      {OUTCOME_LABELS[kind]}
                    </td>
                    <td className="whitespace-nowrap px-3 py-1.5 text-right tabular-nums text-dim">
                      {pa > pb ? <b className="text-fg-bright">{pa}</b> : pa} ·{" "}
                      {pb > pa ? <b className="text-fg-bright">{pb}</b> : pb}
                    </td>
                    <td
                      className="whitespace-nowrap px-3 py-1.5 text-[11.5px] text-dim"
                      title={dayjs(r.playedAt).format("MMM D, h:mm A")}
                    >
                      {dayjs(r.playedAt).fromNow()}
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </section>
  );
}
