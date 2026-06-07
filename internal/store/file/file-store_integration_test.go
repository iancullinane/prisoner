package file

import (
	"testing"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

func TestFileSystemPlayerStore_Contract(t *testing.T) {
	testhelpers.RunPlayerStoreContract(t, func(t *testing.T) types.PlayerStore {
		f, cleanup := createTempFile(t, "[]")
		t.Cleanup(cleanup)
		store, err := NewFileSystemPlayerStore(testhelpers.NoopLogger(), f)
		if err != nil {
			t.Fatalf("could not create store: %v", err)
		}
		return store
	})
}

func TestFileSystemHistoryStore_Contract(t *testing.T) {
	testhelpers.RunHistoryStoreContract(t, func(t *testing.T) types.HistoryStore {
		f, cleanup := createTempFile(t, "[]")
		t.Cleanup(cleanup)
		store, err := NewFileSystemHistoryStore(testhelpers.NoopLogger(), f)
		if err != nil {
			t.Fatalf("could not create store: %v", err)
		}
		return store
	})
}
