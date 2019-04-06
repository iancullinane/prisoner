package entity

type Entity struct {
	Name        string
	Score       int32
	behavior    func(lastMove, oppLastMove string) string
	lastMove    string
	oppLastMove string
}

func New(name string, behavior func(lastMove, oppLastMove string) string) *Entity {
	return &Entity{
		Name:        name,
		Score:       0,
		behavior:    behavior,
		lastMove:    "COOPERATE", // Assume cooperation
		oppLastMove: "COOPERATE", // ...from everyone
	}
}

// Play entity will return either "COOPERATE" or "CHEAT"
func (e *Entity) Play() string {
	move := e.behavior(e.lastMove, e.oppLastMove)
	return move
}

// RecordMoves updates the object with its last move and the opponents move
func (e *Entity) RecordMoves(move, oppMove string) {
	e.lastMove = move
	e.oppLastMove = oppMove
}

// Reset the score for a new round
func (e *Entity) Reset() {
	e.Score = 0
}
