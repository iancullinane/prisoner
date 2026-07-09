export interface Player {
  id: string;
  name: string;
  CreatedAt: string;
}

export type Move = "C" | "B";
export type Result = "T" | "R" | "P" | "S";

export interface Interaction {
  id: string;
  playerA: string;
  playerB: string;
  playerAMove: Move;
  playerBMove: Move;
  playedAt: string;
}

export interface PlayRequest {
  playerA: string;
  playerB: string;
  playerAMove: Move;
  playerBMove: Move;
}

export interface PlayResponse {
  id: string;
  playerAScore: Result;
  playerBScore: Result;
}
