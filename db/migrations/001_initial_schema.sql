-- +goose Up
CREATE TABLE IF NOT EXISTS players (
  id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  wins INT  NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS players;
