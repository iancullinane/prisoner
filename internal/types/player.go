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
	ProtagonistMove   prisoner.Move
	OpponentMove      prisoner.Move
	ProtagonistResult prisoner.Result
	OpponentResult    prisoner.Result
}

func (r Round) String() string {
	return fmt.Sprintf("%v vs %v: %v, %v", r.ProtagonistMove, r.OpponentMove, r.ProtagonistResult, r.OpponentResult)
}
