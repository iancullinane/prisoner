import { useEffect, useState, type FormEvent } from "react";
import { listHistory } from "../api/history";
import { createPlayer, listPlayers } from "../api/players";
import { playerStats } from "../lib/stats";
import type { Player, PrettyInteraction } from "../types";

const TH_CLASS =
  "whitespace-nowrap border-b border-frame-dim bg-black px-3 py-2 text-left text-[10px] font-bold uppercase tracking-[.14em] text-frame";

export default function PlayerManager() {
  const [players, setPlayers] = useState<Player[]>([]);
  const [history, setHistory] = useState<PrettyInteraction[]>([]);
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    listPlayers()
      .then(setPlayers)
      .catch((err: Error) => setError(err.message));
    listHistory()
      .then(setHistory)
      .catch((err: Error) => setError(err.message));
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      const player = await createPlayer(name);
      setPlayers((prev) => [...prev, player]);
      setName("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create player");
    }
  }

  return (
    <div className="flex flex-1 flex-col">
      {error && <p role="alert" className="mb-3 text-betray">{error}</p>}

      <div className="border border-rule">
        <table className="w-full border-collapse text-left text-[12.5px]">
          <thead>
            <tr>
              <th className={TH_CLASS}>Player</th>
              <th className={`${TH_CLASS} text-right`}>Games</th>
              <th className={TH_CLASS}>Cooperation rate</th>
              <th className={`${TH_CLASS} text-right`}>Betrayals</th>
              <th className={`${TH_CLASS} text-right`}>Avg pts</th>
              <th className={TH_CLASS}>Record</th>
            </tr>
          </thead>
          <tbody>
            {players.map((player) => {
              const s = playerStats(player.name, history);
              const pct = Math.round(s.rate * 100);
              const rateColor = s.rate >= 0.5 ? "text-coop" : "text-betray";
              const barColor = s.rate >= 0.5 ? "bg-coop" : "bg-betray";
              return (
                <tr
                  key={player.id}
                  className="border-b border-rule last:border-b-0 even:bg-white/[.014] hover:bg-coop/5"
                >
                  <td className="whitespace-nowrap px-3 py-1.5 text-fg-bright">
                    {player.name}
                  </td>
                  <td className="whitespace-nowrap px-3 py-1.5 text-right tabular-nums">
                    {s.games}
                  </td>
                  <td className="whitespace-nowrap px-3 py-1.5">
                    <span className="mr-2 inline-block h-[7px] w-[74px] bg-rule align-middle">
                      <span
                        className={`block h-full ${barColor}`}
                        style={{ width: `${pct}%` }}
                      />
                    </span>
                    <span className={s.games ? rateColor : "text-dim"}>
                      {s.games ? `${pct}%` : "—"}
                    </span>
                  </td>
                  <td
                    className={`whitespace-nowrap px-3 py-1.5 text-right tabular-nums ${
                      s.betray ? "text-betray" : "text-dim"
                    }`}
                  >
                    {s.betray}
                  </td>
                  <td className="whitespace-nowrap px-3 py-1.5 text-right tabular-nums">
                    {s.games ? s.avgPts.toFixed(2) : "—"}
                  </td>
                  <td className="whitespace-nowrap px-3 py-1.5 text-[11.5px] text-dim">
                    {s.games ? `${s.w}W · ${s.l}L · ${s.d}D` : "no rounds yet"}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>

        <form
          onSubmit={handleSubmit}
          className="flex items-center gap-2 border-t border-rule bg-black px-3 py-2"
        >
          <span className="font-bold text-coop" aria-hidden="true">
            &gt;
          </span>
          <input
            id="player-name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="add player…"
            autoComplete="off"
            aria-label="New player name"
            className="min-w-0 flex-1 border-0 bg-transparent text-fg-bright placeholder:text-dimmer focus:outline-none"
          />
          <span
            className="inline-block h-3.5 w-[7px] animate-blink bg-coop"
            aria-hidden="true"
          />
          <button
            type="submit"
            className="bg-coop px-4 py-1.5 text-[11.5px] font-bold uppercase tracking-[.1em] text-black hover:bg-fg-bright"
          >
            Add
          </button>
        </form>
      </div>
    </div>
  );
}
