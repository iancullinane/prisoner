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

	"github.com/iancullinane/prisoner/api"
	"github.com/iancullinane/prisoner/cmd"
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
