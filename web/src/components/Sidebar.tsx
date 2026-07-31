import { NavLink } from "react-router-dom";

const links = [
  { to: "/players", label: "Players" },
  { to: "/history", label: "History" },
  { to: "/play", label: "Play" },
];

export default function Sidebar() {
  return (
    <nav className="flex gap-4 border-b border-hacker-green/30 bg-hacker-chrome p-4 md:w-48 md:flex-col md:border-b-0 md:border-r">
      {links.map((link) => (
        <NavLink
          key={link.to}
          to={link.to}
          className={({ isActive }) =>
            `font-army px-2 py-1 uppercase tracking-wide ${
              isActive
                ? "border-l-2 border-hacker-ochre text-hacker-ochre"
                : "text-hacker-green hover:text-hacker-ochre"
            }`
          }
        >
          {link.label}
        </NavLink>
      ))}
    </nav>
  );
}
