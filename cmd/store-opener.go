package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/iancullinane/prisoner/db"
	storefile "github.com/iancullinane/prisoner/internal/store/file"
	"github.com/iancullinane/prisoner/internal/store/memory"
	storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
func openStores(ctx context.Context, kind string) (stores, func(), error) {
	switch kind {
	case StoreMemory, "":
		return stores{
			players: memory.NewInMemoryPlayerStore(),
			history: memory.NewInMemoryHistoryStore(),
		}, func() {}, nil

	case StoreFile:
		players, closePlayers, err := openPlayerFileStore()
		if err != nil {
			return stores{}, nil, err
		}
		history, closeHistory, err := openHistoryFileStore()
		if err != nil {
			closePlayers()
			return stores{}, nil, err
		}
		return stores{players: players, history: history},
			func() { closeHistory(); closePlayers() }, nil

	case StorePostgres:
		pool, err := openPostgresPool(ctx)
		if err != nil {
			return stores{}, nil, err
		}
		return stores{
			players: storepostgres.NewPlayerStore(pool),
			history: storepostgres.NewHistoryStore(pool),
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
func openPlayerFileStore() (types.PlayerStore, func(), error) {
	f, err := os.OpenFile(playerFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", playerFileName, err)
	}

	store, err := storefile.NewFileSystemPlayerStore(f)
	if err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("new file system player store: %w", err)
	}
	return store, func() { f.Close() }, nil
}

func openHistoryFileStore() (types.HistoryStore, func(), error) {
	f, err := os.OpenFile(historyFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", historyFileName, err)
	}

	store, err := storefile.NewFileSystemHistoryStore(f)
	if err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("new file system player store: %w", err)
	}
	return store, func() { f.Close() }, nil
}

// ===========================

func openPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
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

	// needed for goose
	stdDB := stdlib.OpenDBFromPool(pool)

	// this is kind of the 'explicit' bit needed to use the embedded migrations
	goose.SetBaseFS(db.Migrations)
	goose.SetDialect("postgres")
	if err := goose.Up(stdDB, "migrations"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	fmt.Println("postgres: connected, migrations applied")
	return pool, nil
}
