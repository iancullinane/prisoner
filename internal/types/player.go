package types

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

const devPlayerName = "Steve"

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

func NewPlayerFromID(id string) (*Player, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &Player{
		ID:   uuid,
		Name: devPlayerName,
	}, nil
}

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	RecordWin(name string) error
}
