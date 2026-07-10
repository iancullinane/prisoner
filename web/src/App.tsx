import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Layout from "./components/Layout";
import PlayersPage from "./pages/PlayersPage";
import HistoryPage from "./pages/HistoryPage";
import PlayPage from "./pages/PlayPage";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Navigate to="/players" replace />} />
          <Route path="/players" element={<PlayersPage />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/play" element={<PlayPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
