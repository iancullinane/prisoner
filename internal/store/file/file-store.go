package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

var ErrInteractionNotFound = errors.New("interaction not found")

const HistoryFile = "history.db.json"

type FileSystemHistoryStore struct {
	database *json.Encoder
	history  types.History
}

func NewFileSystemHistoryStore(file *os.File) (*FileSystemHistoryStore, error) {
	file.Seek(0, io.SeekStart)

	err := initialiseHistoryDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising history db file, %v", err)
	}

	history, err := types.NewHistory(file)
	if err != nil {
		return nil, err
	}

	return &FileSystemHistoryStore{
		database: json.NewEncoder(&tape{file}),
		history:  history,
	}, nil
}

func (fh *FileSystemHistoryStore) GetHistory() (types.History, error) {
	return fh.history, nil
}

func (fh *FileSystemHistoryStore) GetInteraction(interactionID uuid.UUID) (types.Interaction, error) {
	interaction := fh.history.Find(interactionID)
	if interaction == nil {
		return types.Interaction{}, ErrInteractionNotFound
	}
	return *interaction, nil
}

func (f *FileSystemHistoryStore) RecordInteraction(interaction types.Interaction) error {
	f.history = append(f.history, interaction)

	if err := f.database.Encode(f.history); err != nil {
		return fmt.Errorf("encoding history to file: %w", err)
	}

	return nil
}

// MARK: Player Store
// ==================

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

// MARK: Initialize History DB
// ===========================

func initialiseHistoryDBFile(file *os.File) error {
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
