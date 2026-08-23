import PlayerManager from "../components/PlayerManager";
import SectionHeading from "../components/SectionHeading";

export default function PlayersPage() {
  return (
    <section className="flex flex-1 flex-col">
      <SectionHeading>Players</SectionHeading>
      <PlayerManager />
    </section>
  );
}
