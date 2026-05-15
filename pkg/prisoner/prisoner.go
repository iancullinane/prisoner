package prisoner

type Move rune

const (
	Cooperate Move = 'C'
	Betray    Move = 'B'
)

type Result rune

const (
	Temptation Result = 'T' // p1 betrays, p2 cooperates
	Reward     Result = 'R' // both cooperate
	Punish     Result = 'P' // both betray
	Sucker     Result = 'S' // p1 cooperates, p2 betrays
)

func (r Result) String() string {
	return string(r)
}

type Payoff[T any] struct {
	High   T // Temptation payoff
	Low    T // Mutual cooperation
	None   T // Mutual betrayal
	Punish T // Sucker payoff
}

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
		return Punish, Temptation
	case m1 == Betray && m2 == Cooperate:
		return Temptation, Punish
	case m1 == Betray && m2 == Betray:
		return Punish, Punish
	default:
		// Should never be reached
		panic("prisoner.Compute: invalid Move")
	}
}
