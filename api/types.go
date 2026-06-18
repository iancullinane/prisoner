package api

import (
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

type PlayRequest struct {
	PlayerA     uuid.UUID     `json:"player_a"`
	PlayerB     uuid.UUID     `json:"player_b"`
	PlayerAMove prisoner.Move `json:"player_a_move"`
	PlayerBMove prisoner.Move `json:"player_b_move"`
}

type PlayResponse struct {
	ID           uuid.UUID       `json:"ID"`
	PlayerAScore prisoner.Result `json:"player_a_score"`
	PlayerBScore prisoner.Result `json:"player_b_score"`
}
