package prisoner

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

type rune = int32

type runeEnum interface {
	~rune // ← set is {int32, Move, Result, and any `type X rune` within this package}
	// allow both Move and Result to satisfy this interface
	Valid() bool
}

// MARK: Move Types
// ==================

type Move rune

const (
	Cooperate Move = 'C'
	Betray    Move = 'B'
)

var (
	moves = []Move{Cooperate, Betray}
)

func RandomMove() Move {
	oneOrZero := rand.Intn(2)
	return moves[oneOrZero]
}

func (m Move) Valid() bool {
	return m == Cooperate || m == Betray
}

func (m Move) String() string {
	return string(m)
}

func (m Move) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(m))
}

func (m *Move) UnmarshalJSON(data []byte) error {
	return unmarshalRuneEnum(data, m)
}

// MARK: Result Types
// ==================
type Result rune

const (
	Temptation Result = 'T' // p1 betrays, p2 cooperates
	Reward     Result = 'R' // both cooperate
	Punish     Result = 'P' // both betray
	Sucker     Result = 'S' // p1 cooperates, p2 betrays
)

func (r Result) Valid() bool {
	return r == Temptation || r == Reward || r == Punish || r == Sucker
}

func (r Result) String() string {
	return string(r)
}

func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

func (r *Result) UnmarshalJSON(data []byte) error {
	return unmarshalRuneEnum(data, r)
}

// Payoff is a generic type for implementing scores
// each score type could be different, including time
// itself, but mostly just integers...
type Payoff[T any] struct {
	High   T // Temptation payoff
	Low    T // Mutual cooperation
	None   T // Mutual betrayal
	Punish T // Sucker payoff
}

// Classic: T=0, R=-1, P=0, S=-3
var Classic = Payoff[int]{
	High:   0,
	Low:    -1,
	None:   0,
	Punish: -3,
}

// Positive: T=+3, R=+2, P=0, S=-1
var Positive = Payoff[int]{
	High:   3,
	Low:    2,
	None:   0,
	Punish: -1,
}

var Algebraic = Payoff[string]{
	High:   "T",
	Low:    "R",
	None:   "P",
	Punish: "S",
}

// MARK: Compute Functions

func (payoff Payoff[T]) Compute(result Result) T {
	switch result {
	case Temptation:
		return payoff.High
	case Reward:
		return payoff.Low
	case Punish:
		return payoff.None
	case Sucker:
		return payoff.Punish
	default:
		panic("prisoner.Payoff.Compute: invalid reward")
	}
}

func Play(m1, m2 Move) (r1, r2 Result) {
	switch {
	case m1 == Cooperate && m2 == Cooperate:
		return Reward, Reward
	case m1 == Cooperate && m2 == Betray:
		return Sucker, Temptation
	case m1 == Betray && m2 == Cooperate:
		return Temptation, Sucker
	case m1 == Betray && m2 == Betray:
		return Punish, Punish
	default:
		// Should never be reached
		panic("prisoner.Compute: invalid Move")
	}
}

// MARK: Rune Unmarshaller

func unmarshalRuneEnum[T runeEnum](data []byte, dst *T) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) != 1 {
		return fmt.Errorf("invalid %T: %q", *dst, s)
	}
	v := T(s[0])
	if !v.Valid() {
		return fmt.Errorf("invalid %T: %q", *dst, s)
	}
	*dst = v
	return nil
}
