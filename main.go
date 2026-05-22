/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iancullinane/prisoner/api"
	"github.com/iancullinane/prisoner/cmd"
	storefile "github.com/iancullinane/prisoner/internal/store/file"
	"github.com/iancullinane/prisoner/internal/store/memory"
	storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
	"github.com/iancullinane/prisoner/internal/types"
)

//go:embed db/schema.sql
var schemaSQL string

const (
	storeMemory   = "memory"
	storePostgres = "postgres"
	storeFile     = "file"
	dbFileName    = "game.db.json"
)

func main() {
	runServer := flag.Bool("server", false, "set to run the server, otherwise execute cli")
	storeFlag := flag.String("store", storeMemory, "player store backend: memory, file, or postgres")

	flag.Parse()

	if *runServer {
		store, cleanup, err := newPlayerStore(*storeFlag)
		if err != nil {
			log.Fatal(err)
		}
		if cleanup != nil {
			defer cleanup()
		}

		server := api.NewPlayerServer(store)
		log.Printf("%s store: listening on :5001", *storeFlag)
		log.Fatal(http.ListenAndServe(":5001", server))
	}

	cmd.Execute()
}

func newPlayerStore(kind string) (types.PlayerStore, func(), error) {
	switch kind {
	case storeMemory:
		return newMemoryStore(), nil, nil
	case storePostgres:
		ctx := context.Background()
		pool, err := openPostgresPool(ctx)
		if err != nil {
			return nil, nil, err
		}
		return newPostgresStore(pool), func() { pool.Close() }, nil
	case storeFile:
		return newFileStore()
	default:
		return nil, nil, fmt.Errorf("unknown store %q: want %q, %q, or %q", kind, storeMemory, storeFile, storePostgres)
	}
}

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
