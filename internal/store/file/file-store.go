package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/iancullinane/prisoner/internal/types"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   types.League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0, io.SeekStart)

	err := initialisePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := types.NewLeague(file)
	if err != nil {
		return nil, err
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil

}

func (f *FileSystemPlayerStore) GetLeague() (types.League, error) {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league, nil
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, error) {

	player := f.league.Find(name)

	if player != nil {
		return player.Wins, nil
	}

	return 0, nil
}

func (f *FileSystemPlayerStore) RecordWin(name string) error {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, types.Player{Name: name, Wins: 1})
	}

	if err := f.database.Encode(f.league); err != nil {
		return fmt.Errorf("encoding league to file: %w", err)
	}
	return nil
}

// file_system_store.go
func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}

	return nil
}
