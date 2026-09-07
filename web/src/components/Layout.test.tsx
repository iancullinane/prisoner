import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import Layout from "./Layout";

vi.mock("../api/players", () => ({
  listPlayers: vi.fn().mockResolvedValue([]),
}));

vi.mock("../api/history", () => ({
  listHistory: vi.fn().mockResolvedValue([]),
}));

function renderLayout(path = "/players") {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route element={<Layout />}>
          <Route
            path="/players"
            element={
              <div>
                <p>players page</p>
                <input aria-label="scratch input" />
              </div>
            }
          />
          <Route path="/history" element={<p>history page</p>} />
          <Route path="/play" element={<p>play page</p>} />
        </Route>
      </Routes>
    </MemoryRouter>,
  );
}

describe("Layout keyboard navigation", () => {
  it("switches views with the number keys", async () => {
    const user = userEvent.setup();
    renderLayout();

    await user.keyboard("2");
    expect(await screen.findByText("history page")).toBeInTheDocument();

    await user.keyboard("3");
    expect(await screen.findByText("play page")).toBeInTheDocument();

    await user.keyboard("1");
    expect(await screen.findByText("players page")).toBeInTheDocument();
  });

  it("does not navigate while typing in a form field", async () => {
    const user = userEvent.setup();
    renderLayout();

    await user.click(screen.getByLabelText("scratch input"));
    await user.keyboard("2");

    expect(screen.getByText("players page")).toBeInTheDocument();
    expect(screen.queryByText("history page")).not.toBeInTheDocument();
    expect(screen.getByLabelText("scratch input")).toHaveValue("2");
  });
});
