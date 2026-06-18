package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"

    storefile "github.com/iancullinane/prisoner/internal/store/file"
    "github.com/iancullinane/prisoner/internal/store/memory"
    storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
    "github.com/iancullinane/prisoner/internal/types"
    "github.com/jackc/pgx/v5/pgxpool"
)

const (
	StoreMemory     = "memory"
	StorePostgres   = "postgres"
	StoreFile       = "file"
	playerFileName  = "player.db.json"
	historyFileName = "history.db.json"
)

// stores bundles the backends a command needs so they can share one
// underlying connection (e.g. a single postgres pool).
type stores struct {
	players types.PlayerStore
	history types.HistoryStore
}

// openStores selects an in-memory, file-based, or database backed storage
// system and constructs both the player and history stores from a single
// connection. The returned cleanup is always safe to call.
func openStores(ctx context.Context, kind string, logger *slog.Logger) (stores, func(), error) {
	switch kind {
	case StoreMemory, "":
		return stores{
			players: memory.NewInMemoryPlayerStore(logger),
			history: memory.NewInMemoryHistoryStore(logger),
		}, func() {}, nil

	case StoreFile:
		players, closePlayers, err := openPlayerFileStore(logger)
		if err != nil {
			return stores{}, nil, err
		}
		history, closeHistory, err := openHistoryFileStore(logger)
		if err != nil {
			closePlayers()
			return stores{}, nil, err
		}
		return stores{players: players, history: history},
			func() { closeHistory(); closePlayers() }, nil

	case StorePostgres:
		pool, err := openPostgresPool(ctx, logger)
		if err != nil {
			return stores{}, nil, err
		}
		return stores{
			players: storepostgres.NewPlayerStore(logger, pool),
			history: storepostgres.NewHistoryStore(logger, pool),
		}, func() { pool.Close() }, nil

	default:
		return stores{}, nil, fmt.Errorf(
			"unknown store %q: want %q, %q, or %q",
			kind, StoreMemory, StoreFile, StorePostgres,
		)
	}
}

// openPlayerFileStore will open or create a new json file to use as storage
// behind the application
func openPlayerFileStore(logger *slog.Logger) (types.PlayerStore, func(), error) {
	f, err := os.OpenFile(playerFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", playerFileName, err)
	}

	store, err := storefile.NewFileSystemPlayerStore(logger, f)
	if err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("new file system player store: %w", err)
	}
	return store, func() { f.Close() }, nil
}

func openHistoryFileStore(logger *slog.Logger) (types.HistoryStore, func(), error) {
	f, err := os.OpenFile(historyFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", historyFileName, err)
	}

	store, err := storefile.NewFileSystemHistoryStore(logger, f)
	if err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("new file system player store: %w", err)
	}
	return store, func() { f.Close() }, nil
}

// ===========================

func openPostgresPool(ctx context.Context, logger *slog.Logger) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required for postgres store")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("postgres ping: %w", err)
    }

    logger.Info("postgres connected")
    return pool, nil
}
