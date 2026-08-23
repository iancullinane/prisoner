import type { Move } from "../types";

// Color + glyph + letter together, so the move is never encoded by color alone.
export default function MoveChip({ move }: { move: Move }) {
  const coop = move === "C";
  return (
    <span
      aria-label={coop ? "cooperate" : "betray"}
      className={`inline-block min-w-[52px] border px-1.5 text-center text-[12px] font-bold tracking-[.06em] ${
        coop
          ? "border-coop-deep bg-coop/10 text-coop"
          : "border-betray-deep bg-betray/10 text-betray"
      }`}
    >
      {coop ? "✓ C" : "✕ B"}
    </span>
  );
}
