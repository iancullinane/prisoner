package game

import (
	"github.com/iancullinane/prisoner/src/dilemma"
	"github.com/iancullinane/prisoner/src/entity"
	"github.com/iancullinane/prisoner/src/utils"
)

type Game struct {
	entities []*entity.Entity
	behavior *entity.BehaviorFactory
}

// behaviors := entity.NewBehaviorFactory()
func New(numEntities int) *Game {

	entityList := make([]*entity.Entity, numEntities)

	return &Game{
		entities: entityList,
		behavior: entity.NewBehaviorFactory(),
	}
}

func (g *Game) GetEntities() []*entity.Entity {
	return g.entities
}

func (g *Game) AddEntity(entIn *entity.Entity) {
	tmpEntity := entity.New(
		utils.GetRandomName(),
		g.behavior.GetRandomBehavior(),
	)
	g.entities = append(g.entities, tmpEntity)
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
	utils.SortByScore(g.entities)
}

func (g *Game) GetGameID() string {
	return "GameID"
}
