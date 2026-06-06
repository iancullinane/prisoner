-- +goose Up
CREATE TABLE interactions (
  id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  player_a_id   UUID        NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  player_b_id   UUID        NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  player_a_move TEXT        NOT NULL,
  player_b_move TEXT        NOT NULL,
  played_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX interactions_player_a_idx ON interactions (player_a_id);
CREATE INDEX interactions_player_b_idx ON interactions (player_b_id);

-- +goose Down
DROP TABLE IF EXISTS interactions;
