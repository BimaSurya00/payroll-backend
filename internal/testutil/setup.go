package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"hris/config"
)

// NewTestDB creates a connection to test database
func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Connect to test database
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf(
		"postgres://test:%s@%s:%s@localhost:5432/hris_test?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
	))
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	t.Cleanup(func() {
		// Cleanup tables for test isolation
		CleanupTables(t, pool,
			"employees",
			"leave_requests",
			"leave_balances",
			"payroll_items",
			"payrolls",
			"attendance_records",
			"overtime_requests",
			"audit_logs",
		)
	})

	return pool
}

// CleanupTables truncates specified tables for test isolation
func CleanupTables(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()

	for _, table := range tables {
		_, err := pool.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Errorf("Failed to truncate table %s: %v", table, err)
		}
	}
}
