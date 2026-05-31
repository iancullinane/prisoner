package api

import (
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

// MARK: HistoryStore
// =========================================

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
// =========================================

type StubPlayerStore struct {
	scores            map[string]int
	players           []types.Player
	createPlayerCalls []string
}

func (s *StubPlayerStore) CreatePlayer(name string) (types.Player, error) {
	s.createPlayerCalls = append(s.createPlayerCalls, name)

	id, err := uuid.NewRandom()
	if err != nil {
		return types.Player{}, err
	}

	player := types.Player{
		ID:   id,
		Name: name,
	}

	if s.players == nil {
		s.players = append(s.players, player)
	}

	return player, nil
}

func (s *StubPlayerStore) GetPlayerByID(id uuid.UUID) (types.Player, error) {
	player := types.Players(s.players).FindByID(id)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}
	return *player, nil
}

func (s *StubPlayerStore) GetPlayerByName(name string) (types.Player, error) {
	player := types.Players(s.players).FindByName(name)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}
	return *player, nil
}
