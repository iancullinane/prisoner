package memory

import "github.com/iancullinane/prisoner/internal/types"

// in_memory_player_store.go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

// in_memory_player_store.go
func (i *InMemoryPlayerStore) GetLeague() types.League {
	var league []types.Player
	for name, wins := range i.store {
		league = append(league, types.Player{Name: name, Wins: wins})
	}
	return league
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
