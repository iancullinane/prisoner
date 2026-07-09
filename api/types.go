package api

import (
	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

type PlayRequest struct {
	PlayerA     uuid.UUID     `json:"playerA"`
	PlayerB     uuid.UUID     `json:"playerB"`
	PlayerAMove prisoner.Move `json:"playerAMove"`
	PlayerBMove prisoner.Move `json:"playerBMove"`
}

type PlayResponse struct {
	ID           uuid.UUID       `json:"id"`
	PlayerAScore prisoner.Result `json:"playerAScore"`
	PlayerBScore prisoner.Result `json:"playerBScore"`
}
