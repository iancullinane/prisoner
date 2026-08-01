package store

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

var ErrNotImplemented = errors.New("pretty history not implemented")

func GetPrettyHistoryFromStores(playerID *uuid.UUID, playerStore types.PlayerStore, historyStore types.HistoryStore) (types.PrettyHistory, error) {
	var history types.History
	var err error
	if playerID != nil {
		history, err = historyStore.GetHistoryByPlayerID(*playerID)
	} else {
		history, err = historyStore.GetHistory()
	}
	if err != nil {
		return nil, fmt.Errorf("getting history: %w", err)
	}

	pretty := make(types.PrettyHistory, 0, len(history))
	for _, interaction := range history {
		playerA, err := playerStore.GetPlayerByID(interaction.PlayerA)
		if err != nil {
			return nil, fmt.Errorf("resolving player %s: %w", interaction.PlayerA, err)
		}
		playerB, err := playerStore.GetPlayerByID(interaction.PlayerB)
		if err != nil {
			return nil, fmt.Errorf("resolving player %s: %w", interaction.PlayerB, err)
		}

		pretty = append(pretty, types.PrettyInteraction{
			PlayerAName: playerA.Name,
			PlayerBName: playerB.Name,
			PlayerAMove: interaction.PlayerAMove,
			PlayerBMove: interaction.PlayerBMove,
			PlayedAt:    interaction.PlayedAt,
		})
	}
	return pretty, nil
}
