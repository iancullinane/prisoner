import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import Sidebar from "./Sidebar";
import WindowFrame from "./WindowFrame";

const KEY_ROUTES: Record<string, string> = {
  "1": "/players",
  "2": "/history",
  "3": "/play",
};

export default function Layout() {
  const navigate = useNavigate();

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.metaKey || e.ctrlKey || e.altKey) return;
      const target = e.target as Element | null;
      if (target && /input|select|textarea/i.test(target.tagName)) return;
      const route = KEY_ROUTES[e.key];
      if (route) navigate(route);
    }
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, [navigate]);

  return (
    <WindowFrame>
      <Sidebar />
      <main className="flex flex-1 flex-col overflow-y-auto p-6">
        <Outlet />
      </main>
    </WindowFrame>
  );
}
