package memory

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

// MARK: HistoryStore
// ------------------------------------------------------------

type InMemoryHistoryStore struct {
	logger *slog.Logger
	store  []types.Interaction
}

func NewInMemoryHistoryStore(logger *slog.Logger) *InMemoryHistoryStore {
	return &InMemoryHistoryStore{
		logger: logger.With(slog.String("component", "memory-history-store")),
		store:  []types.Interaction{},
	}
}

func (i *InMemoryHistoryStore) GetHistory() (types.History, error) {
	return i.store, nil
}

func (i *InMemoryHistoryStore) RecordInteraction(interaction types.Interaction) error {
	i.store = append(i.store, interaction)
	i.logger.Debug("recorded interaction",
		slog.String("id", interaction.ID.String()),
		slog.Int("history_len", len(i.store)),
	)
	return nil
}

func (i *InMemoryHistoryStore) GetHistoryByPlayerID(playerID uuid.UUID) (types.History, error) {
	historyByPlayer := make(types.History, 0)
	for _, interaction := range i.store {
		if interaction.PlayerA == playerID || interaction.PlayerB == playerID {
			historyByPlayer = append(historyByPlayer, interaction)
		}
	}
	return historyByPlayer, nil
}

// MARK: PlayerStore
// ------------------------------------------------------------

// in_memory_player_store.go
func NewInMemoryPlayerStore(logger *slog.Logger) *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		logger:  logger.With(slog.String("component", "memory-player-store")),
		players: types.Players{},
	}
}

type InMemoryPlayerStore struct {
	logger  *slog.Logger
	players types.Players
}

func (i *InMemoryPlayerStore) GetOrCreatePlayer(name string) (types.Player, error) {
	player := i.players.FindByName(name)
	if player != nil {
		return *player, nil
	}

	newPlayer := types.NewPlayer(name)
	i.players = append(i.players, *newPlayer)
	i.logger.Debug("created player",
		slog.String("name", newPlayer.Name),
		slog.String("id", newPlayer.ID.String()),
	)
	return *newPlayer, nil
}

func (i *InMemoryPlayerStore) GetPlayerByID(id uuid.UUID) (types.Player, error) {
	player := i.players.FindByID(id)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}
	return *player, nil
}

func (i *InMemoryPlayerStore) GetPlayerByName(name string) (types.Player, error) {
	player := i.players.FindByName(name)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}
	return *player, nil
}

func (i *InMemoryPlayerStore) GetAllPlayers() (types.Players, error) {
	return i.players, nil
}

func (i *InMemoryPlayerStore) GetRandomPlayer() (types.Player, error) {
	return i.players.GetRandomPlayer()
}

func (i *InMemoryPlayerStore) GetRandomPlayerExcept(exceptID uuid.UUID) (types.Player, error) {
	return i.players.GetRandomPlayerExcept(exceptID)
}
