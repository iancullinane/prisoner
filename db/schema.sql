CREATE TABLE players (
  id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name text      NOT NULL UNIQUE,
  wins int32     NOT NULL DEFAULT 0
);
