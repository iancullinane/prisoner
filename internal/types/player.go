package types

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
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

//
