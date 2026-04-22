package entity

import "testing"

func TestPersonalityMoves(t *testing.T) {
	tests := []struct {
		name string
		fn   func(Memory) string
		mem  Memory
		want string
	}{
		{"AlwaysCooperate ignores memory", AlwaysCooperate, Memory{lastMove: "CHEAT", oppLastMove: "CHEAT", betrayed: 9}, "COOPERATE"},
		{"AlwaysCheat ignores memory", AlwaysCheat, Memory{lastMove: "COOPERATE", oppLastMove: "COOPERATE", betrayed: 0}, "CHEAT"},
		{"CopyCat mirrors opponent cooperate", CopyCat, Memory{oppLastMove: "COOPERATE"}, "COOPERATE"},
		{"CopyCat mirrors opponent cheat", CopyCat, Memory{oppLastMove: "CHEAT"}, "CHEAT"},
		{"Revenge cooperates when never betrayed", Revenge, Memory{betrayed: 0}, "COOPERATE"},
		{"Revenge cheats after one betrayal", Revenge, Memory{betrayed: 1}, "CHEAT"},
		{"Revenge cheats after many betrayals", Revenge, Memory{betrayed: 100}, "CHEAT"},
		{"Tolerant cooperates at five betrayals", Tolerant, Memory{betrayed: 5}, "COOPERATE"},
		{"Tolerant cheats after more than five", Tolerant, Memory{betrayed: 6}, "CHEAT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.mem)
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestBehaviorFactoryGetBehaviorByName(t *testing.T) {
	f := NewBehaviorFactory()
	tests := []struct {
		key  string
		want string
	}{
		{"cheater", "cheater"},
		{"niceguy", "niceguy"},
		{"copycat", "copycat"},
		{"revenge", "revenge"},
		{"tolerant", "tolerant"},
		{"random", "random"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			b := f.GetBehaviorByName(tt.key)
			if b.GetName() != tt.want {
				t.Fatalf("GetName() got %q want %q", b.GetName(), tt.want)
			}
		})
	}
}

func TestEntityPlay(t *testing.T) {
	cooperate := Behavior{behaviorName: "niceguy", behavior: AlwaysCooperate}
	cheat := Behavior{behaviorName: "cheater", behavior: AlwaysCheat}
	copycat := Behavior{behaviorName: "copycat", behavior: CopyCat}
	revenge := Behavior{behaviorName: "revenge", behavior: Revenge}
	tolerant := Behavior{behaviorName: "tolerant", behavior: Tolerant}

	tests := []struct {
		name  string
		b     Behavior
		setup func(*Entity)
		want  string
	}{
		{
			name:  "AlwaysCooperate after harsh history",
			b:     cooperate,
			setup: func(e *Entity) { e.RecordMoves("CHEAT", "CHEAT") },
			want:  "COOPERATE",
		},
		{
			name:  "AlwaysCheat after cooperation",
			b:     cheat,
			setup: func(e *Entity) { e.RecordMoves("COOPERATE", "COOPERATE") },
			want:  "CHEAT",
		},
		{
			name:  "CopyCat follows recorded opponent move",
			b:     copycat,
			setup: func(e *Entity) { e.RecordMoves("COOPERATE", "CHEAT") },
			want:  "CHEAT",
		},
		{
			name:  "Revenge cooperates before opponent cheats",
			b:     revenge,
			setup: func(e *Entity) { e.RecordMoves("COOPERATE", "COOPERATE") },
			want:  "COOPERATE",
		},
		{
			name:  "Revenge cheats after opponent cheat recorded",
			b:     revenge,
			setup: func(e *Entity) { e.RecordMoves("COOPERATE", "CHEAT") },
			want:  "CHEAT",
		},
		{
			name: "Tolerant cooperates within limit",
			b:    tolerant,
			setup: func(e *Entity) {
				for i := 0; i < 5; i++ {
					e.RecordMoves("COOPERATE", "CHEAT")
				}
			},
			want: "COOPERATE",
		},
		{
			name: "Tolerant cheats beyond limit",
			b:    tolerant,
			setup: func(e *Entity) {
				for i := 0; i < 6; i++ {
					e.RecordMoves("COOPERATE", "CHEAT")
				}
			},
			want: "CHEAT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New("p", tt.b)
			tt.setup(e)
			got := e.Play()
			if got != tt.want {
				t.Fatalf("Play() got %q want %q", got, tt.want)
			}
		})
	}
}

func TestEntityRecordMoves(t *testing.T) {
	tests := []struct {
		name         string
		steps        []struct{ move, opp string }
		wantLast     string
		wantOpp      string
		wantBetrayed int
	}{
		{
			name:         "single cooperate round",
			steps:        []struct{ move, opp string }{{"COOPERATE", "COOPERATE"}},
			wantLast:     "COOPERATE",
			wantOpp:      "COOPERATE",
			wantBetrayed: 0,
		},
		{
			name:         "opponent cheat increments betrayed",
			steps:        []struct{ move, opp string }{{"COOPERATE", "CHEAT"}},
			wantLast:     "COOPERATE",
			wantOpp:      "CHEAT",
			wantBetrayed: 1,
		},
		{
			name: "two opponent cheats stack",
			steps: []struct{ move, opp string }{
				{"COOPERATE", "CHEAT"},
				{"CHEAT", "CHEAT"},
			},
			wantLast:     "CHEAT",
			wantOpp:      "CHEAT",
			wantBetrayed: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New("p", Behavior{behaviorName: "n", behavior: AlwaysCooperate})
			for _, s := range tt.steps {
				e.RecordMoves(s.move, s.opp)
			}
			if e.lastMove != tt.wantLast {
				t.Fatalf("lastMove got %q want %q", e.lastMove, tt.wantLast)
			}
			if e.oppLastMove != tt.wantOpp {
				t.Fatalf("oppLastMove got %q want %q", e.oppLastMove, tt.wantOpp)
			}
			if e.betrayed != tt.wantBetrayed {
				t.Fatalf("betrayed got %d want %d", e.betrayed, tt.wantBetrayed)
			}
		})
	}
}

func TestEntityNewAndGetters(t *testing.T) {
	b := Behavior{behaviorName: "cheater", behavior: AlwaysCheat}
	e := New("alice", b)
	if e.Name != "alice" {
		t.Fatalf("Name got %q want alice", e.Name)
	}
	if e.Score != 0 {
		t.Fatalf("Score got %d want 0", e.Score)
	}
	if e.GetBehaviorName() != "cheater" {
		t.Fatalf("GetBehaviorName got %q want cheater", e.GetBehaviorName())
	}
	gotB := e.GetBehavior()
	if gotB.behaviorName != b.behaviorName {
		t.Fatalf("GetBehavior name got %q want %q", gotB.behaviorName, b.behaviorName)
	}
}

func TestEntityReset(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Entity)
	}{
		{"zero score unchanged", func(*Entity) {}},
		{"nonzero score cleared", func(e *Entity) { e.Score = 42 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New("p", Behavior{behaviorName: "n", behavior: AlwaysCooperate})
			tt.setup(e)
			e.Reset()
			if e.Score != 0 {
				t.Fatalf("Score after Reset got %d want 0", e.Score)
			}
		})
	}
}
