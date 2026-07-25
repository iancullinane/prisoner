import PlayerManager from "../components/PlayerManager";

export default function PlayersPage() {
  return (
    <section className="flex flex-1 flex-col">
      <h1 className="font-army mb-4 text-2xl text-hacker-green">Players</h1>
      <PlayerManager />
    </section>
  );
}
