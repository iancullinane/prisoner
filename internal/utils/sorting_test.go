package utils

import (
	"testing"

	"github.com/iancullinane/prisoner/internal/entity"
)

func TestSortByScore(t *testing.T) {
	a := &entity.Entity{Score: 3}
	b := &entity.Entity{Score: 1}
	c := &entity.Entity{Score: 2}
	in := []*entity.Entity{a, b, c}
	got := SortByScore(in)
	if len(got) != 3 {
		t.Fatalf("len got %d", len(got))
	}
	if got[0].Score != 1 || got[1].Score != 2 || got[2].Score != 3 {
		t.Errorf("got scores %d, %d, %d want 1, 2, 3", got[0].Score, got[1].Score, got[2].Score)
	}
	if got[0] != b || got[1] != c || got[2] != a {
		t.Errorf("want same entities reordered, not new allocations")
	}
	if in[0] != a || in[1] != b || in[2] != c {
		t.Error("input slice should not be reordered")
	}
}
