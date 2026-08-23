import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import MoveChip from "./MoveChip";

describe("MoveChip", () => {
  it("renders glyph and letter for cooperate", () => {
    render(<MoveChip move="C" />);
    expect(screen.getByText("✓ C")).toBeInTheDocument();
    expect(screen.getByLabelText("cooperate")).toBeInTheDocument();
  });

  it("renders glyph and letter for betray", () => {
    render(<MoveChip move="B" />);
    expect(screen.getByText("✕ B")).toBeInTheDocument();
    expect(screen.getByLabelText("betray")).toBeInTheDocument();
  });
});
