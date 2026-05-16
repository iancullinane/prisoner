package testhelpers

import (
	"reflect"
	"testing"

	"github.com/iancullinane/prisoner/internal/types"
)

func AssertLeague(t testing.TB, got, want []types.Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
