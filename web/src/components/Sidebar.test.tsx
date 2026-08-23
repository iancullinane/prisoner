import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import Sidebar from "./Sidebar";

function renderSidebar(path = "/players") {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Sidebar />
    </MemoryRouter>,
  );
}

describe("Sidebar", () => {
  it("renders links for players, history, and play with key hints", () => {
    renderSidebar();

    expect(screen.getByRole("link", { name: /players/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /history/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /play$/i })).toBeInTheDocument();

    expect(screen.getByText("[1]")).toBeInTheDocument();
    expect(screen.getByText("[2]")).toBeInTheDocument();
    expect(screen.getByText("[3]")).toBeInTheDocument();
  });

  it("marks the active route with aria-current and leaves the rest unmarked", () => {
    renderSidebar("/history");

    expect(screen.getByRole("link", { name: /history/i })).toHaveAttribute(
      "aria-current",
      "page",
    );
    expect(screen.getByRole("link", { name: /players/i })).not.toHaveAttribute(
      "aria-current",
    );
  });

  it("shows the move color legend", () => {
    renderSidebar();

    expect(screen.getByText("cooperate")).toBeInTheDocument();
    expect(screen.getByText("betray")).toBeInTheDocument();
  });
});
