-- +goose Up
CREATE TABLE interactions (
  id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  protagonist_id   UUID        NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  opponent_id      UUID        NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  protagonist_move TEXT        NOT NULL,
  opponent_move    TEXT        NOT NULL,
  played_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX interactions_protagonist_idx ON interactions (protagonist_id);
CREATE INDEX interactions_opponent_idx    ON interactions (opponent_id);

-- +goose Down
DROP TABLE IF EXISTS interactions;
