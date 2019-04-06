package entity

import (
	"log"
	"math/rand"
	"time"
)

type Memory struct {
	lastMove    string
	oppLastMove string
	betrayed    bool
}

// AlwaysCooperate will always act altruistically, giving 3
func AlwaysCooperate(mem Memory) string {
	return "COOPERATE"
}

// AlwaysCheat will never add +1
func AlwaysCheat(mem Memory) string {
	return "CHEAT"
}

// CopyCat will blindly do whatever their opponent did last time, excpet the
// first round where they will COOOPERATE
func CopyCat(mem Memory) string {
	log.Printf("Copy cat will copy: %s", mem.oppLastMove)
	return mem.oppLastMove
}

// func Revenge(lastMove, oppLastMove string) string {

// }

// Random is random
func Random(mem Memory) string {

	rand.Seed(time.Now().UnixNano())
	if rand1() == true {
		return "COOPERATE"
	}

	return "CHEAT"
}

func rand1() bool {
	val := rand.Float32()
	return val < 0.4
}
