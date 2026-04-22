package dilemma

import "testing"

func TestResults(t *testing.T) {
	var tests = []struct {
		name  string
		move1 string
		move2 string
		want1 int32
		want2 int32
	}{
		{"cooperate/cooperate", "COOPERATE", "COOPERATE", 2, 2},
		{"cheat/cooperate", "CHEAT", "COOPERATE", 3, -1},
		{"cooperate/cheat", "COOPERATE", "CHEAT", -1, 3},
		{"cheat/cheat", "CHEAT", "CHEAT", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, r2 := Compute(tt.move1, tt.move2)
			assertCorrectResult(t, r1, r2, tt.want1, tt.want2)
		})
	}
}

func assertCorrectResult(t testing.TB, result1, result2, expected1, expected2 int32) {
	t.Helper()
	if result1 != expected1 {
		t.Errorf("first result wrong %d want %d", result1, expected1)
	}
	if result2 != expected2 {
		t.Errorf("second result wrong %d want %d", result1, expected1)
	}
}
