package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"hris/database/seeder"
	"hris/internal/user/repository"
)

func main() {
	ctx := context.Background()

	// Connect to database
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/fiber_app?sslmode=disable")
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
