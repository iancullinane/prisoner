package pkg

import (
	"log"

	"github.com/iancullinane/prisoner/internal/dilemma"
	"github.com/iancullinane/prisoner/internal/entity"
	"github.com/iancullinane/prisoner/internal/utils"
)

type client struct {
}

func New() *client {
	return &client{}
}

func (c *client) GetPersonalities(n int) []*entity.Entity {

	f := entity.NewBehaviorFactory()
	entityList := make([]*entity.Entity, n)

	for i := 0; i < len(entityList); i++ {

		tmpEntity := entity.New(
			utils.GetRandomName(),
			f.GetRandomBehavior(),
		)
		entityList[i] = tmpEntity

	}

	return entityList
}

func (c *client) PlayTournament(e []*entity.Entity, n int) {
	dilemma.PlayOneTournament(e, n)
}

func (c *client) PrintResults(players []*entity.Entity) {
	resultsSorted := utils.SortByScore(players)
	log.Printf("%-10s\t%s\t%s", "Name", "Score", "Personality")
	log.Printf("------------------------------------")
	for _, e := range resultsSorted {
		log.Printf("%-10s\t%d\t%s", e.Name, e.Score, e.GetBehaviorName())
	}
}
