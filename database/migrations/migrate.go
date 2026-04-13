package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	conn, err := pgx.Connect(nil, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(nil)

	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		migrationDir = "/migrations"
	}

	ensureMigrationsTable(conn)
	applied := getAppliedMigrations(conn)

	files, err := filepath.Glob(filepath.Join(migrationDir, "*.up.sql"))
	if err != nil {
		log.Fatalf("Failed to read migration files: %v\n", err)
	}

	sort.Strings(files)

	count := 0
	for _, file := range files {
		name := filepath.Base(file)
		if applied[name] {
			fmt.Printf("Skipping %s (already applied)\n", name)
			continue
		}

		sql, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read %s: %v\n", name, err)
		}

		tx, err := conn.Begin(nil)
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v\n", err)
		}

		_, err = tx.Exec(nil, string(sql))
		if err != nil {
			tx.Rollback(nil)
			log.Fatalf("Migration %s failed: %v\n", name, err)
		}

		_, err = tx.Exec(nil, "INSERT INTO schema_migrations (filename) VALUES ($1)", name)
		if err != nil {
			tx.Rollback(nil)
			log.Fatalf("Failed to record migration %s: %v\n", name, err)
		}

		if err := tx.Commit(nil); err != nil {
			log.Fatalf("Failed to commit migration %s: %v\n", name, err)
		}

		fmt.Printf("Applied %s\n", name)
		count++
	}

	fmt.Printf("Done. %d migration(s) applied.\n", count)
}

func ensureMigrationsTable(conn *pgx.Conn) {
	_, err := conn.Exec(nil, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v\n", err)
	}
}

func getAppliedMigrations(conn *pgx.Conn) map[string]bool {
	rows, err := conn.Query(nil, "SELECT filename FROM schema_migrations")
	if err != nil {
		return make(map[string]bool)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		applied[strings.TrimSpace(name)] = true
	}
	return applied
}
