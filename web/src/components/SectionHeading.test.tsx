import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import SectionHeading from "./SectionHeading";

describe("SectionHeading", () => {
  it("renders a level-one heading with its text", () => {
    render(<SectionHeading>History</SectionHeading>);
    expect(
      screen.getByRole("heading", { level: 1, name: "History" }),
    ).toBeInTheDocument();
  });
});
