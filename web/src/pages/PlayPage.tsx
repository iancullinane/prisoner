import { useEffect, useState, type FormEvent } from "react";
import { listPlayers } from "../api/players";
import { playRound } from "../api/play";
import PayoffMatrix from "../components/PayoffMatrix";
import SectionHeading from "../components/SectionHeading";
import { resultToPoints } from "../lib/stats";
import type { Move, Player, PlayResponse } from "../types";

const LEGEND_CLASS =
  "px-1.5 text-[10px] font-bold uppercase tracking-[.16em] text-frame";
const SELECT_CLASS =
  "w-full border border-rule bg-black px-2 py-1.5 text-fg-bright focus:border-coop focus:outline-none focus:ring-1 focus:ring-coop";
const SEAT_LABEL_CLASS =
  "mb-1 block text-[10px] uppercase tracking-[.12em] text-dim";

function MoveButtons({
  label,
  move,
  onChange,
}: {
  label: string;
  move: Move;
  onChange: (m: Move) => void;
}) {
  return (
    <div role="group" aria-label={label} className="flex justify-center gap-3">
      <button
        type="button"
        aria-pressed={move === "C"}
        onClick={() => onChange("C")}
        className="w-[170px] border border-coop-deep bg-coop/5 py-3 text-center text-coop hover:bg-coop/15 aria-pressed:border-coop aria-pressed:bg-coop/20"
      >
        <span className="block text-[20px] font-bold" aria-hidden="true">
          ✓
        </span>
        <span className="block text-[11px] uppercase tracking-[.14em]">
          Cooperate
        </span>
      </button>
      <button
        type="button"
        aria-pressed={move === "B"}
        onClick={() => onChange("B")}
        className="w-[170px] border border-betray-deep bg-betray/5 py-3 text-center text-betray hover:bg-betray/15 aria-pressed:border-betray aria-pressed:bg-betray/20"
      >
        <span className="block text-[20px] font-bold" aria-hidden="true">
          ✕
        </span>
        <span className="block text-[11px] uppercase tracking-[.14em]">
          Betray
        </span>
      </button>
    </div>
  );
}

export default function PlayPage() {
  const [players, setPlayers] = useState<Player[]>([]);
  const [playerA, setPlayerA] = useState("");
  const [playerB, setPlayerB] = useState("");
  const [moveA, setMoveA] = useState<Move>("C");
  const [moveB, setMoveB] = useState<Move>("C");
  const [result, setResult] = useState<PlayResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    listPlayers()
      .then(setPlayers)
      .catch((err: Error) => setError(err.message));
  }, []);

  const selfPlay = playerA !== "" && playerA === playerB;
  const nameOf = (id: string, fallback: string) =>
    players.find((p) => p.id === id)?.name ?? fallback;

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      const response = await playRound({
        playerA,
        playerB,
        playerAMove: moveA,
        playerBMove: moveB,
      });
      setResult(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to play round");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section>
      <SectionHeading>Play</SectionHeading>
      {error && <p role="alert" className="mb-3 text-betray">{error}</p>}

      <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_290px]">
        <form onSubmit={handleSubmit}>
          <fieldset className="mb-4 border border-rule px-4 pb-4 pt-3">
            <legend className={LEGEND_CLASS}>Matchup</legend>
            <div className="grid grid-cols-[1fr_34px_1fr] items-end gap-2.5">
              <div>
                <label htmlFor="player-a" className={SEAT_LABEL_CLASS}>
                  Player A
                </label>
                <select
                  id="player-a"
                  value={playerA}
                  onChange={(e) => setPlayerA(e.target.value)}
                  className={SELECT_CLASS}
                >
                  <option value="">Select player</option>
                  {players.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="pb-1.5 text-center font-bold text-frame">vs</div>
              <div>
                <label htmlFor="player-b" className={SEAT_LABEL_CLASS}>
                  Player B
                </label>
                <select
                  id="player-b"
                  value={playerB}
                  onChange={(e) => setPlayerB(e.target.value)}
                  className={SELECT_CLASS}
                >
                  <option value="">Select player</option>
                  {players.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            {selfPlay && (
              <p className="mt-3 border-l-2 border-frame bg-frame/5 px-3 py-2 text-[11.5px] text-betray">
                A player can't face themselves.
              </p>
            )}
          </fieldset>

          <fieldset className="mb-4 border border-rule px-4 pb-5 pt-3">
            <legend className={LEGEND_CLASS}>Round</legend>
            <div className="flex flex-col gap-5 pt-1">
              <div>
                <p className="mb-2 text-center text-[11px] uppercase tracking-[.16em] text-frame">
                  {nameOf(playerA, "Player A")}'s move
                </p>
                <MoveButtons label="Move for A" move={moveA} onChange={setMoveA} />
              </div>
              <div>
                <p className="mb-2 text-center text-[11px] uppercase tracking-[.16em] text-frame">
                  {nameOf(playerB, "Player B")}'s move
                </p>
                <MoveButtons label="Move for B" move={moveB} onChange={setMoveB} />
              </div>
              <button
                type="submit"
                disabled={submitting || !playerA || !playerB || selfPlay}
                className="mx-auto w-[230px] bg-coop py-2 text-[11.5px] font-bold uppercase tracking-[.1em] text-black hover:bg-fg-bright disabled:cursor-not-allowed disabled:bg-transparent disabled:text-dim disabled:outline disabled:outline-1 disabled:outline-rule"
              >
                {submitting ? "Playing…" : "Play round"}
              </button>
            </div>
          </fieldset>

          {result && (
            <p className="text-center text-[13px]">
              Score — {nameOf(playerA, "Player A")}:{" "}
              <b className="text-fg-bright">
                {resultToPoints(result.playerAScore)} pts
              </b>{" "}
              · {nameOf(playerB, "Player B")}:{" "}
              <b className="text-fg-bright">
                {resultToPoints(result.playerBScore)} pts
              </b>
            </p>
          )}
        </form>

        <div>
          <fieldset className="border border-rule px-4 pb-4 pt-3">
            <legend className={LEGEND_CLASS}>Payoff</legend>
            <PayoffMatrix />
            <p className="mt-3 text-[11px] leading-[1.75] text-dim">
              Betraying always pays more <i>this</i> round — which is why both
              of you do it, and both of you end up with 1.
              <br />
              <br />
              Over many rounds, reputation is the only thing that beats the
              math.
            </p>
          </fieldset>
        </div>
      </div>
    </section>
  );
}
