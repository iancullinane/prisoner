package types

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

type History []Interaction

type HistoryStore interface {
	GetHistory() (History, error)
	RecordInteraction(interaction Interaction) error
	// GetInteraction(id uuid.UUID) *Interaction
	GetHistoryByPlayerID(playerID uuid.UUID) (History, error)
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
	ID          uuid.UUID     `json:"id"`
	PlayerA     uuid.UUID     `json:"playerA"`
	PlayerB     uuid.UUID     `json:"playerB"`
	PlayerAMove prisoner.Move `json:"playerAMove"`
	PlayerBMove prisoner.Move `json:"playerBMove"`
	PlayedAt    time.Time     `json:"playedAt"`
}

func NewInteraction(playerA, playerB uuid.UUID, playerAMove, playerBMove prisoner.Move) Interaction {
	return Interaction{
		ID:          uuid.New(),
		PlayerA:     playerA,
		PlayerB:     playerB,
		PlayerAMove: playerAMove,
		PlayerBMove: playerBMove,
		PlayedAt:    time.Now(),
	}
}

func NewInteractionFromJSON(r io.Reader) (Interaction, error) {
	var i Interaction
	err := json.NewDecoder(r).Decode(&i)
	if err != nil {
		return Interaction{}, err
	}

	return i, err
}

func (h History) Find(name uuid.UUID) *Interaction {
	for i, p := range h {
		if p.ID == name {
			return &h[i]
		}
	}
	return nil
}

func (r Interaction) PrintScore(scoring prisoner.Payoff[int]) string {
	result1, result2 := prisoner.Play(r.PlayerAMove, r.PlayerBMove)
	return fmt.Sprintf("%v : %v", scoring.Compute(result1), scoring.Compute(result2))
}

func (r Interaction) String() string {
	return fmt.Sprintf("%s vs %s: %s, %s", r.PlayerA, r.PlayerB, r.PlayerAMove, r.PlayerBMove)
}

// PrintInteraction renders the interaction as a human-readable line, or as
// indented JSON (including the computed score) when asJSON is true.
func (r Interaction) PrintInteraction(scoring prisoner.Payoff[int], asJSON bool) string {
	if asJSON {
		out := struct {
			Interaction
			Score string `json:"score"`
			// PlayerAScore int `json:"playerAScore"`
			// PlayerBScore int `json:"playerBScore"`
		}{Interaction: r, Score: r.PrintScore(scoring)}

		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return fmt.Sprintf("error marshaling interaction: %v", err)
		}
		return string(b)
	}
	return fmt.Sprintf("%s: %s", r, r.PrintScore(scoring))
}
