
-- name: CreatePlayer :one
INSERT INTO players (name) VALUES ($1)
RETURNING id, name;

-- name: GetPlayerByName :one
SELECT id, name FROM players WHERE name = $1;

-- name: GetPlayerByID :one
SELECT id, name FROM players WHERE id = $1;

-- name: ListPlayers :many
SELECT id, name FROM players
ORDER BY name;

-- name: RecordInteraction :exec
INSERT INTO interactions (protagonist_id, opponent_id, protagonist_move, opponent_move)
VALUES ($1, $2, $3, $4);

-- name: GetHistory :many
SELECT id, protagonist_id, opponent_id, protagonist_move, opponent_move, played_at
FROM interactions
ORDER BY played_at DESC;

