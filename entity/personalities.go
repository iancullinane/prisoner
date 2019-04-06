package entity

import (
	"log"
	"math/rand"
	"time"
)

// AlwaysCooperate will always act altruistically, giving 3
func AlwaysCooperate(lastMove, oppLastMove string) string {
	return "COOPERATE"
}

// AlwaysCheat will never add +1
func AlwaysCheat(lastMove, oppLastMove string) string {
	return "CHEAT"
}

// CopyCat will blindly do whatever their opponent did last time, excpet the
// first round where they will COOOPERATE
func CopyCat(lastMove, oppLastMove string) string {
	log.Printf("Copy cat will copy: %s", oppLastMove)
	return oppLastMove
}

// Random is random
func Random(lastMove, oppLastMove string) string {

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
