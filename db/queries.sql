
-- name: GetOrCreatePlayer :one
INSERT INTO players (name) VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING id, name, created_at;

-- name: GetPlayerByName :one
SELECT id, name, created_at FROM players WHERE name = $1;

-- name: GetPlayerByID :one
SELECT id, name, created_at FROM players WHERE id = $1;

-- name: ListPlayers :many
SELECT id, name, created_at FROM players
ORDER BY name;

-- name: GetRandomPlayer :one
SELECT id, name, created_at FROM players ORDER BY RANDOM() LIMIT 1;

-- name: GetRandomPlayerExcept :one
SELECT id, name, created_at FROM players WHERE id != $1 ORDER BY RANDOM() LIMIT 1;

-- name: RecordInteraction :exec
INSERT INTO interactions (player_a_id, player_b_id, player_a_move, player_b_move, played_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetHistory :many
SELECT id, player_a_id, player_b_id, player_a_move, player_b_move, played_at
FROM interactions
ORDER BY played_at DESC;

-- name: GetHistoryByPlayerID :many
SELECT id, player_a_id, player_b_id, player_a_move, player_b_move, played_at
FROM interactions
WHERE player_a_id = $1 OR player_b_id = $1
ORDER BY played_at DESC;
