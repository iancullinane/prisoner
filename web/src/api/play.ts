import type { PlayRequest, PlayResponse } from "../types";

const BASE_URL = import.meta.env.VITE_API_BASE_URL;

export async function playRound(req: PlayRequest): Promise<PlayResponse> {
  const response = await fetch(`${BASE_URL}/api/v1/play`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(req),
  });
  if (!response.ok) {
    throw new Error(`failed to play round: ${response.status}`);
  }
  return response.json();
}
