package entity

import (
	"math/rand"
	"time"
)

type BehaviorFactory struct {
	behaviors []func(Memory) string
}

type Memory struct {
	lastMove    string
	oppLastMove string
	betrayed    int
}

func NewBehaviorFactory() *BehaviorFactory {
	return &BehaviorFactory{
		behaviors: []func(Memory) string{
			AlwaysCheat,
			AlwaysCooperate,
			CopyCat,
			Revenge,
			Random,
		},
	}
}

func (f *BehaviorFactory) GetRandomBehavior() func(Memory) string {

	return f.behaviors[rand.Intn(len(f.behaviors))]
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
	return mem.oppLastMove
}

// Revenge will always cheat if it has been cheated against once
func Revenge(mem Memory) string {

	if mem.betrayed > 0 {
		return "CHEAT"
	}

	return "COOPERATE"
}

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
	return val < 0.5
}
