package types

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

type Player struct {
	ID   uuid.UUID
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

// Round represents a single interaction between two players
// resultint in betrayals or cooperations
type Round struct {
	Protagonist       uuid.UUID
	Opponent          uuid.UUID
	ProtagonistMove   prisoner.Move
	OpponentMove      prisoner.Move
	ProtagonistResult prisoner.Result
	OpponentResult    prisoner.Result
}

func NewRound(protagonist, opponent uuid.UUID, protagonistMove, opponentMove prisoner.Move) Round {
	return Round{
		Protagonist:     protagonist,
		Opponent:        opponent,
		ProtagonistMove: protagonistMove,
		OpponentMove:    opponentMove,
	}
}

func (r Round) PrintScore(scoring prisoner.Payoff[int]) string {
	return fmt.Sprintf("%s vs %s: %s, %s", r.Protagonist, r.Opponent, r.ProtagonistMove, r.OpponentMove)
}

func (r Round) String() string {
	return fmt.Sprintf("%s vs %s: %s, %s", r.Protagonist, r.Opponent, r.ProtagonistMove, r.OpponentMove)
}
