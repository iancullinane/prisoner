package main

import (
	"flag"
	"log"
	"sort"

	"github.com/iancullinane/prisoner/utils"

	"github.com/iancullinane/prisoner/dilemma"
	"github.com/iancullinane/prisoner/entity"
)

var value string

func init() {

	flag.StringVar(&value, "value", "", "")
}

func main() {

	entities := make([]*entity.Entity, 0)
	behaviors := entity.NewBehaviorFactory()

	for i := 1; i <= 10; i++ {

		tmpEntity := entity.New(
			utils.GetRandomName(),
			behaviors.GetRandomBehavior(),
		)
		entities = append(entities, tmpEntity)
	}

	dilemma.PlayOneTournament(entities, 10)

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Score < entities[j].Score
	})

	log.Printf("%-10s\t%s\t%s", "Name", "Score", "Personality")
	log.Printf("------------------------------------")
	for _, e := range entities {
		log.Printf("%-10s\t%d\t%s", e.Name, e.Score, e.GetBehaviorName())
	}

}
