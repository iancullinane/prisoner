//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/iancullinane/prisoner/db"
	"github.com/iancullinane/prisoner/internal/store/testhelpers"
	"github.com/iancullinane/prisoner/internal/types"
)

var integrationPool *pgxpool.Pool

// TestMain is a special name used for "setting" up your tests,
// which somehow I had no idea about until now....
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

	// Run goose migrations — In production, migrations are run "manually"
	// via the cmd/migrate command, this just keeps the test container up
	// to date during local development, so this is isolated
	stdDB := stdlib.OpenDBFromPool(pool)
	if err := runMigrations(stdDB); err != nil {
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

func runMigrations(stdDB *sql.DB) error {
	goose.SetBaseFS(db.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(stdDB, "migrations")
}

func TestPostgresPlayerStore_Contract(t *testing.T) {
	testhelpers.RunPlayerStoreContract(t, func(t *testing.T) types.PlayerStore {
		// CASCADE is required because interactions has a FK referencing players.
		if _, err := integrationPool.Exec(context.Background(), "TRUNCATE players CASCADE"); err != nil {
			t.Fatalf("could not truncate players table: %v", err)
		}
		return NewPlayerStore(testhelpers.NoopLogger(), integrationPool)
	})
}

func TestPostgresHistoryStore_Contract(t *testing.T) {
	testhelpers.RunHistoryStoreContract(t, func(t *testing.T) types.HistoryStore {
		ctx := context.Background()
		// Player store contract leaves "Alice" in the table; reset everything first.
		if _, err := integrationPool.Exec(ctx, "TRUNCATE players, interactions CASCADE"); err != nil {
			t.Fatalf("could not reset tables: %v", err)
		}
		if err := seedTestPlayers(ctx, integrationPool); err != nil {
			t.Fatalf("seed test players: %v", err)
		}
		return NewHistoryStore(testhelpers.NoopLogger(), integrationPool)
	})
}

func seedTestPlayers(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO players (id, name) VALUES
			($1, $2),
			($3, $4)
		ON CONFLICT (id) DO NOTHING`,
		testhelpers.TestPlayerOneID, testhelpers.TestPlayerOneName,
		testhelpers.TestPlayerTwoID, testhelpers.TestPlayerTwoName,
	)
	return err
}
