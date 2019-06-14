package main

import (
	"flag"
	"log"

	"github.com/iancullinane/prisoner/src/game"
	"github.com/iancullinane/prisoner/src/ui"
	"github.com/iancullinane/prisoner/src/utils"
)

var numOfEnts, numOfRounds int

func init() {

	flag.IntVar(&numOfEnts, "entities", 150, "")
	flag.IntVar(&numOfRounds, "rounds", 5, "")
}

func main() {
	flag.Parse()

	ui := ui.Init()

	game := game.New(ui, 10)
	game.FillEntities()

	game.PlayOneTournament(10)

	utils.SortByScore(game.GetEntities())

	log.Printf("%-10s\t%s\t%s", "Name", "Score", "Personality")
	log.Printf("------------------------------------")
	for _, e := range game.GetEntities() {
		log.Printf("%-10s\t%d\t%s", e.Name, e.Score, e.GetBehaviorName())
	}

}
