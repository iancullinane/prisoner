package main

import (
	"flag"
	"log"

	"github.com/iancullinane/prisoner/src/game"
	"github.com/iancullinane/prisoner/src/utils"
)

var numOfEnts, numOfRounds int

func init() {

	flag.IntVar(&numOfEnts, "entities", 150, "")
	flag.IntVar(&numOfRounds, "rounds", 5, "")
}

func main() {
	flag.Parse()

	game := game.New(10)
	game.FillEntities()

	game.PlayOneTournament(10)

	// dilemma.PlayOneTournament(entities, numOfRounds)

	utils.SortByScore(game.GetEntities())

	log.Printf("%-10s\t%s\t%s", "Name", "Score", "Personality")
	log.Printf("------------------------------------")
	for _, e := range game.GetEntities() {
		log.Printf("%-10s\t%d\t%s", e.Name, e.Score, e.GetBehaviorName())
	}

}
