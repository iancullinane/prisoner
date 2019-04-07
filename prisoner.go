package main

import (
	"flag"
	"log"

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

	for _, e := range entities {
		log.Printf("%-10s\t%d\t%s", e.Name, e.Score, e.GetBehaviorName())
	}

}
