CREATE TABLE IF NOT EXISTS players (
  id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name text      NOT NULL UNIQUE,
  wins int     NOT NULL DEFAULT 0
);
