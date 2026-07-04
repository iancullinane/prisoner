package types

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/pkg/prisoner"
)

func TestNewInteractionSetsPlayedAt(t *testing.T) {
	before := time.Now()
	interaction := NewInteraction(uuid.New(), uuid.New(), prisoner.Cooperate, prisoner.Betray)
	after := time.Now()

	if interaction.PlayedAt.IsZero() {
		t.Fatal("expected PlayedAt to be set, got zero value")
	}
	if interaction.PlayedAt.Before(before) || interaction.PlayedAt.After(after) {
		t.Errorf("expected PlayedAt within [%v, %v], got %v", before, after, interaction.PlayedAt)
	}
}
