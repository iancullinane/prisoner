import { NavLink } from "react-router-dom";

const links = [
  { to: "/players", label: "Players", key: "1" },
  { to: "/history", label: "History", key: "2" },
  { to: "/play", label: "Play", key: "3" },
];

export default function Sidebar() {
  return (
    <nav className="flex gap-2 border-b border-rule bg-panel-2 px-2 py-3 md:w-[210px] md:flex-col md:gap-0 md:border-b-0 md:border-r md:px-0 md:py-3.5">
      {links.map((link) => (
        <NavLink
          key={link.to}
          to={link.to}
          className={({ isActive }) =>
            `flex items-center gap-2.5 border-l-[3px] px-3 py-2 text-[11.5px] uppercase tracking-[.13em] ${
              isActive
                ? "border-frame bg-frame/5 text-frame"
                : "border-transparent text-dim hover:bg-coop/5 hover:text-fg"
            }`
          }
        >
          <span className="text-[10.5px] text-dimmer" aria-hidden="true">
            [{link.key}]
          </span>
          {link.label}
        </NavLink>
      ))}
      <div className="mt-auto hidden px-4 text-[10.5px] leading-[1.9] text-dimmer md:block">
        <div>
          <span className="inline-block w-3.5 text-center font-bold text-coop">
            C
          </span>{" "}
          cooperate
        </div>
        <div>
          <span className="inline-block w-3.5 text-center font-bold text-betray">
            B
          </span>{" "}
          betray
        </div>
      </div>
    </nav>
  );
}
