package memory

import (
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
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) error {
	i.store[name]++
	return nil
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, error) {
	return i.store[name], nil
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
