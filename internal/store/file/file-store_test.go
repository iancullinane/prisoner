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

		got, err := store.GetPlayerScore("Chris")
		assertNoError(t, err)
		want := 33

		assertScoreEquals(t, got, want)

	})

	t.Run("record player win", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store, err := NewFileSystemPlayerStore(newTempStore)
		assertNoError(t, err)

		err = store.RecordWin("Chris")
		assertNoError(t, err)

		got, err := store.GetPlayerScore("Chris")
		assertNoError(t, err)
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("record new player win", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store, err := NewFileSystemPlayerStore(newTempStore)
		assertNoError(t, err)

		err = store.RecordWin("Ian")
		assertNoError(t, err)

		got, err := store.GetPlayerScore("Ian")
		assertNoError(t, err)
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

}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an error but didn't want one %v", err)
	}
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
