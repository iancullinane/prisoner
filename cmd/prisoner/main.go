package main

import (
	"flag"

	pkg "github.com/iancullinane/prisoner/pkg/prisoner"
)

var numOfEnts, numOfRounds int

func init() {

	flag.IntVar(&numOfEnts, "entities", 150, "")
	flag.IntVar(&numOfRounds, "rounds", 5, "")
}

func main() {
	flag.Parse()

	prisoner := pkg.New()
	players := prisoner.GetPersonalities(numOfEnts)
	prisoner.PlayTournament(players, numOfRounds)
	prisoner.PrintResults(players)
}
