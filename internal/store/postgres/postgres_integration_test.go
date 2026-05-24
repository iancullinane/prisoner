//go:build integration

package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

var integrationPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	c, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("prisoner"),
		tcpostgres.WithUsername("prisoner"),
		tcpostgres.WithPassword("prisoner"),
		tcpostgres.BasicWaitStrategies(),
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

func TestPostgresPlayerStore_Contract(t *testing.T) {
	testhelpers.RunPlayerStoreContract(t, func() types.PlayerStore {
		// Truncate for a clean slate on each subtest.
		if _, err := integrationPool.Exec(context.Background(), "TRUNCATE players"); err != nil {
			t.Fatalf("could not truncate players table: %v", err)
		}
		return NewPlayerStore(integrationPool)
	})
}
