package file

import (
	"io"

	"github.com/iancullinane/prisoner/api"
	"github.com/iancullinane/prisoner/internal/types"
)

type FilesSystemPlayerStore struct {
	database io.ReadSeeker
}

func NewFilesSystemPlayerStore(database io.ReadSeeker) *FilesSystemPlayerStore {
	return &FilesSystemPlayerStore{database: database}
}

func (f *FilesSystemPlayerStore) GetLeague() []types.Player {
	f.database.Seek(0, io.SeekStart) // Always read from the beginning
	league, _ := api.NewLeague(f.database)
	return league
}

func (f *FilesSystemPlayerStore) GetPlayerScore(name string) int {
	return 5
}
