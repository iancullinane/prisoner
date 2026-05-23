package testhelpers

import (
	"reflect"
	"testing"

	"github.com/iancullinane/prisoner/internal/types"
)

func AssertLeague(t testing.TB, got, want []types.Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

// RunPlayerStoreContract runs a behavioural contract test suite against any
// PlayerStore implementation. newStore is called before each subtest so each
// case starts with a clean, empty store.
func RunPlayerStoreContract(t *testing.T, newStore func() types.PlayerStore) {
	t.Helper()

	t.Run("RecordWin increments score", func(t *testing.T) {
		store := newStore()
		store.RecordWin("Alice")
		store.RecordWin("Alice")
		got := store.GetPlayerScore("Alice")
		if got != 2 {
			t.Errorf("got score %d, want 2", got)
		}
	})

	t.Run("GetPlayerScore returns 0 for unknown player", func(t *testing.T) {
		store := newStore()
		got := store.GetPlayerScore("Nobody")
		if got != 0 {
			t.Errorf("got score %d, want 0", got)
		}
	})

	t.Run("RecordWin adds new player to league", func(t *testing.T) {
		store := newStore()
		store.RecordWin("Bob")
		league := store.GetLeague()
		var found *types.Player
		for i := range league {
			if league[i].Name == "Bob" {
				found = &league[i]
				break
			}
		}
		if found == nil {
			t.Fatalf("expected Bob in league, got %v", league)
		}
		if found.Wins != 1 {
			t.Errorf("got wins %d, want 1", found.Wins)
		}
	})

	t.Run("GetLeague returns all recorded players", func(t *testing.T) {
		store := newStore()
		store.RecordWin("Alice")
		store.RecordWin("Bob")
		league := store.GetLeague()
		if len(league) != 2 {
			t.Errorf("got %d players in league, want 2", len(league))
		}
	})
}
