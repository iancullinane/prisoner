package memory

import (
	"testing"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

func TestInMemoryPlayerStore_Contract(t *testing.T) {
	testhelpers.RunPlayerStoreContract(t, func() types.PlayerStore {
		return NewInMemoryPlayerStore()
	})
}
