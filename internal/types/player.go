package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
)

const devPlayerName = "Steve"

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type Player struct {
	ID   uuid.UUID
	Name string
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
	// CreatePlayer(name string) (Player, error)
	GetOrCreatePlayer(name string) (Player, error)
	GetPlayerByID(id uuid.UUID) (Player, error)
	GetPlayerByName(name string) (Player, error)
	GetAllPlayers() (Players, error)
	GetRandomPlayer() (Player, error)
	GetRandomPlayerExcept(exceptID uuid.UUID) (Player, error)
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

func (l Players) FindByID(id uuid.UUID) *Player {
	for i, p := range l {
		if p.ID == id {
			return &l[i]
		}
	}
	return nil
}

func (l Players) FindByName(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

func (l Players) GetAllPlayers() Players {
	return l
}

func (l Players) String() string {
	var buf strings.Builder
	for i, p := range l {
		if i > 0 {
			buf.WriteString(",\n")
		}
		fmt.Fprintf(&buf, "%s\t\t(%s)", p.Name, p.ID.String())
	}
	return buf.String()
}

func (l Players) GetRandomPlayer() (Player, error) {
	if len(l) == 0 {
		return Player{}, ErrPlayerNotFound
	}
	return l[rand.Intn(len(l))], nil
}

func (l Players) GetRandomPlayerExcept(exceptID uuid.UUID) (Player, error) {
	candidates := make(Players, 0, len(l))
	for _, p := range l {
		if p.ID != exceptID {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) == 0 {
		return Player{}, ErrPlayerNotFound
	}
	return candidates[rand.Intn(len(candidates))], nil
}
