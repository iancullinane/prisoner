package game

import (
	"github.com/iancullinane/prisoner/internal/dilemma"
	"github.com/iancullinane/prisoner/internal/entity"
	"github.com/iancullinane/prisoner/internal/utils"
)

type Game struct {
	entities []*entity.Entity
	behavior *entity.BehaviorFactory
}

func New(numEntities int) *Game {

	entityList := make([]*entity.Entity, numEntities)

	return &Game{
		entities: entityList,
		behavior: entity.NewBehaviorFactory(),
	}
}

func (g *Game) GetEntities() []*entity.Entity {
	return utils.SortByScore(g.entities)
}

func (g *Game) FillEntities() {

	for i := 0; i < len(g.entities); i++ {

		tmpEntity := entity.New(
			utils.GetRandomName(),
			g.behavior.GetRandomBehavior(),
		)
		g.entities[i] = tmpEntity

	}
}

func (g *Game) PlayOneTournament(numOfRounds int) {

	dilemma.PlayOneTournament(g.entities, numOfRounds)
	g.entities = utils.SortByScore(g.entities)
}

func (g *Game) GetGameID() string {
	return "GameID"
}
