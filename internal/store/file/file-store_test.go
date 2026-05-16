package file

import (
	"strings"
	"testing"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Chris", "Wins": 10},
			{"Name": "Cleo", "Wins": 10}
		]`)

		store := NewFilesSystemPlayerStore(database)

		got := store.GetLeague()

		want := []types.Player{
			{Name: "Chris", Wins: 10},
			{Name: "Cleo", Wins: 10},
		}

		testhelpers.AssertLeague(t, got, want)
		got = store.GetLeague()
		testhelpers.AssertLeague(t, got, want)

	})

}
