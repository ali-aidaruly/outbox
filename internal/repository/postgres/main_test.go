package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	postgres2 "github.com/testcontainers/testcontainers-go/modules/postgres"
	"os/exec"
)

// Global test variables
var TestDB *sql.DB
var testContainer *postgres2.PostgresContainer

// TestMain sets up the test container before running any tests.
func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgres2.Run(ctx, "postgres:16-alpine")
	if err != nil {
		log.Fatalf("Failed to start test container: %v", err)
	}
	testContainer = container

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	TestDB, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate test container: %v", err)
	}

	os.Exit(code)
}

func runMigrations(dbURL string) error {
	migrationFiles, err := filepath.Glob("migrations/postgres/*.sql")
	if err != nil {
		return err
	}

	for _, file := range migrationFiles {
		cmd := exec.Command("psql", dbURL, "-f", file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
