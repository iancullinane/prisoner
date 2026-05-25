package api

import (
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

type StubHistoryStore struct {
	history                types.History
	recordInteractionCalls []types.Interaction
}

func (h *StubHistoryStore) GetHistory() (types.History, error) {
	return h.history, nil
}

func (h *StubHistoryStore) RecordInteraction(interaction types.Interaction) error {
	h.recordInteractionCalls = append(h.recordInteractionCalls, interaction)
	return nil
}

// MARK: PlayerStore

type StubPlayerStore struct {
	scores   map[string]int
	players  map[uuid.UUID]types.Player
	winCalls []string
	league   types.League
}

func (s *StubPlayerStore) GetPlayerScore(name string) (int, error) {
	score := s.scores[name]
	return score, nil
}

func (s *StubPlayerStore) RecordWin(name string) error {
	s.winCalls = append(s.winCalls, name)
	return nil
}

func (s *StubPlayerStore) GetLeague() (types.League, error) {
	return s.league, nil
}
