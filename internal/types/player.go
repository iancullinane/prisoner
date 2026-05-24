package types

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

type Player struct {
	ID   uuid.UUID
	Name string
	Wins int
}

func NewPlayer(name string) *Player {
	if name == "" {
		name = randomdata.FirstName(randomdata.RandomGender)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("failed to generate uuid for player: %v", err))
	}

	return &Player{
		ID:   id,
		Name: name,
	}
}

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	RecordWin(name string) error
	GetLeague() (League, error)
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
	Protagonist       Player
	Opponent          Player
	ProtagonistMove   prisoner.Move
	OpponentMove      prisoner.Move
	ProtagonistResult prisoner.Result
	OpponentResult    prisoner.Result
}

func NewRound(protagonist, opponent *Player, protagonistMove, opponentMove prisoner.Move) Round {
	return Round{
		Protagonist:     *protagonist,
		Opponent:        *opponent,
		ProtagonistMove: protagonistMove,
		OpponentMove:    opponentMove,
	}
}

func (r Round) PrintScore(scoring prisoner.Payoff[int]) string {
	return fmt.Sprintf("%s vs %s: %s, %s", r.Protagonist.Name, r.Opponent.Name, r.ProtagonistMove, r.OpponentMove)
}

func (r Round) String() string {
	return fmt.Sprintf("%s vs %s: %s, %s", r.Protagonist.Name, r.Opponent.Name, r.ProtagonistMove, r.OpponentMove)
}
