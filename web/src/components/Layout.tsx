import { Outlet } from "react-router-dom";
import Sidebar from "./Sidebar";
import WindowFrame from "./WindowFrame";

export default function Layout() {
  return (
    <WindowFrame>
      <Sidebar />
      <main className="flex flex-1 flex-col overflow-y-auto p-6">
        <Outlet />
      </main>
    </WindowFrame>
  );
}
