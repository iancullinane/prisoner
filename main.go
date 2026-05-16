/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iancullinane/prisoner/cmd"
	"github.com/iancullinane/prisoner/cmd/api"
	"github.com/iancullinane/prisoner/internal/store/memory"
	storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
)

//go:embed db/schema.sql
var schemaSQL string

func main() {

	runServer := flag.Bool("server", false, "set to run the server, otherwise execute cli")

	flag.Parse()

	if *runServer {

		var server *api.PlayerServer
		if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
			ctx := context.Background()
			pool, err := pgxpool.New(ctx, dsn)
			if err != nil {
				log.Fatal(err)
			}
			defer pool.Close()
			if err := pool.Ping(ctx); err != nil {
				log.Fatalf("postgres ping: %v", err)
			}
			if _, err := pool.Exec(ctx, schemaSQL); err != nil {
				log.Fatalf("apply schema: %v", err)
			}
			log.Print("postgres: connected, schema applied, listening :5001")
			server = api.NewPlayerServer(storepostgres.NewPlayerStore(pool))
		} else {
			server = api.NewPlayerServer(memory.NewInMemoryPlayerStore())
		}
		log.Fatal(http.ListenAndServe(":5001", server))
	}

	cmd.Execute()
}
