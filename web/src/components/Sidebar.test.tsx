import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import Sidebar from "./Sidebar";

describe("Sidebar", () => {
  it("renders links for players, history, and play", () => {
    render(
      <MemoryRouter initialEntries={["/players"]}>
        <Sidebar />
      </MemoryRouter>,
    );

    expect(screen.getByRole("link", { name: /players/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /history/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /^play$/i })).toBeInTheDocument();
  });

  it("marks the active route with aria-current and leaves the rest unmarked", () => {
    render(
      <MemoryRouter initialEntries={["/history"]}>
        <Sidebar />
      </MemoryRouter>,
    );

    expect(screen.getByRole("link", { name: /history/i })).toHaveAttribute(
      "aria-current",
      "page",
    );
    expect(screen.getByRole("link", { name: /players/i })).not.toHaveAttribute(
      "aria-current",
    );
  });
});
