package main

import (
	"flag"
	"log"

	"github.com/iancullinane/prisoner/dilemma"
	"github.com/iancullinane/prisoner/entity"
)

var value string

func init() {

	flag.StringVar(&value, "value", "", "")
}

func main() {

	log.Println("Prisoners dilemma")
	entA := entity.New("steve", entity.Random)
	entB := entity.New("steve2", entity.Revenge)

	dilemma.PlayRepeated(entA, entB, 10)

	log.Printf("Player 1: %d\tPlayer 2: %d", entA.Score, entB.Score)

}
