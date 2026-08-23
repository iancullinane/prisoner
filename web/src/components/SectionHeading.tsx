import type { ReactNode } from "react";

export default function SectionHeading({ children }: { children: ReactNode }) {
  return (
    <h1 className="mb-4 inline-block self-start bg-frame px-2.5 py-0.5 text-xs font-bold uppercase tracking-[.18em] text-black">
      {children}
    </h1>
  );
}
