/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
    "context"
    "fmt"
    "log/slog"
    "os"

    "database/sql"

    "github.com/iancullinane/prisoner/db"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5/stdlib"
    "github.com/pressly/goose/v3"
    "github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Database migration commands",
    Long:  `Run database migrations up, down, or check status`,
}

var migrateUpCmd = &cobra.Command{
    Use:   "up",
    Short: "Run all pending migrations",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runMigration(func(stdDB *sql.DB) error {
            return goose.Up(stdDB, "migrations")
        })
    },
}

var migrateDownCmd = &cobra.Command{
    Use:   "down",
    Short: "Roll back the last migration",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runMigration(func(stdDB *sql.DB) error {
            return goose.Down(stdDB, "migrations")
        })
    },
}

var migrateStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show migration status",
    RunE: func(cmd *cobra.Command, args []string) error {
        return runMigration(func(stdDB *sql.DB) error {
            return goose.Status(stdDB, "migrations")
        })
    },
}

func runMigration(fn func(*sql.DB) error) error {
    ctx := context.Background()
    log := logger.With(slog.String("service", "migrate"))

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        return fmt.Errorf("DATABASE_URL is required")
    }

    pool, err := pgxpool.New(ctx, dsn)
    if err != nil {
        return fmt.Errorf("connect to database: %w", err)
    }
    defer pool.Close()

    if err := pool.Ping(ctx); err != nil {
        return fmt.Errorf("ping database: %w", err)
    }

    stdDB := stdlib.OpenDBFromPool(pool)
    defer stdDB.Close()

    goose.SetBaseFS(db.Migrations)
    goose.SetDialect("postgres")

    if err := fn(stdDB); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    log.Info("migration completed successfully")
    return nil
}

func init() {
    rootCmd.AddCommand(migrateCmd)
    migrateCmd.AddCommand(migrateUpCmd)
    migrateCmd.AddCommand(migrateDownCmd)
    migrateCmd.AddCommand(migrateStatusCmd)
}
