package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"

	"os/exec"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)

	if err != nil {
		log.Fatalf("Failed to start test container: %v", err)
	}

	var mappedPort nat.Port
	for i := 0; i < 10; i++ {
		mappedPort, err = container.MappedPort(ctx, "5432/tcp")
		if err == nil {
			break
		}
		log.Printf("Waiting for mapped port... attempt %d", i+1)
		time.Sleep(1 * time.Second)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container host: %v", err)
	}
	dbURL := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, mappedPort.Port())

	if err := waitForPostgres(host, mappedPort.Port()); err != nil {
		log.Fatalf("PostgreSQL did not become ready: %v", err)
	}

	testDB, err = sql.Open("pgx", dbURL)
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

func waitForPostgres(host, port string) error {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		cmd := exec.Command("pg_isready", "-h", host, "-p", port, "-U", "testuser", "-d", "testdb")
		if err := cmd.Run(); err == nil {
			return nil
		}

		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("PostgreSQL did not become ready after %d retries", maxRetries)
}

func runMigrations(dbURL string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	migrationDir := filepath.Join(workingDir, "../../../migrations/postgres")

	migrationFiles := []string{
		"0001_create_uuid_extension.sql",
		"0002_create_outbox_table.sql",
		"0003_add_indexes.sql",

		"_testdata.sql",
	}

	for _, file := range migrationFiles {
		migrationPath := filepath.Join(migrationDir, file)
		cmd := exec.Command("psql", dbURL, "-f", migrationPath)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
