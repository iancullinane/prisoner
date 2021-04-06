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
	players := prisoner.GetPersonalities(10)
	prisoner.PlayTournament(players, 5)
	prisoner.PrintResults(players)
}
