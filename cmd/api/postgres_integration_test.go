//go:build integration

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	storepostgres "github.com/iancullinane/prisoner/internal/store/postgres"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var integrationPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	c, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("prisoner"),
		tcpostgres.WithUsername("prisoner"),
		tcpostgres.WithPassword("prisoner"),
	)
	if err != nil {
		panic(err)
	}

	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		panic(err)
	}

	schemaPath, err := findRepoFile("db/schema.sql")
	if err != nil {
		panic(err)
	}
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		panic(err)
	}
	if _, err := pool.Exec(ctx, string(schemaSQL)); err != nil {
		panic(err)
	}

	integrationPool = pool
	code := m.Run()

	pool.Close()
	if err := c.Terminate(context.Background()); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func findRepoFile(rel string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return filepath.Join(dir, rel), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func TestRecordingWinsAndRetrievingThem_Postgres(t *testing.T) {
	store := storepostgres.NewPlayerStore(integrationPool)
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest("players", player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest("players", player))
		assertResponseStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertResponseStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{Name: "Peppesadr", Wins: 4},
		}
		assertLeague(t, got, want)
	})
}
