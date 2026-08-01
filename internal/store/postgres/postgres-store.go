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
	sqlcdb "github.com/iancullinane/prisoner/prisonerdb"
)

type HistoryStore struct {
	logger  *slog.Logger
	pool    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewHistoryStore(logger *slog.Logger, pool *pgxpool.Pool) *HistoryStore {
	return &HistoryStore{
		logger:  logger.With(slog.String("component", "postgres-history-store")),
		pool:    pool,
		queries: sqlcdb.New(pool),
	}
}

func (s *HistoryStore) GetHistory() (types.History, error) {
	ctx := context.Background()
	rows, err := s.queries.GetHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	out := make(types.History, 0, len(rows))
	for _, r := range rows {
		out = append(out, interactionFromRow(r))
	}
	return out, nil
}

func (s *HistoryStore) GetPrettyHistory(playerID *uuid.UUID) (types.PrettyHistory, error) {
	ctx := context.Background()
	var arg uuid.NullUUID
	if playerID != nil {
		arg = uuid.NullUUID{UUID: *playerID, Valid: true}
	}
	rows, err := s.queries.GetPrettyHistory(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("get pretty history: %w", err)
	}
	out := make(types.PrettyHistory, 0, len(rows))
	for _, r := range rows {
		out = append(out, prettyFromRow(r))
	}
	return out, nil
}

func (s *HistoryStore) RecordInteraction(interaction types.Interaction) error {
	ctx := context.Background()
	params := sqlcdb.RecordInteractionParams{
		PlayerAID:   interaction.PlayerA,
		PlayerBID:   interaction.PlayerB,
		PlayerAMove: string(interaction.PlayerAMove),
		PlayerBMove: string(interaction.PlayerBMove),
		PlayedAt:    pgtype.Timestamptz{Time: interaction.PlayedAt, Valid: true},
	}
	if err := s.queries.RecordInteraction(ctx, params); err != nil {
		return err
	}
	s.logger.Debug("recorded interaction",
		slog.String("id", interaction.ID.String()),
	)
	return nil
}

func (s *HistoryStore) GetHistoryByPlayerID(playerID uuid.UUID) (types.History, error) {
	ctx := context.Background()
	playerHistory, err := s.queries.GetHistoryByPlayerID(ctx, playerID)
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
	logger  *slog.Logger
	pool    *pgxpool.Pool
	queries *sqlcdb.Queries
}

func NewPlayerStore(logger *slog.Logger, pool *pgxpool.Pool) *PlayerStore {
	return &PlayerStore{
		logger:  logger.With(slog.String("component", "postgres-player-store")),
		pool:    pool,
		queries: sqlcdb.New(pool),
	}
}

func playerFromRow(p sqlcdb.Player) types.Player {
	return types.Player{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Time,
	}
}

func (s *PlayerStore) GetOrCreatePlayer(name string) (types.Player, error) {
	ctx := context.Background()
	player, err := s.queries.GetOrCreatePlayer(ctx, name)
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
	player, err := s.queries.GetPlayerByID(ctx, id)
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
	player, err := s.queries.GetPlayerByName(ctx, name)
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
	players, err := s.queries.ListPlayers(ctx)
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
	player, err := s.queries.GetRandomPlayer(ctx)
	if err != nil {
		return types.Player{}, fmt.Errorf("get random player: %w", err)
	}
	return playerFromRow(player), nil
}

func (s *PlayerStore) GetRandomPlayerExcept(exceptID uuid.UUID) (types.Player, error) {
	ctx := context.Background()
	player, err := s.queries.GetRandomPlayerExcept(ctx, exceptID)
	if err != nil {
		return types.Player{}, fmt.Errorf("get random player except %q: %w", exceptID, err)
	}
	return playerFromRow(player), nil
}
