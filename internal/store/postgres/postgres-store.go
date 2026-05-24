package postgres

import (
	"context"
	"errors"
	"fmt"

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

func (s *PlayerStore) GetPlayerScore(name string) (int, error) {
	ctx := context.Background()
	row, err := s.q.GetPlayerByName(ctx, name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("get player %q: %w", name, err)
	}
	return int(row.Wins), nil
}

func (s *PlayerStore) RecordWin(name string) error {
	ctx := context.Background()
	if err := s.q.UpsertPlayerWin(ctx, name); err != nil {
		return fmt.Errorf("record win for %q: %w", name, err)
	}
	return nil
}

func (s *PlayerStore) GetLeague() (types.League, error) {
	ctx := context.Background()
	rows, err := s.q.ListPlayers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	out := make([]types.Player, 0, len(rows))
	for _, r := range rows {
		out = append(out, types.Player{Name: r.Name, Wins: int(r.Wins)})
	}
	return out, nil
}
