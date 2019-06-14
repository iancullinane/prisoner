package main

import (
	"flag"
	"log"

	"github.com/iancullinane/prisoner/src/utils"

	"github.com/iancullinane/prisoner/src/dilemma"
	"github.com/iancullinane/prisoner/src/entity"
)

var numOfEnts, numOfRounds int

func init() {

	flag.IntVar(&numOfEnts, "entities", 150, "")
	flag.IntVar(&numOfRounds, "rounds", 5, "")
}

func main() {

	flag.Parse()

	entities := make([]*entity.Entity, 0)
	behaviors := entity.NewBehaviorFactory()

	for i := 1; i <= numOfEnts; i++ {

		tmpEntity := entity.New(
			utils.GetRandomName(),
			behaviors.GetRandomBehavior(),
		)
		entities = append(entities, tmpEntity)
	}

	dilemma.PlayOneTournament(entities, numOfRounds)

	utils.SortByScore(entities)

	log.Printf("%-10s\t%s\t%s", "Name", "Score", "Personality")
	log.Printf("------------------------------------")
	for _, e := range entities {
		log.Printf("%-10s\t%d\t%s", e.Name, e.Score, e.GetBehaviorName())
	}

}
