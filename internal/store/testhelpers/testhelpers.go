package testhelpers

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

// NoopLogger returns a logger that discards all output, for store constructors
// under test where log output is not the thing being asserted.
func NoopLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

const (
	TestPlayerOne     = "00000000-0000-0000-0000-000000000000"
	TestPlayerTwo     = "00000000-0000-0000-0000-111111111111"
	TestPlayerThree   = "00000000-0000-0000-0000-222222222222"
	TestPlayerOneName = "Alice"
	TestPlayerTwoName = "Bob"
)

var (
	TestPlayerOneID   = uuid.MustParse(TestPlayerOne)
	TestPlayerTwoID   = uuid.MustParse(TestPlayerTwo)
	TestPlayerThreeID = uuid.MustParse(TestPlayerThree)
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
func AssertPlayersContain(t testing.TB, got types.Players, want ...types.Player) {
	t.Helper()
	for _, player := range want {
		if got.FindByID(player.ID) == nil {
			t.Errorf("got %v, missing player %+v", got, player)
		}
	}
}

// MARK: Contracts

func RunHistoryStoreContract(t *testing.T, newStore func(t *testing.T) types.HistoryStore) {
	t.Helper()

	t.Run("record interaction", func(t *testing.T) {
		store := newStore(t)
		move1, move2 := prisoner.Cooperate, prisoner.Cooperate
		mustRecordInteraction(t, store, TestPlayerOneID, TestPlayerTwoID, move1, move2)
		history := mustGetHistory(t, store)
		if len(history) != 1 {
			t.Errorf("got %d interactions, want 1", len(history))
		}
		if history[0].PlayerA != TestPlayerOneID {
			t.Errorf("got player a %v, want %v", history[0].PlayerA, TestPlayerOneID)
		}
		if history[0].PlayerB != TestPlayerTwoID {
			t.Errorf("got player b %v, want %v", history[0].PlayerB, TestPlayerTwoID)
		}
		if history[0].PlayerAMove != move1 {
			t.Errorf("got player a move %v, want %v", history[0].PlayerAMove, move1)
		}
		if history[0].PlayerBMove != move2 {
			t.Errorf("got player b move %v, want %v", history[0].PlayerBMove, move2)
		}

	})

	t.Run("history of unknown player is empty list", func(t *testing.T) {
		store := newStore(t)
		mustRecordInteraction(t, store, TestPlayerOneID, TestPlayerTwoID, prisoner.Cooperate, prisoner.Cooperate)

		history := mustGetHistoryByPlayerID(t, store, TestPlayerThreeID)
		if len(history) != 0 {
			t.Errorf("got %d interactions, want 0", len(history))
		}

		history = mustGetHistoryByPlayerID(t, store, TestPlayerOneID)
		if len(history) != 1 {
			t.Errorf("got %d interactions, want 1", len(history))
		}
	})
}

// ========================================================

// RunPlayerStoreContract runs a behavioural contract test suite against any
// PlayerStore implementation. newStore is called before each subtest so each
// case starts with a clean, empty store.
func RunPlayerStoreContract(t *testing.T, newStore func(t *testing.T) types.PlayerStore) {
	t.Helper()

	t.Run("GetOrCreatePlayer creates a player", func(t *testing.T) {
		store := newStore(t)

		player := mustGetOrCreatePlayer(t, store, "Alice")
		want := types.Player{ID: player.ID, Name: "Alice"}
		AssertPlayer(t, player, want)
	})

	t.Run("GetOrCreatePlayer returns existing player", func(t *testing.T) {
		store := newStore(t)

		player1 := mustGetOrCreatePlayer(t, store, "Alice")
		player2 := mustGetOrCreatePlayer(t, store, "Alice")

		if player1.ID != player2.ID {
			t.Errorf("got different IDs %v and %v, want same ID", player1.ID, player2.ID)
		}
		AssertPlayer(t, player1, player2)
	})

	t.Run("GetPlayerByID returns existing player", func(t *testing.T) {
		store := newStore(t)

		created := mustGetOrCreatePlayer(t, store, TestPlayerOneName)
		got := mustGetPlayer(t, store, created.ID)
		AssertPlayer(t, got, created)
	})

	t.Run("GetPlayerByID returns error for unknown player", func(t *testing.T) {
		store := newStore(t)

		_, err := store.GetPlayerByID(TestPlayerThreeID)
		assertPlayerNotFound(t, err)
	})

	t.Run("GetPlayerByName returns existing player", func(t *testing.T) {
		store := newStore(t)

		created := mustGetOrCreatePlayer(t, store, TestPlayerOneName)
		got := mustGetPlayerByName(t, store, TestPlayerOneName)
		AssertPlayer(t, got, created)
	})

	t.Run("GetPlayerByName returns error for unknown player", func(t *testing.T) {
		store := newStore(t)

		_, err := store.GetPlayerByName("Nobody")
		assertPlayerNotFound(t, err)
	})

	t.Run("GetAllPlayers returns created players", func(t *testing.T) {
		store := newStore(t)

		alice := mustGetOrCreatePlayer(t, store, TestPlayerOneName)
		bob := mustGetOrCreatePlayer(t, store, TestPlayerTwoName)
		players, err := store.GetAllPlayers()
		AssertNoError(t, err)
		AssertPlayersContain(t, players, alice, bob)
	})
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

func mustGetPlayerByName(t *testing.T, store types.PlayerStore, name string) types.Player {
	t.Helper()
	player, err := store.GetPlayerByName(name)
	if err != nil {
		t.Fatalf("GetPlayerByName(%v) returned error: %v", name, err)
	}
	return player
}

func assertPlayerNotFound(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected player not found error, got nil")
	}
}

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

func mustGetHistoryByPlayerID(t *testing.T, store types.HistoryStore, playerID uuid.UUID) types.History {
	t.Helper()
	history, err := store.GetHistoryByPlayerID(playerID)
	if err != nil {
		t.Fatalf("GetHistoryByPlayerID(%v) returned error: %v", playerID, err)
	}
	return history
}
