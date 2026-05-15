package api

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

func NewPostgresMemoryStore() *PostgresMemoryStore {
	return &PostgresMemoryStore{map[string]int{}}
}

type PostgresMemoryStore struct {
	store map[string]int
}

func (p *PostgresMemoryStore) RecordWin(name string) {

}

func (p *PostgresMemoryStore) GetPlayerScore(name string) int {
	return p.store[name]
}
