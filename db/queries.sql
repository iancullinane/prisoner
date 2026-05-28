-- name: GetPlayerByName :one
SELECT id, name, wins FROM players WHERE name = $1;

-- name: ListPlayers :many
SELECT id, name, wins FROM players
ORDER BY name;

-- name: UpsertPlayerWin :exec
INSERT INTO players (name, wins) VALUES ($1, 1)
ON CONFLICT (name) DO UPDATE SET wins = players.wins + 1;

-- name: RecordInteraction :exec
INSERT INTO interactions (protagonist_id, opponent_id, protagonist_move, opponent_move)
VALUES ($1, $2, $3, $4);

-- name: GetHistory :many
SELECT id, protagonist_id, opponent_id, protagonist_move, opponent_move, played_at
FROM interactions
ORDER BY played_at DESC;
