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

	t.Run("get a placer", func(t *testing.T) {
		store, err := NewFileSystemPlayerStore(newTempStore)
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

		_, err := NewFileSystemPlayerStore(database)

		testhelpers.AssertNoError(t, err)
	})

}

func TestFileSystemPlayerStore_Contract(t *testing.T) {
	testhelpers.RunPlayerStoreContract(t, func() types.PlayerStore {
		f, cleanup := createTempFile(t, "[]")
		t.Cleanup(cleanup)
		store, err := NewFileSystemPlayerStore(f)
		if err != nil {
			t.Fatalf("could not create store: %v", err)
		}
		return store
	})
}

// MARK: History Store Tests

const goldenHistoryData = `[
	{
		"ID": "22222222-2222-cccc-4444-444444444444",
		"Protagonist": "00000000-0000-aaaa-2222-222222222222",
		"Opponent": "11111111-1111-bbbb-3333-333333333333",
		"ProtagonistMove": "C",
		"OpponentMove": "B"
	}
]`

func TestFileSystemHistoryStore(t *testing.T) {
	newTempStore, closer := createTempFile(t, goldenHistoryData)
	defer closer()
	t.Run("single interaction", func(t *testing.T) {
		store, err := NewFileSystemHistoryStore(newTempStore)
		testhelpers.AssertNoError(t, err)

		got, err := store.GetInteraction(interactionIDOne)
		// testhelpers.AssertNoError(t, err)

		want := types.Interaction{
			ID:              interactionIDOne,
			Protagonist:     player1,
			Opponent:        player2,
			ProtagonistMove: prisoner.Cooperate,
			OpponentMove:    prisoner.Betray,
		}

		testhelpers.AssertInteraction(t, got, want)

	})

	// file_system_store_test.go
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		testhelpers.AssertNoError(t, err)
	})

}
