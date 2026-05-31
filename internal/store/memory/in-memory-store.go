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

func (i *InMemoryPlayerStore) RecordWin(name string) error {
	return nil
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	return 1, nil
}

func (i *InMemoryPlayerStore) CreatePlayer(name string) (types.Player, error) {
	player := types.NewPlayer(name)
	i.players = append(i.players, *player)
	return *player, nil
}

func (i *InMemoryPlayerStore) GetPlayerByID(id uuid.UUID) (types.Player, error) {
	player := i.players.FindByID(id)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}
	return *player, nil
}

// MARK: Postgres
// -------------------------------

// func NewPostgresMemoryStore() *PostgresMemoryStore {
// 	return &PostgresMemoryStore{map[string]int{}}
// }

// type PostgresMemoryStore struct {
// 	store map[string]int
// }

// func (p *PostgresMemoryStore) RecordWin(name string) {
// 	p.store[name]++
// }

// func (p *PostgresMemoryStore) GetPlayerScore(name string) int {
// 	return p.store[name]
// }
