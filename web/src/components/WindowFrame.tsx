import { useEffect, useState, type ReactNode } from "react";
import { useLocation } from "react-router-dom";

function formatClock(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

export default function WindowFrame({ children }: { children: ReactNode }) {
  const { pathname } = useLocation();
  const [now, setNow] = useState(() => new Date());

  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  return (
    <div className="fixed inset-0 flex items-center justify-center md:p-[max(1rem,3vmin)]">
      <div className="flex h-full w-full max-h-[860px] max-w-[1280px] flex-col border border-hacker-ochre bg-hacker-bg shadow-[0_0_24px_rgb(43_217_74_/_0.15)]">
        <header className="flex items-center justify-between border-b border-hacker-ochre bg-hacker-chrome px-3 py-1.5 font-army text-sm uppercase tracking-wide text-hacker-ochre">
          <span>PRISONER — ~{pathname}</span>
          <time>{formatClock(now)}</time>
        </header>
        <div className="flex min-h-0 flex-1 flex-col md:flex-row">
          {children}
        </div>
      </div>
    </div>
  );
}
