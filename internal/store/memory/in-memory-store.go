package memory

import (
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

// MARK: HistoryStore
// ------------------------------------------------------------

type InMemoryHistoryStore struct {
	store []types.Interaction
}

func NewInMemoryHistoryStore() *InMemoryHistoryStore {
	return &InMemoryHistoryStore{[]types.Interaction{}}
}

func (i *InMemoryHistoryStore) GetHistory() (types.History, error) {
	return i.store, nil
}

func (i *InMemoryHistoryStore) RecordInteraction(interaction types.Interaction) error {
	i.store = append(i.store, interaction)
	return nil
}

// MARK: PlayerStore
// ------------------------------------------------------------

// in_memory_player_store.go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		players: types.Players{},
	}
}

type InMemoryPlayerStore struct {
	players types.Players
}

func (i *InMemoryPlayerStore) GetOrCreatePlayer(name string) (types.Player, error) {
	player := i.players.FindByName(name)
	if player != nil {
		return *player, nil
	}

	newPlayer := types.NewPlayer(name)
	i.players = append(i.players, *newPlayer)
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
