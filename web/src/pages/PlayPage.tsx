import { useEffect, useState, type FormEvent } from "react";
import { listPlayers } from "../api/players";
import { playRound } from "../api/play";
import type { Move, Player, PlayResponse } from "../types";

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
      <h1 className="font-army mb-4 text-2xl text-hacker-green">Play</h1>
      {error && <p role="alert">{error}</p>}

      <form onSubmit={handleSubmit} className="flex max-w-md flex-col gap-4">
        <label htmlFor="player-a">Player A</label>
        <select
          id="player-a"
          value={playerA}
          onChange={(e) => setPlayerA(e.target.value)}
          className="border border-hacker-green/50 bg-hacker-chrome text-hacker-green"
        >
          <option value="">Select player</option>
          {players.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        <div role="group" aria-label="Move for A">
          <button type="button" aria-pressed={moveA === "C"} onClick={() => setMoveA("C")}>
            Cooperate
          </button>
          <button type="button" aria-pressed={moveA === "B"} onClick={() => setMoveA("B")}>
            Betray
          </button>
        </div>

        <label htmlFor="player-b">Player B</label>
        <select
          id="player-b"
          value={playerB}
          onChange={(e) => setPlayerB(e.target.value)}
          className="border border-hacker-green/50 bg-hacker-chrome text-hacker-green"
        >
          <option value="">Select player</option>
          {players.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        <div role="group" aria-label="Move for B">
          <button type="button" aria-pressed={moveB === "C"} onClick={() => setMoveB("C")}>
            Cooperate
          </button>
          <button type="button" aria-pressed={moveB === "B"} onClick={() => setMoveB("B")}>
            Betray
          </button>
        </div>

        <button
          type="submit"
          disabled={submitting || !playerA || !playerB || selfPlay}
          className="bg-hacker-green px-4 py-2 font-army text-hacker-bg disabled:opacity-40"
        >
          {submitting ? "Playing…" : "Play round"}
        </button>
      </form>

      {result && (
        <p className="mt-4">
          Score — Player A: {result.playerAScore}, Player B: {result.playerBScore}
        </p>
      )}
    </section>
  );
}
