package file

import (
	"testing"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

const goldenTest = `[
		{"Name": "Chris", "Wins": 33},
		{"Name": "Cleo", "Wins": 10}
	]`

func TestFileSystemStore(t *testing.T) {
	newTempStore, closer := createTempFile(t, goldenTest)
	defer closer()
	t.Run("get score for a placer", func(t *testing.T) {
		store, err := NewFileSystemPlayerStore(newTempStore)
		assertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, got, want)

	})

	t.Run("league from a reader", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store, err := NewFileSystemPlayerStore(newTempStore)
		assertNoError(t, err)

		got := store.GetLeague()

		want := []types.Player{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}

		testhelpers.AssertLeague(t, got, want)
		// should get the same result every time
		got = store.GetLeague()
		testhelpers.AssertLeague(t, got, want)

	})

	t.Run("record player win", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store, err := NewFileSystemPlayerStore(newTempStore)
		assertNoError(t, err)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("record new player win", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store, err := NewFileSystemPlayerStore(newTempStore)
		assertNoError(t, err)

		store.RecordWin("Ian")

		got := store.GetPlayerScore("Ian")
		want := 1
		assertScoreEquals(t, got, want)
	})

	// file_system_store_test.go
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})

	// file_system_store_test.go
	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := types.League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}

		testhelpers.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		testhelpers.AssertLeague(t, got, want)
	})

}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't want one %v", err)
	}
}
