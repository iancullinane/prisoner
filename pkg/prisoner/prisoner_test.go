package prisoner

import (
	"fmt"
	"testing"
	"time"
)

var traditionalIntRewardSpec = []struct {
	result Result
	want   int
}{
	{Temptation, 2},
	{Reward, 1},
	{Punish, 0},
	{Sucker, -1},
}

func TestRefactor_Reward_resultMapsToInt(t *testing.T) {
	traditional := Payoff[int]{
		High:   2,
		Low:    1,
		None:   0,
		Punish: -1,
	}
	for _, tc := range traditionalIntRewardSpec {
		t.Run(string(tc.result), func(t *testing.T) {
			got := traditional.Compute(tc.result)
			if got != tc.want {
				t.Errorf("For(%q): got %d, want %d", tc.result, got, tc.want)
			}
		})
	}
}

var timePayoffSpec = []struct {
	result Result
	want   time.Duration
}{
	{Temptation, 0},
	{Reward, 365 * 24 * time.Hour},
	{Punish, 2 * 365 * 24 * time.Hour},
	{Sucker, 3 * 365 * 24 * time.Hour},
}

func TestRefactor_Reward_resultMapsToTime(t *testing.T) {
	traditional := Payoff[time.Duration]{
		High:   0,
		Low:    365 * 24 * time.Hour,
		None:   2 * 365 * 24 * time.Hour,
		Punish: 3 * 365 * 24 * time.Hour,
	}
	for _, tc := range timePayoffSpec {
		t.Run(tc.result.String(), func(t *testing.T) {
			got := traditional.Compute(tc.result)
			if got != tc.want {
				t.Errorf("Compute(%s): got %v, want %v", tc.result.String(), got, tc.want)
			}
		})
	}
}

func TestRefactor_Result(t *testing.T) {
	var tests = []struct {
		name string
		m1   Move
		m2   Move
		r1   Result
		r2   Result
	}{
		{"both_cooperate", Cooperate, Cooperate, Reward, Reward},
		{"both_betray", Betray, Betray, Punish, Punish},
		{"p1_cooperate_p2_betray", Cooperate, Betray, Punish, Temptation},
		{"p1_betray_p2_cooperate", Betray, Cooperate, Temptation, Punish},
	}
	for _, tc := range tests {
		got1, got2 := Play(tc.m1, tc.m2)

		assertResult(t, got1, tc.r1, got2, tc.r2)
	}
}

func assertResult(t *testing.T, got1, want1, got2, want2 Result) {
	t.Helper()
	p1 := "player1"
	p2 := "player2"

	if got1 != want1 {
		p1 += fmt.Sprintf(": got %s, want %s", got1, want1)
	}
	if got2 != want2 {
		p2 += fmt.Sprintf(": got %s, want %s", got2, want2)
	}

	if p1 != "player1" || p2 != "player2" {
		t.Errorf("%s %s", p1, p2)
	}

}
