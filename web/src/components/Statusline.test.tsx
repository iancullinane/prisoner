import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import Statusline from "./Statusline";

describe("Statusline", () => {
  it("shows the view-switch key hints", () => {
    render(<Statusline />);
    expect(screen.getByText("switch view")).toBeInTheDocument();
    expect(screen.getByText("1")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
    expect(screen.getByText("3")).toBeInTheDocument();
  });

  it("shows the mode indicator", () => {
    render(<Statusline />);
    expect(screen.getByText("— NORMAL —")).toBeInTheDocument();
  });
});
