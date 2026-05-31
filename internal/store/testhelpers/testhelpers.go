package testhelpers

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

const (
	TestPlayerOne = "00000000-0000-aaaa-2222-222222222222"
	TestPlayerTwo = "11111111-1111-bbbb-3333-333333333333"
)

var (
	player1, _ = uuid.Parse(TestPlayerOne)
	player2, _ = uuid.Parse(TestPlayerTwo)
)

// MARK: Assertations

func AssertPlayer(t testing.TB, got, want types.Player) {
	t.Helper()

	if got.ID == uuid.Nil {
		t.Errorf("got nil ID, want valid UUID")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func AssertInteraction(t testing.TB, got, want types.Interaction) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't want one %v", err)
	}
}

func AssertLeague(t testing.TB, got, want []types.Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

// MARK: Contracts

func RunHistoryStoreContract(t *testing.T, newStore func() types.HistoryStore) {
	t.Helper()

	// a, b := mustTwoUUIDs()

	t.Run("record interaction", func(t *testing.T) {
		store := newStore()
		player1, player2 := mustTwoUUIDs()
		move1, move2 := prisoner.Cooperate, prisoner.Cooperate
		mustRecordInteraction(t, store, player1, player2, move1, move2)
		history := mustGetHistory(t, store)
		if len(history) != 1 {
			t.Errorf("got %d interactions, want 1", len(history))
		}
		if history[0].Protagonist != player1 {
			t.Errorf("got protagonist %v, want %v", history[0].Protagonist, player1)
		}
		if history[0].Opponent != player2 {
			t.Errorf("got opponent %v, want %v", history[0].Opponent, player2)
		}
		if history[0].ProtagonistMove != move1 {
			t.Errorf("got protagonist move %v, want %v", history[0].ProtagonistMove, move1)
		}
		if history[0].OpponentMove != move2 {
			t.Errorf("got opponent move %v, want %v", history[0].OpponentMove, move2)
		}

		// got := mustGetPlayerScore(t, store, "Alice")
		// if got != 2 {
		// 	t.Errorf("got score %d, want 2", got)
		// }
	})
}

// ========================================================

// RunPlayerStoreContract runs a behavioural contract test suite against any
// PlayerStore implementation. newStore is called before each subtest so each
// case starts with a clean, empty store.
func RunPlayerStoreContract(t *testing.T, newStore func() types.PlayerStore) {
	t.Helper()

	t.Run("GetOrCreatePlayer creates a player", func(t *testing.T) {
		store := newStore()

		player := mustGetOrCreatePlayer(t, store, "Alice")
		want := types.Player{ID: player.ID, Name: "Alice"}
		AssertPlayer(t, player, want)
	})

	t.Run("GetOrCreatePlayer returns existing player", func(t *testing.T) {
		store := newStore()

		player1 := mustGetOrCreatePlayer(t, store, "Alice")
		player2 := mustGetOrCreatePlayer(t, store, "Alice")
		
		if player1.ID != player2.ID {
			t.Errorf("got different IDs %v and %v, want same ID", player1.ID, player2.ID)
		}
		AssertPlayer(t, player1, player2)
	})

	// t.Run("GetPlayerScore returns 0 for unknown player", func(t *testing.T) {
	// 	store := newStore()
	// 	got := mustGetPlayerScore(t, store, "Nobody")
	// 	if got != 0 {
	// 		t.Errorf("got score %d, want 0", got)
	// 	}
	// })

	// t.Run("RecordWin adds new player to league", func(t *testing.T) {
	// 	store := newStore()
	// 	mustRecordWin(t, store, "Bob")
	// 	league := mustGetLeague(t, store)
	// 	var found *types.Player
	// 	for i := range league {
	// 		if league[i].Name == "Bob" {
	// 			found = &league[i]
	// 			break
	// 		}
	// 	}
	// 	if found == nil {
	// 		t.Fatalf("expected Bob in league, got %v", league)
	// 	}
	// 	if found.Wins != 1 {
	// 		t.Errorf("got wins %d, want 1", found.Wins)
	// 	}
	// })

	// t.Run("GetLeague returns all recorded players", func(t *testing.T) {
	// 	store := newStore()
	// 	mustRecordWin(t, store, "Alice")
	// 	mustRecordWin(t, store, "Bob")
	// 	league := mustGetLeague(t, store)
	// 	if len(league) != 2 {
	// 		t.Errorf("got %d players in league, want 2", len(league))
	// 	}
	// })
}

// MARK: Musts

func mustGetOrCreatePlayer(t *testing.T, store types.PlayerStore, name string) types.Player {
	t.Helper()
	player, err := store.GetOrCreatePlayer(name)
	if err != nil {
		t.Fatalf("GetOrCreatePlayer(%v) returned error: %v", name, err)
	}
	return player
}

func mustGetPlayer(t *testing.T, store types.PlayerStore, id uuid.UUID) types.Player {
	t.Helper()
	player, err := store.GetPlayerByID(id)
	if err != nil {
		t.Fatalf("GetPlayerByID(%v) returned error: %v", id, err)
	}
	return player
}

// func mustGetLeague(t *testing.T, store types.PlayerStore) types.League {
// 	t.Helper()
// 	league, err := store.GetLeague()
// 	if err != nil {
// 		t.Fatalf("GetLeague() returned error: %v", err)
// 	}
// 	return league
// }

func mustRecordInteraction(t *testing.T, store types.HistoryStore, player1, player2 uuid.UUID, move1, move2 prisoner.Move) {
	t.Helper()
	interaction := types.NewInteraction(player1, player2, move1, move2)
	if err := store.RecordInteraction(interaction); err != nil {
		t.Fatalf("RecordInteraction(%v) returned error: %v", interaction, err)
	}
}

func mustGetHistory(t *testing.T, store types.HistoryStore) types.History {
	t.Helper()
	history, err := store.GetHistory()
	if err != nil {
		t.Fatalf("GetHistory() returned error: %v", err)
	}
	return history
}

func mustTwoUUIDs() (a, b uuid.UUID) {
	a, _ = uuid.Parse(TestPlayerOne)
	b, _ = uuid.Parse(TestPlayerTwo)
	return
}
