package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/iancullinane/prisoner/internal/types"
)

var ErrInteractionNotFound = errors.New("interaction not found")
var ErrNotImplemented = errors.New("file store not implemented")

const seedPlayerName = "Luigi"

var seedPlayerID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

// MARK: History Store
// ===================
const HistoryFile = "history.db.json"

type FileSystemHistoryStore struct {
	logger   *slog.Logger
	database *json.Encoder
	history  types.History
}

func NewFileSystemHistoryStore(logger *slog.Logger, file *os.File) (*FileSystemHistoryStore, error) {
	file.Seek(0, io.SeekStart)

	err := initialiseHistoryDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising history db file, %v", err)
	}

	history, err := types.NewHistory(file)
	if err != nil {
		return nil, err
	}

	logger = logger.With(slog.String("component", "file-history-store"))
	logger.Debug("opened history store",
		slog.String("file", file.Name()),
		slog.Int("history_len", len(history)),
	)

	return &FileSystemHistoryStore{
		logger:   logger,
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

	f.logger.Debug("recorded interaction",
		slog.String("id", interaction.ID.String()),
		slog.Int("history_len", len(f.history)),
	)
	return nil
}

func (f *FileSystemHistoryStore) GetHistoryByPlayerID(playerID uuid.UUID) (types.History, error) {
	historyByPlayer := make(types.History, 0)
	for _, interaction := range f.history {
		if interaction.PlayerA == playerID || interaction.PlayerB == playerID {
			historyByPlayer = append(historyByPlayer, interaction)
		}
	}
	return historyByPlayer, nil
}

// MARK: Player Store
// ==================

type FileSystemStore struct {
	logger   *slog.Logger
	database *json.Encoder
	players  types.Players
}

func NewFileSystemPlayerStore(logger *slog.Logger, file *os.File) (*FileSystemStore, error) {
	file.Seek(0, io.SeekStart)

	err := initialisePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	players, err := types.NewPlayers(file)
	if err != nil {
		return nil, err
	}

	logger = logger.With(slog.String("component", "file-player-store"))
	logger.Debug("opened player store",
		slog.String("file", file.Name()),
		slog.Int("player_count", len(players)),
	)

	store := &FileSystemStore{
		logger:   logger,
		database: json.NewEncoder(&tape{file}),
		players:  players,
	}

	if store.players.FindByName(seedPlayerName) == nil {
		store.players = append(store.players, *types.NewPlayerWithID(seedPlayerID, seedPlayerName))
		if err := store.database.Encode(store.players); err != nil {
			return nil, fmt.Errorf("seeding %s player to file: %w", seedPlayerName, err)
		}
	}

	return store, nil

}

func (f *FileSystemStore) GetOrCreatePlayer(name string) (types.Player, error) {
	player := f.players.FindByName(name)
	if player != nil {
		return *player, nil
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return types.Player{}, fmt.Errorf("generate UUID: %w", err)
	}

	newPlayer := types.NewPlayerWithID(id, name)

	f.players = append(f.players, *newPlayer)

	if err := f.database.Encode(f.players); err != nil {
		return types.Player{}, fmt.Errorf("encoding players to file: %w", err)
	}

	f.logger.Debug("created player",
		slog.String("name", newPlayer.Name),
		slog.String("id", newPlayer.ID.String()),
	)
	return *newPlayer, nil
}

func (f *FileSystemStore) GetAllPlayers() (types.Players, error) {
	players := f.players.GetAllPlayers()
	return players, nil
}

func (f *FileSystemStore) GetPlayerByID(id uuid.UUID) (types.Player, error) {
	player := f.players.FindByID(id)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}

	return *player, nil
}

func (f *FileSystemStore) GetPlayerByName(name string) (types.Player, error) {
	player := f.players.FindByName(name)
	if player == nil {
		return types.Player{}, types.ErrPlayerNotFound
	}
	return *player, nil
}

func (f *FileSystemStore) GetRandomPlayer() (types.Player, error) {
	return f.players.GetRandomPlayer()
}

func (f *FileSystemStore) GetRandomPlayerExcept(exceptID uuid.UUID) (types.Player, error) {
	return f.players.GetRandomPlayerExcept(exceptID)
}

// MARK: Initialize Player DB
// ===========================

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
