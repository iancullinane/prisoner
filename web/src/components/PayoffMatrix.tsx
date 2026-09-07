import { PAYOFF } from "../lib/stats";

const TH_CLASS =
  "border border-rule px-1.5 py-1.5 text-[10.5px] font-medium uppercase tracking-[.1em] text-dim";

function Cell({
  pair,
  caption,
  tone,
}: {
  pair: [number, number];
  caption: string;
  tone: "cc" | "bb" | "mix";
}) {
  const colors = {
    cc: ["text-coop", "text-coop-deep"],
    bb: ["text-frame", "text-frame-dim"],
    mix: ["text-betray", "text-betray-deep"],
  }[tone];
  return (
    <td
      className={`border border-rule px-1.5 py-1.5 text-center font-bold tabular-nums ${colors[0]}`}
    >
      {pair[0]} · {pair[1]}
      <small
        className={`block text-[10px] font-normal tracking-[.08em] ${colors[1]}`}
      >
        {caption}
      </small>
    </td>
  );
}

export default function PayoffMatrix() {
  return (
    <table
      aria-label="Payoff matrix"
      className="w-full border-collapse text-[12.5px]"
    >
      <tbody>
        <tr>
          <th className={TH_CLASS} />
          <th className={TH_CLASS} scope="col">
            They cooperate
          </th>
          <th className={TH_CLASS} scope="col">
            They betray
          </th>
        </tr>
        <tr>
          <th className={TH_CLASS} scope="row">
            You cooperate
          </th>
          <Cell pair={PAYOFF.CC} caption="mutual trust" tone="cc" />
          <Cell pair={PAYOFF.CB} caption="you're played" tone="mix" />
        </tr>
        <tr>
          <th className={TH_CLASS} scope="row">
            You betray
          </th>
          <Cell pair={PAYOFF.BC} caption="you cash in" tone="mix" />
          <Cell pair={PAYOFF.BB} caption="mutual ruin" tone="bb" />
        </tr>
      </tbody>
    </table>
  );
}
