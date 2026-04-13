package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"hris/database/seeder"
	"hris/internal/user/repository"
)

func main() {
	ctx := context.Background()

	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "postgres")
	dbPass := getEnv("POSTGRES_PASSWORD", "postgres")
	dbName := getEnv("POSTGRES_DATABASE", "fiber_app")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	fmt.Println("Connected to database")

	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)

	// Run seeders
	superUserSeeder := seeder.NewSuperUserSeeder(userRepo)
	if err := superUserSeeder.Seed(ctx); err != nil {
		log.Fatalf("Failed to seed super user: %v", err)
	}

	fmt.Println("\n✅ All seeders completed successfully!")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
