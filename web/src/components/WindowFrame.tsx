import { useEffect, useState, type ReactNode } from "react";
import { useLocation } from "react-router-dom";
import { listHistory } from "../api/history";
import { listPlayers } from "../api/players";
import Statusline from "./Statusline";

function formatClock(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

export default function WindowFrame({ children }: { children: ReactNode }) {
  const { pathname } = useLocation();
  const [now, setNow] = useState(() => new Date());
  const [rounds, setRounds] = useState<number | null>(null);
  const [playerCount, setPlayerCount] = useState<number | null>(null);

  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  // titlebar totals refresh on navigation, not after in-page plays
  useEffect(() => {
    let cancelled = false;
    listHistory()
      .then((h) => !cancelled && setRounds(h.length))
      .catch(() => {});
    listPlayers()
      .then((p) => !cancelled && setPlayerCount(p.length))
      .catch(() => {});
    return () => {
      cancelled = true;
    };
  }, [pathname]);

  return (
    <div className="fixed inset-x-0 top-0 flex h-dvh items-center justify-center md:p-[max(1rem,3vmin)]">
      <div className="flex h-full w-full max-h-[860px] max-w-[1280px] flex-col border border-frame bg-panel shadow-[0_0_0_1px_rgba(232,163,61,.14),0_0_40px_rgba(74,222,128,.05)]">
        <header className="flex items-center justify-between gap-4 border-b border-frame bg-black px-3.5 py-1.5 text-[11px] font-bold uppercase tracking-[.14em] text-frame">
          <span>prisoner — ~{pathname}</span>
          <span className="flex items-center gap-4 text-frame-dim">
            <span className="max-sm:hidden">
              rounds <b className="text-frame">{rounds ?? "—"}</b>
            </span>
            <span className="max-sm:hidden">
              players <b className="text-frame">{playerCount ?? "—"}</b>
            </span>
            <time>{formatClock(now)}</time>
          </span>
        </header>
        <div className="flex min-h-0 flex-1 flex-col md:flex-row">{children}</div>
        <Statusline />
      </div>
    </div>
  );
}
