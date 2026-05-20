package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iancullinane/prisoner/internal/types"
	sqlcdb "github.com/iancullinane/prisoner/prisonerdb"
)

type PlayerStore struct {
	pool *pgxpool.Pool
	q    *sqlcdb.Queries
}

func NewPlayerStore(pool *pgxpool.Pool) *PlayerStore {
	return &PlayerStore{pool: pool, q: sqlcdb.New(pool)}
}

func (s *PlayerStore) GetPlayerScore(name string) int {
	ctx := context.Background()
	row, err := s.q.GetPlayerByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0
		}
		panic(err)
	}
	return int(row.Wins)
}

func (s *PlayerStore) RecordWin(name string) {
	ctx := context.Background()
	if err := s.q.UpsertPlayerWin(ctx, name); err != nil {
		panic(err)
	}
}

func (s *PlayerStore) GetLeague() types.League {
	ctx := context.Background()
	rows, err := s.q.ListPlayers(ctx)
	if err != nil {
		panic(err)
	}
	out := make([]types.Player, 0, len(rows))
	for _, r := range rows {
		out = append(out, types.Player{Name: r.Name, Wins: int(r.Wins)})
	}
	return out
}
