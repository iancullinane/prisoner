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
		store := NewFileSystemPlayerStore(newTempStore)

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, got, want)

	})

	t.Run("league from a reader", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store := NewFileSystemPlayerStore(newTempStore)

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

		store := NewFileSystemPlayerStore(newTempStore)

		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("record new player win", func(t *testing.T) {
		newTempStore, closer := createTempFile(t, goldenTest)
		defer closer()

		store := NewFileSystemPlayerStore(newTempStore)

		store.RecordWin("Ian")

		got := store.GetPlayerScore("Ian")
		want := 1
		assertScoreEquals(t, got, want)
	})

}
