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
		mustRecordWin(t, store, "Alice")
		mustRecordWin(t, store, "Alice")
		got := mustGetPlayerScore(t, store, "Alice")
		if got != 2 {
			t.Errorf("got score %d, want 2", got)
		}
	})

	t.Run("GetPlayerScore returns 0 for unknown player", func(t *testing.T) {
		store := newStore()
		got := mustGetPlayerScore(t, store, "Nobody")
		if got != 0 {
			t.Errorf("got score %d, want 0", got)
		}
	})

	t.Run("RecordWin adds new player to league", func(t *testing.T) {
		store := newStore()
		mustRecordWin(t, store, "Bob")
		league := mustGetLeague(t, store)
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
		mustRecordWin(t, store, "Alice")
		mustRecordWin(t, store, "Bob")
		league := mustGetLeague(t, store)
		if len(league) != 2 {
			t.Errorf("got %d players in league, want 2", len(league))
		}
	})
}

func mustRecordWin(t *testing.T, store types.PlayerStore, name string) {
	t.Helper()
	if err := store.RecordWin(name); err != nil {
		t.Fatalf("RecordWin(%q) returned error: %v", name, err)
	}
}

func mustGetPlayerScore(t *testing.T, store types.PlayerStore, name string) int {
	t.Helper()
	score, err := store.GetPlayerScore(name)
	if err != nil {
		t.Fatalf("GetPlayerScore(%q) returned error: %v", name, err)
	}
	return score
}

func mustGetLeague(t *testing.T, store types.PlayerStore) types.League {
	t.Helper()
	league, err := store.GetLeague()
	if err != nil {
		t.Fatalf("GetLeague() returned error: %v", err)
	}
	return league
}
