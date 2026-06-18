package types

import (
	"strings"
	"testing"
)

// This currently FAILS — String() prints to stdout and returns "".
func TestPlayersString(t *testing.T) {
	players := Players{*NewPlayer("Alice")}
	if got := players.String(); !strings.Contains(got, "Alice") {
		t.Errorf("String() = %q, want it to contain the player name", got)
	}
}
