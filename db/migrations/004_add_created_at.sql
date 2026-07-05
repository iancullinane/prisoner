-- +goose Up
ALTER TABLE players ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE players DROP COLUMN created_at;
