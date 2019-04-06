package entity

type Entity struct {
	Name        string
	Score       int32
	behavior    func(mem Memory) string
	lastMove    string
	oppLastMove string
	betrayed    bool
}

func New(name string, behavior func(mem Memory) string) *Entity {
	return &Entity{
		Name:        name,
		Score:       0,
		behavior:    behavior,
		lastMove:    "COOPERATE", // Assume cooperation
		oppLastMove: "COOPERATE", // ...from everyone
		betrayed:    false,
	}
}

// Play entity will return either "COOPERATE" or "CHEAT"
func (e *Entity) Play() string {

	memory := Memory{
		lastMove:    e.lastMove,
		oppLastMove: e.oppLastMove,
		betrayed:    e.betrayed,
	}

	move := e.behavior(memory)
	return move
}

// RecordMoves updates the object with its last move and the opponents move
func (e *Entity) RecordMoves(move, oppMove string) {
	e.lastMove = move
	e.oppLastMove = oppMove

	if oppMove == "CHEAT" {
		e.betrayed = true
	}
}

// Reset the score for a new round
func (e *Entity) Reset() {
	e.Score = 0
}
