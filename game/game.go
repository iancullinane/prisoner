package game

import (
	"github.com/iancullinane/prisoner/dilemma"
	"github.com/iancullinane/prisoner/entity"
	"github.com/iancullinane/prisoner/utils"
)

type Game struct {
	entities []*entity.Entity
}

func New(numEntities int) *Game {

	entityList := make([]*entity.Entity, numEntities)

	return &Game{
		entities: entityList,
	}
}

func (g *Game) GetEntities() []*entity.Entity {
	return g.entities
}

func (g *Game) AddEntity(entIn *entity.Entity) {
	g.entities = append(g.entities, entIn)
}

func (g *Game) PlayOneTournament(numOfRounds int) {
	dilemma.PlayOneTournament(g.entities, numOfRounds)
	utils.SortByScore(g.entities)
}
