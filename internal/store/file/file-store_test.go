package file

import (
	"strings"
	"testing"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

func TestFileSystemStore(t *testing.T) {
	database := strings.NewReader(`[
		{"Name": "Chris", "Wins": 33},
		{"Name": "Cleo", "Wins": 10}
	]`)
	t.Run("get score for a placer", func(t *testing.T) {
		store := NewFilesSystemPlayerStore(database)

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, got, want)

	})

	t.Run("league from a reader", func(t *testing.T) {
		store := NewFilesSystemPlayerStore(database)

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

}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("score not equal got %v want %v", got, want)
	}
}
