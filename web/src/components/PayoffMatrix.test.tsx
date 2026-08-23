import { describe, it, expect } from "vitest";
import { render, screen, within } from "@testing-library/react";
import PayoffMatrix from "./PayoffMatrix";

describe("PayoffMatrix", () => {
  it("renders all four payoff cells with their captions", () => {
    render(<PayoffMatrix />);

    const table = screen.getByRole("table", { name: /payoff/i });
    expect(within(table).getByText("3 · 3")).toBeInTheDocument();
    expect(within(table).getByText("0 · 5")).toBeInTheDocument();
    expect(within(table).getByText("5 · 0")).toBeInTheDocument();
    expect(within(table).getByText("1 · 1")).toBeInTheDocument();
    expect(within(table).getByText("mutual trust")).toBeInTheDocument();
    expect(within(table).getByText("you're played")).toBeInTheDocument();
    expect(within(table).getByText("you cash in")).toBeInTheDocument();
    expect(within(table).getByText("mutual ruin")).toBeInTheDocument();
  });
});
