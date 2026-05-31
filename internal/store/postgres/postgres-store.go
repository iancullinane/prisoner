package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
	sqlcdb "github.com/iancullinane/prisoner/prisonerdb"
)

type HistoryStore struct {
	pool *pgxpool.Pool
	q    *sqlcdb.Queries
}

func NewHistoryStore(pool *pgxpool.Pool) *HistoryStore {
	return &HistoryStore{pool: pool, q: sqlcdb.New(pool)}
}

func (s *HistoryStore) GetHistory() (types.History, error) {
	ctx := context.Background()
	rows, err := s.q.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	out := make(types.History, 0, len(rows))
	for _, r := range rows {
		out = append(out, types.Interaction{
			ID:              uuid.UUID(r.ID.Bytes),
			Protagonist:     uuid.UUID(r.ProtagonistID.Bytes),
			Opponent:        uuid.UUID(r.OpponentID.Bytes),
			ProtagonistMove: prisoner.Move([]rune(r.ProtagonistMove)[0]),
			OpponentMove:    prisoner.Move([]rune(r.OpponentMove)[0]),
		})
	}
	return out, nil
}

func (s *HistoryStore) RecordInteraction(interaction types.Interaction) error {
	ctx := context.Background()
	params := sqlcdb.RecordInteractionParams{
		ProtagonistID:   pgtype.UUID{Bytes: interaction.Protagonist, Valid: true},
		OpponentID:      pgtype.UUID{Bytes: interaction.Opponent, Valid: true},
		ProtagonistMove: string(interaction.ProtagonistMove),
		OpponentMove:    string(interaction.OpponentMove),
	}
	if err := s.q.RecordInteraction(ctx, params); err != nil {
		return err
	}
	return nil
}

type PlayerStore struct {
	pool *pgxpool.Pool
	q    *sqlcdb.Queries
}

func NewPlayerStore(pool *pgxpool.Pool) *PlayerStore {
	return &PlayerStore{pool: pool, q: sqlcdb.New(pool)}
}

func (s *PlayerStore) CreatePlayer(name string) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.CreatePlayer(ctx, name)
	if err != nil {
		return types.Player{}, fmt.Errorf("get player %q: %w", name, err)
	}

	return types.Player{
		ID:   uuid.UUID(player.ID.Bytes),
		Name: player.Name,
	}, nil
}

func (s *PlayerStore) GetPlayerByID(id uuid.UUID) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetPlayerByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return types.Player{}, fmt.Errorf("get player by id %q: %w", id, err)
	}

	return types.Player{
		ID:   uuid.UUID(player.ID.Bytes),
		Name: player.Name,
	}, nil
}

func (s *PlayerStore) GetPlayerByName(name string) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetPlayerByName(ctx, name)
	if err != nil {
		return types.Player{}, fmt.Errorf("get player by name %q: %w", name, err)
	}

	return types.Player{
		ID:   uuid.UUID(player.ID.Bytes),
		Name: player.Name,
	}, nil
}
