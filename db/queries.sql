-- name: GetPlayerByName :one
SELECT id, name, wins FROM players WHERE name = $1;

-- name: ListPlayers :many
SELECT id, name, wins FROM players
ORDER BY name;

-- name: UpsertPlayerWin :exec
INSERT INTO players (name, wins) VALUES ($1, 1)
ON CONFLICT (name) DO UPDATE SET wins = players.wins + 1;
