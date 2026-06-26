package prisoner

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestPayoff_IntCompute(t *testing.T) {
	intSpec := []struct {
		name   string
		payoff Payoff[int]
		result Result
		want   int
	}{
		{"Classic/Temptation", Classic, Temptation, 0},
		{"Classic/Reward", Classic, Reward, -1},
		{"Classic/Punish", Classic, Punish, 0},
		{"Classic/Sucker", Classic, Sucker, -3},
		{"Positive/Temptation", Positive, Temptation, 3},
		{"Positive/Reward", Positive, Reward, 2},
		{"Positive/Punish", Positive, Punish, 0},
		{"Positive/Sucker", Positive, Sucker, -1},
	}
	for _, tc := range intSpec {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.payoff.Compute(tc.result); got != tc.want {
				t.Errorf("%s.Compute(%q): got %d, want %d", tc.name, tc.result, got, tc.want)
			}
		})
	}

	stringSpec := []struct {
		name   string
		payoff Payoff[string]
		result Result
		want   string
	}{
		{"Algebraic/Temptation", Algebraic, Temptation, "T"},
		{"Algebraic/Reward", Algebraic, Reward, "R"},
		{"Algebraic/Punish", Algebraic, Punish, "P"},
		{"Algebraic/Sucker", Algebraic, Sucker, "S"},
	}
	for _, tc := range stringSpec {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.payoff.Compute(tc.result); got != tc.want {
				t.Errorf("%s.Compute(%q): got %q, want %q", tc.name, tc.result, got, tc.want)
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

func Test_Reward_resultMapsToTime(t *testing.T) {
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

func TestPlay(t *testing.T) {
	var tests = []struct {
		name string
		m1   Move
		m2   Move
		r1   Result
		r2   Result
	}{
		{"both_cooperate", Cooperate, Cooperate, Reward, Reward},
		{"both_betray", Betray, Betray, Punish, Punish},
		{"p1_cooperate_p2_betray", Cooperate, Betray, Sucker, Temptation},
		{"p1_betray_p2_cooperate", Betray, Cooperate, Temptation, Sucker},
	}
	for _, tc := range tests {
		got1, got2 := Play(tc.m1, tc.m2)

		assertResult(t, got1, tc.r1, got2, tc.r2)
	}
}

func TestMove_JSON(t *testing.T) {
	cases := []struct {
		move Move
		json string
	}{
		{Cooperate, `"C"`},
		{Betray, `"B"`},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.move)
		if err != nil || string(got) != c.json {
			t.Errorf("marshal %v: got %s err=%v, want %s", c.move, got, err, c.json)
		}
		var back Move
		if err := json.Unmarshal([]byte(c.json), &back); err != nil || back != c.move {
			t.Errorf("unmarshal %s: got %v err=%v, want %v", c.json, back, err, c.move)
		}
	}

	for _, in := range []string{`""`, `"CC"`, `123`, `"X"`} {
		var m Move
		if err := json.Unmarshal([]byte(in), &m); err == nil {
			t.Errorf("unmarshal %s: expected error", in)
		}
	}
}

func TestResult_JSON(t *testing.T) {
	cases := []struct {
		result Result
		json   string
	}{
		{Temptation, `"T"`},
		{Reward, `"R"`},
		{Punish, `"P"`},
		{Sucker, `"S"`},
	}
	for _, c := range cases {
		got, err := json.Marshal(c.result)
		if err != nil || string(got) != c.json {
			t.Errorf("marshal %v: got %s err=%v, want %s", c.result, got, err, c.json)
		}
		var back Result
		if err := json.Unmarshal([]byte(c.json), &back); err != nil || back != c.result {
			t.Errorf("unmarshal %s: got %v err=%v, want %v", c.json, back, err, c.result)
		}
	}

	for _, in := range []string{`""`, `"TR"`, `123`, `"X"`} {
		var r Result
		if err := json.Unmarshal([]byte(in), &r); err == nil {
			t.Errorf("unmarshal %s: expected error", in)
		}
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
