import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import WindowFrame from "./WindowFrame";

export default function Layout() {
  return (
    <WindowFrame>
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-6">
        <Outlet />
      </main>
    </WindowFrame>
  );
}
