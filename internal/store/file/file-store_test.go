package file

import (
	"testing"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

// testID1, _ = uuid.Parse("00000000-0000-aaaa-2222-222222222222")
// testID2, _ = uuid.Parse("11111111-1111-bbbb-3333-333333333333")
var (
	player1, _          = uuid.Parse("00000000-0000-aaaa-2222-222222222222")
	player2, _          = uuid.Parse("11111111-1111-bbbb-3333-333333333333")
	interactionIDOne, _ = uuid.Parse("22222222-2222-cccc-4444-444444444444")
)

// =========================

const goldenPlayerStoreData = `[
		{"ID": "00000000-0000-aaaa-2222-222222222222", "Name": "Chris"},
		{"ID": "11111111-1111-bbbb-3333-333333333333", "Name": "Cleo"}
	]`

func TestFileSystemPlayerStore(t *testing.T) {
	newTempStore, closer := createTempFile(t, goldenPlayerStoreData)
	defer closer()

	t.Run("get a player", func(t *testing.T) {
		store, err := NewFileSystemPlayerStore(testhelpers.NoopLogger(), newTempStore)
		testhelpers.AssertNoError(t, err)

		got, err := store.GetPlayerByID(player1)
		testhelpers.AssertNoError(t, err)
		want := types.Player{ID: player1, Name: "Chris"}

		testhelpers.AssertPlayer(t, got, want)

	})

	// file_system_store_test.go
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(testhelpers.NoopLogger(), database)

		testhelpers.AssertNoError(t, err)
	})

	t.Run("persists a created player across store instances", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "[]")
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(testhelpers.NoopLogger(), database)
		testhelpers.AssertNoError(t, err)

		created, err := store.GetOrCreatePlayer("Alice")
		testhelpers.AssertNoError(t, err)

		reopened, err := NewFileSystemPlayerStore(testhelpers.NoopLogger(), database)
		testhelpers.AssertNoError(t, err)

		got, err := reopened.GetPlayerByName("Alice")
		testhelpers.AssertNoError(t, err)
		testhelpers.AssertPlayer(t, got, created)
	})

}

// MARK: History Store Tests

const goldenHistoryData = `[
	{
		"ID": "22222222-2222-cccc-4444-444444444444",
		"PlayerA": "00000000-0000-aaaa-2222-222222222222",
		"PlayerB": "11111111-1111-bbbb-3333-333333333333",
		"PlayerAMove": "C",
		"PlayerBMove": "B"
	}
]`

func TestFileSystemHistoryStore(t *testing.T) {
	newTempStore, closer := createTempFile(t, goldenHistoryData)
	defer closer()
	t.Run("single interaction", func(t *testing.T) {
		store, err := NewFileSystemHistoryStore(testhelpers.NoopLogger(), newTempStore)
		testhelpers.AssertNoError(t, err)

		got, err := store.GetInteraction(interactionIDOne)
		// testhelpers.AssertNoError(t, err)

		want := types.Interaction{
			ID:          interactionIDOne,
			PlayerA:     player1,
			PlayerB:     player2,
			PlayerAMove: prisoner.Cooperate,
			PlayerBMove: prisoner.Betray,
		}

		testhelpers.AssertInteraction(t, got, want)

	})

	// file_system_store_test.go
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(testhelpers.NoopLogger(), database)

		testhelpers.AssertNoError(t, err)
	})

}
