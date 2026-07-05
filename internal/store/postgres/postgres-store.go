package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
	sqlcdb "github.com/iancullinane/prisoner/prisonerdb"
)

type HistoryStore struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	q      *sqlcdb.Queries
}

func NewHistoryStore(logger *slog.Logger, pool *pgxpool.Pool) *HistoryStore {
	return &HistoryStore{
		logger: logger.With(slog.String("component", "postgres-history-store")),
		pool:   pool,
		q:      sqlcdb.New(pool),
	}
}

// interaction from row
func interactionFromRow(i sqlcdb.Interaction) types.Interaction {
	return types.Interaction{
		ID:          i.ID.Bytes,
		PlayerA:     i.PlayerAID.Bytes,
		PlayerB:     i.PlayerBID.Bytes,
		PlayerAMove: prisoner.Move([]rune(i.PlayerAMove)[0]),
		PlayerBMove: prisoner.Move([]rune(i.PlayerBMove)[0]),
		PlayedAt:    i.PlayedAt.Time,
	}
}

func (s *HistoryStore) GetHistory() (types.History, error) {
	ctx := context.Background()
	rows, err := s.q.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	out := make(types.History, 0, len(rows))
	for _, r := range rows {
		out = append(out, interactionFromRow(r))
	}
	return out, nil
}

func (s *HistoryStore) RecordInteraction(interaction types.Interaction) error {
	ctx := context.Background()
	params := sqlcdb.RecordInteractionParams{
		PlayerAID:   pgtype.UUID{Bytes: interaction.PlayerA, Valid: true},
		PlayerBID:   pgtype.UUID{Bytes: interaction.PlayerB, Valid: true},
		PlayerAMove: string(interaction.PlayerAMove),
		PlayerBMove: string(interaction.PlayerBMove),
		PlayedAt:    pgtype.Timestamptz{Time: interaction.PlayedAt, Valid: true},
	}
	if err := s.q.RecordInteraction(ctx, params); err != nil {
		return err
	}
	s.logger.Debug("recorded interaction",
		slog.String("id", interaction.ID.String()),
	)
	return nil
}

func (s *HistoryStore) GetHistoryByPlayerID(playerID uuid.UUID) (types.History, error) {
	ctx := context.Background()
	playerHistory, err := s.q.GetHistoryByPlayerID(ctx, pgtype.UUID{Bytes: playerID, Valid: true})
	if err != nil {
		return types.History{}, err
	}

	out := make(types.History, 0, len(playerHistory))
	for _, r := range playerHistory {
		out = append(out, interactionFromRow(r))
	}
	return out, nil
}

type PlayerStore struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	q      *sqlcdb.Queries
}

func NewPlayerStore(logger *slog.Logger, pool *pgxpool.Pool) *PlayerStore {
	return &PlayerStore{
		logger: logger.With(slog.String("component", "postgres-player-store")),
		pool:   pool,
		q:      sqlcdb.New(pool),
	}
}

func playerFromRow(p sqlcdb.Player) types.Player {
	return types.Player{
		ID:        uuid.UUID(p.ID.Bytes),
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Time,
	}
}

func (s *PlayerStore) GetOrCreatePlayer(name string) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetOrCreatePlayer(ctx, name)
	if err != nil {
		return types.Player{}, fmt.Errorf("get or create player %q: %w", name, err)
	}

	out := playerFromRow(player)
	s.logger.Debug("get or create player",
		slog.String("name", out.Name),
		slog.String("id", out.ID.String()),
	)
	return out, nil
}

func (s *PlayerStore) GetPlayerByID(id uuid.UUID) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetPlayerByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Player{}, types.ErrPlayerNotFound
	}
	if err != nil {
		return types.Player{}, fmt.Errorf("get player by id %q: %w", id, err)
	}

	return playerFromRow(player), nil
}

func (s *PlayerStore) GetPlayerByName(name string) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetPlayerByName(ctx, name)
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Player{}, types.ErrPlayerNotFound
	}
	if err != nil {
		return types.Player{}, fmt.Errorf("get player by name %q: %w", name, err)
	}

	return playerFromRow(player), nil
}

func (s *PlayerStore) GetAllPlayers() (types.Players, error) {
	ctx := context.Background()
	players, err := s.q.ListPlayers(ctx)
	if err != nil {
		return types.Players{}, fmt.Errorf("could not list players %w", err)
	}

	out := make(types.Players, 0, len(players))
	for _, p := range players {
		out = append(out, playerFromRow(p))
	}
	return out, nil
}

func (s *PlayerStore) GetRandomPlayer() (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetRandomPlayer(ctx)
	if err != nil {
		return types.Player{}, fmt.Errorf("get random player: %w", err)
	}
	return playerFromRow(player), nil
}

func (s *PlayerStore) GetRandomPlayerExcept(exceptID uuid.UUID) (types.Player, error) {
	ctx := context.Background()
	player, err := s.q.GetRandomPlayerExcept(ctx, pgtype.UUID{Bytes: exceptID, Valid: true})
	if err != nil {
		return types.Player{}, fmt.Errorf("get random player except %q: %w", exceptID, err)
	}
	return playerFromRow(player), nil
}
