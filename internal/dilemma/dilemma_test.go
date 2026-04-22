package dilemma

import (
	"testing"

	"github.com/iancullinane/prisoner/internal/entity"
)

func TestDilemma_Once(t *testing.T) {

	behaviorFactory := entity.NewBehaviorFactory()
	var tests = []struct {
		name    string
		entity1 *entity.Entity
		entity2 *entity.Entity
		want1   int32
		want2   int32
	}{
		{
			"cooperate/cooperate",
			entity.New("cooperate", behaviorFactory.GetBehaviorByName("niceguy")),
			entity.New("cheat", behaviorFactory.GetBehaviorByName("cheater")),
			-1,
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Once(tt.entity1, tt.entity2)

			if tt.entity1.Score != tt.want1 {
				t.Errorf("entity1 score wrong %d want %d", tt.entity1.Score, tt.want1)
			}
			if tt.entity2.Score != tt.want2 {
				t.Errorf("entity2 score wrong %d want %d", tt.entity2.Score, tt.want2)
			}
		})
	}
}
