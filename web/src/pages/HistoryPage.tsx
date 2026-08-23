import { useEffect, useState } from "react";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { listHistory, listHistoryForPlayer } from "../api/history";
import { listPlayers } from "../api/players";
import SectionHeading from "../components/SectionHeading";
import type { Player, PrettyInteraction } from "../types";

dayjs.extend(relativeTime);

export default function HistoryPage() {
  const [history, setHistory] = useState<PrettyInteraction[]>([]);
  const [players, setPlayers] = useState<Player[]>([]);
  const [selectedPlayerId, setSelectedPlayerId] = useState("");
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

  return (
    <section>
      <SectionHeading>History</SectionHeading>
      {error && <p role="alert">{error}</p>}

      <label htmlFor="history-player-filter" className="mr-2">
        Filter by player
      </label>
      <select
        id="history-player-filter"
        value={selectedPlayerId}
        onChange={(e) => setSelectedPlayerId(e.target.value)}
        className="mb-4 border border-hacker-green/50 bg-hacker-chrome text-hacker-green"
      >
        <option value="">All players</option>
        {players.map((player) => (
          <option key={player.id} value={player.id}>
            {player.name}
          </option>
        ))}
      </select>

      <table className="w-full border-collapse text-left">
        <thead>
          <tr className="text-sm uppercase text-hacker-ochre">
            <th className="border-b border-hacker-green/30 py-2">Player A</th>
            <th className="border-b border-hacker-green/30 py-2">Player B</th>
            <th className="border-b border-hacker-green/30 py-2">Moves</th>
            <th className="border-b border-hacker-green/30 py-2">Played At</th>
          </tr>
        </thead>
        <tbody>
          {history.map((interaction, index) => (
            <tr key={index}>
              <td className="py-1">{interaction.playerAName}</td>
              <td className="py-1">{interaction.playerBName}</td>
              <td className="py-1">
                {interaction.playerAMove} / {interaction.playerBMove}
              </td>
              <td className="py-1" title={dayjs(interaction.playedAt).format("MMM D, h:mm A")}>
                {dayjs(interaction.playedAt).fromNow()}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
