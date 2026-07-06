import { useEffect, useState, type FormEvent } from "react";
import { createPlayer, listPlayers } from "../api/players";
import type { Player } from "../types";

export default function PlayerManager() {
  const [players, setPlayers] = useState<Player[]>([]);
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    listPlayers()
      .then(setPlayers)
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
    <div>
      {error && <p role="alert">{error}</p>}
      <ul>
        {players.map((player) => (
          <li key={player.id}>{player.name}</li>
        ))}
      </ul>
      <form onSubmit={handleSubmit}>
        <label htmlFor="player-name">Name</label>
        <input
          id="player-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <button type="submit">Add player</button>
      </form>
    </div>
  );
}
