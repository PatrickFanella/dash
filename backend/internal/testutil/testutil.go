package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/patrickfanella/dash/backend/internal/database"
	"github.com/patrickfanella/dash/backend/internal/models"
)

// SetupPool reads TEST_DATABASE_URL, runs migrations, and returns a connected
// pool + queries. Returns nil if TEST_DATABASE_URL is not set. The caller
// should defer pool.Close() and call os.Exit(m.Run()).
func SetupPool(migrationsPath string) (*pgxpool.Pool, *models.Queries) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		fmt.Println("TEST_DATABASE_URL not set, skipping integration tests")
		return nil, nil
	}

	if err := database.RunMigrations(dbURL, migrationsPath); err != nil {
		fmt.Fprintf(os.Stderr, "migrations failed: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := database.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}

	return pool, models.New(pool)
}

// SkipIfNoDB skips the test if no database pool is available.
func SkipIfNoDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	if pool == nil {
		t.Skip("TEST_DATABASE_URL not set")
	}
}

// TruncateAll truncates all application tables with CASCADE.
func TruncateAll(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"TRUNCATE service_section_mappings, services, sections CASCADE")
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

// FormatUUID converts a [16]byte UUID to its standard string representation.
func FormatUUID(b [16]byte) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
