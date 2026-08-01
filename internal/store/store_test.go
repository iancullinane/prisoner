package store_test

import (
	"reflect"
	"testing"

	"github.com/iancullinane/prisoner/internal/store"
	"github.com/iancullinane/prisoner/internal/store/memory"
	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

func TestGetPrettyHistoryFromStores(t *testing.T) {
	t.Run("resolves player names onto interactions", func(t *testing.T) {
		logger := testhelpers.NoopLogger()
		players := memory.NewInMemoryPlayerStore(logger)
		history := memory.NewInMemoryHistoryStore(logger)

		alice, err := players.GetOrCreatePlayer(testhelpers.TestPlayerOneName)
		testhelpers.AssertNoError(t, err)
		bob, err := players.GetOrCreatePlayer(testhelpers.TestPlayerTwoName)
		testhelpers.AssertNoError(t, err)

		interaction := types.NewInteraction(alice.ID, bob.ID, prisoner.Cooperate, prisoner.Betray)
		testhelpers.AssertNoError(t, history.RecordInteraction(interaction))

		got, err := store.GetPrettyHistoryFromStores(nil, players, history)
		testhelpers.AssertNoError(t, err)

		want := types.PrettyHistory{{
			PlayerAName: testhelpers.TestPlayerOneName,
			PlayerBName: testhelpers.TestPlayerTwoName,
			PlayerAMove: prisoner.Cooperate,
			PlayerBMove: prisoner.Betray,
			PlayedAt:    interaction.PlayedAt,
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot  %+v\nwant %+v", got, want)
		}
	})

	t.Run("errors when an interaction references an unknown player", func(t *testing.T) {
		logger := testhelpers.NoopLogger()
		players := memory.NewInMemoryPlayerStore(logger)
		history := memory.NewInMemoryHistoryStore(logger)

		interaction := types.NewInteraction(testhelpers.TestPlayerOneID, testhelpers.TestPlayerTwoID, prisoner.Betray, prisoner.Betray)
		testhelpers.AssertNoError(t, history.RecordInteraction(interaction))

		_, err := store.GetPrettyHistoryFromStores(nil, players, history)
		if err == nil {
			t.Fatal("got no error, want one for unknown player")
		}
	})

	t.Run("filters to interactions involving the given player", func(t *testing.T) {
		logger := testhelpers.NoopLogger()
		players := memory.NewInMemoryPlayerStore(logger)
		history := memory.NewInMemoryHistoryStore(logger)

		alice, err := players.GetOrCreatePlayer(testhelpers.TestPlayerOneName)
		testhelpers.AssertNoError(t, err)
		bob, err := players.GetOrCreatePlayer(testhelpers.TestPlayerTwoName)
		testhelpers.AssertNoError(t, err)
		carol, err := players.GetOrCreatePlayer("Carol")
		testhelpers.AssertNoError(t, err)

		aliceVsBob := types.NewInteraction(alice.ID, bob.ID, prisoner.Cooperate, prisoner.Betray)
		testhelpers.AssertNoError(t, history.RecordInteraction(aliceVsBob))
		bobVsCarol := types.NewInteraction(bob.ID, carol.ID, prisoner.Betray, prisoner.Betray)
		testhelpers.AssertNoError(t, history.RecordInteraction(bobVsCarol))

		got, err := store.GetPrettyHistoryFromStores(&alice.ID, players, history)
		testhelpers.AssertNoError(t, err)

		want := types.PrettyHistory{{
			PlayerAName: testhelpers.TestPlayerOneName,
			PlayerBName: testhelpers.TestPlayerTwoName,
			PlayerAMove: prisoner.Cooperate,
			PlayerBMove: prisoner.Betray,
			PlayedAt:    aliceVsBob.PlayedAt,
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot  %+v\nwant %+v", got, want)
		}
	})
}
