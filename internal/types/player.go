package types

import (
	"encoding/json"
	"fmt"
	"io"

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

type Players []Player

func NewPlayers(rdr io.Reader) (Players, error) {
	var players Players
	err := json.NewDecoder(rdr).Decode(&players)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return players, err
}

func (l Players) Find(id uuid.UUID) *Player {
	for i, p := range l {
		if p.ID == id {
			return &l[i]
		}
	}
	return nil
}
