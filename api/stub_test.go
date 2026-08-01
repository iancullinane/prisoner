package api

import (
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

type StubPlayerHistoryStore struct {
	players []types.Player
	history types.History
}

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

func (h *StubHistoryStore) GetHistoryByPlayerID(playerID uuid.UUID) (types.History, error) {
	historyByPlayer := make(types.History, 0)
	for _, interaction := range h.history {
		if interaction.PlayerA == playerID || interaction.PlayerB == playerID {
			historyByPlayer = append(historyByPlayer, interaction)
		}
	}
	return historyByPlayer, nil
}

// MARK: PlayerStore
// =========================================

type StubPlayerStore struct {
	scores                 map[string]int
	players                []types.Player
	getOrCreatePlayerCalls []string
	getOrCreatePlayerError error
}

func (s *StubPlayerStore) GetOrCreatePlayer(name string) (types.Player, error) {
	s.getOrCreatePlayerCalls = append(s.getOrCreatePlayerCalls, name)

	if s.getOrCreatePlayerError != nil {
		return types.Player{}, s.getOrCreatePlayerError
	}

	player := types.Players(s.players).FindByName(name)
	if player != nil {
		return *player, nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return types.Player{}, err
	}

	player = &types.Player{
		ID:   id,
		Name: name,
	}

	if s.players == nil {
		s.players = append(s.players, *player)
	}

	return *player, nil
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

func (s *StubPlayerStore) GetAllPlayers() (types.Players, error) {
	return types.Players(s.players), nil
}

func (s *StubPlayerStore) GetRandomPlayer() (types.Player, error) {
	return types.Players(s.players).GetRandomPlayer()
}

func (s *StubPlayerStore) GetRandomPlayerExcept(exceptID uuid.UUID) (types.Player, error) {
	return types.Players(s.players).GetRandomPlayerExcept(exceptID)
}
