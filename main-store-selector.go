package main

import (
	"context"
	"fmt"
	"log"
	"os"

	storefile "github.com/iancullinane/prisoner/internal/store/file"
	"github.com/iancullinane/prisoner/internal/store/memory"
	storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
	"github.com/iancullinane/prisoner/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newMemoryStore() types.PlayerStore {
	return memory.NewInMemoryPlayerStore()
}

func newPostgresStore(pool *pgxpool.Pool) types.PlayerStore {
	return storepostgres.NewPlayerStore(pool)
}

func newFileStore() (types.PlayerStore, func(), error) {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", dbFileName, err)
	}

	store, err := storefile.NewFileSystemPlayerStore(db)
	if err != nil {
		return nil, nil, fmt.Errorf("new file system player store: %w", err)
	}
	return store, func() { db.Close() }, nil
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

	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	log.Print("postgres: connected, schema applied")
	return pool, nil
}
