package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/iancullinane/prisoner/db"
	storefile "github.com/iancullinane/prisoner/internal/store/file"
	"github.com/iancullinane/prisoner/internal/store/memory"
	storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StoreMemory   = "memory"
	StorePostgres = "postgres"
	StoreFile     = "file"
	dbFileName    = "game.db.json"
)

// openHistoryStore lets us optionally set an in-memory, file-based, or database backed
// storage system.
func openHistoryStore(ctx context.Context, kind string) (types.HistoryStore, func(), error) {
	switch kind {
	case StoreMemory, "":
		return memory.NewInMemoryHistoryStore(), nil, nil
	case StorePostgres:
		fmt.Println("postgres history store not implemented")
		return nil, nil, nil
	case StoreFile:
		fmt.Println("file history store not implemented")
		return nil, nil, nil
	default:
		fmt.Println("unknown history store %q: want %q, %q, or %q", kind, StoreMemory, StoreFile, StorePostgres)
		return nil, nil, nil
	}
}

// openPlayerStore lets us optionally set an in-memory, file-based, or database backed
// storage system.
func openPlayerStore(ctx context.Context, kind string) (types.PlayerStore, func(), error) {
	switch kind {
	case StoreMemory, "":
		return memory.NewInMemoryPlayerStore(), nil, nil
	case StorePostgres:
		pool, err := openPostgresPool(ctx)
		if err != nil {
			return nil, nil, err
		}
		return storepostgres.NewPlayerStore(pool), func() { pool.Close() }, nil
	case StoreFile:
		return openFileStore()
	default:
		return nil, nil, fmt.Errorf(
			"unknown store %q: want %q, %q, or %q",
			kind, StoreMemory, StoreFile, StorePostgres,
		)
	}
}

// openFileStore will open or create a new json file to use as storage
// behind the application
func openFileStore() (types.PlayerStore, func(), error) {
	f, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", dbFileName, err)
	}

	store, err := storefile.NewFileSystemPlayerStore(f)
	if err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("new file system player store: %w", err)
	}
	return store, func() { f.Close() }, nil
}

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

	if _, err := pool.Exec(ctx, db.SchemaSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	log.Print("postgres: connected, schema applied")
	return pool, nil
}
