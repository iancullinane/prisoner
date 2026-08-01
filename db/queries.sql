
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



-- name: GetPrettyHistory :many
SELECT pa.name as player_a_name, pb.name as player_b_name, player_a_move, player_b_move, played_at
FROM interactions
INNER JOIN players pa ON interactions.player_a_id = pa.id
INNER JOIN players pb ON interactions.player_b_id = pb.id
WHERE sqlc.narg('player_id')::UUID IS NULL
    OR pa.id = sqlc.narg('player_id')
    OR pb.id = sqlc.narg('player_id')
ORDER BY interactions.played_at DESC;
