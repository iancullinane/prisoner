import type { Player } from "../types";

const BASE_URL = import.meta.env.VITE_API_BASE_URL;

export async function listPlayers(): Promise<Player[]> {
  const response = await fetch(`${BASE_URL}/players`);
  if (!response.ok) {
    throw new Error(`failed to list players: ${response.status}`);
  }
  return response.json();
}

export async function createPlayer(name: string): Promise<Player> {
  const response = await fetch(`${BASE_URL}/players/${encodeURIComponent(name)}`, {
    method: "POST",
  });
  if (!response.ok) {
    throw new Error(`failed to create player: ${response.status}`);
  }
  return response.json();
}
