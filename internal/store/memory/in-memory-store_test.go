package memory

import (
	"testing"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

// MARK: Integration Test
// =========================================
func TestInMemoryPlayerStore_Contract(t *testing.T) {
	testhelpers.RunPlayerStoreContract(t, func(t *testing.T) types.PlayerStore {
		return NewInMemoryPlayerStore()
	})
}

func TestInMemoryHistoryStore_Contract(t *testing.T) {
	testhelpers.RunHistoryStoreContract(t, func(t *testing.T) types.HistoryStore {
		return NewInMemoryHistoryStore()
	})
}
