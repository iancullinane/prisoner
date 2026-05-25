package types

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

type History []Interaction

type HistoryStore interface {
	GetHistory() (History, error)
	RecordInteraction(interaction Interaction) error
}

func NewHistory(rdr io.Reader) (History, error) {
	var history History
	err := json.NewDecoder(rdr).Decode(&history)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return history, err
}

type Interaction struct {
	ID              uuid.UUID
	Protagonist     uuid.UUID
	Opponent        uuid.UUID
	ProtagonistMove prisoner.Move
	OpponentMove    prisoner.Move
	// CreatedAt       time.Time
}

func NewInteraction(protagonist, opponent uuid.UUID, protagonistMove, opponentMove prisoner.Move) Interaction {
	return Interaction{
		ID:              uuid.New(),
		Protagonist:     protagonist,
		Opponent:        opponent,
		ProtagonistMove: protagonistMove,
		OpponentMove:    opponentMove,
		// CreatedAt:       time.Now(),
	}
}

func (r Interaction) PrintScore(scoring prisoner.Payoff[int]) string {
	result1, result2 := prisoner.Play(r.ProtagonistMove, r.OpponentMove)
	return fmt.Sprintf("%v : %v", scoring.Compute(result1), scoring.Compute(result2))
}

func (r Interaction) String() string {
	return fmt.Sprintf("%s vs %s: %s, %s", r.Protagonist, r.Opponent, r.ProtagonistMove, r.OpponentMove)
}
