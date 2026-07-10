import type { Interaction } from "../types";

const BASE_URL = import.meta.env.VITE_API_BASE_URL;

export async function listHistory(): Promise<Interaction[]> {
  const response = await fetch(`${BASE_URL}/api/v1/history`);
  if (!response.ok) {
    throw new Error(`failed to list history: ${response.status}`);
  }
  return response.json();
}

export async function listHistoryForPlayer(id: string): Promise<Interaction[]> {
  const response = await fetch(`${BASE_URL}/api/v1/history/${encodeURIComponent(id)}`);
  if (!response.ok) {
    throw new Error(`failed to list history for player: ${response.status}`);
  }
  return response.json();
}
