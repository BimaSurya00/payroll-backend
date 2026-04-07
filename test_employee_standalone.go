// This is a test utility file for employee conversion testing
// Run with: go run test_employee.go
// Note: This file should not be compiled with main.go
// Use: go run test_employee.go (standalone)

//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/employee/helper"
	"hris/internal/employee/repository"
)

func main() {
	TestEmployeeConversion()
}

func TestEmployeeConversion() {
	// Connect to database
	connString := "postgres://postgres:postgres@localhost:5432/fiber_app"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test repository FindAll
	repo := repository.NewEmployeeRepository(pool)

	fmt.Println("Testing FindAll...")
	employees, total, err := repo.FindAll(ctx, 1, 15, "")
	if err != nil {
		log.Fatalf("FindAll failed: %v", err)
	}

	fmt.Printf("✓ Found %d employees (total: %d)\n", len(employees), total)

	// Test converter
	fmt.Println("\nTesting converter...")
	for i, emp := range employees {
		fmt.Printf("\n[%d] Converting employee: %s\n", i+1, emp.Position)
		fmt.Printf("  - EmploymentStatus: '%s'\n", emp.EmploymentStatus)
		fmt.Printf("  - JobLevel: '%s'\n", emp.JobLevel)
		fmt.Printf("  - Gender: '%s'\n", emp.Gender)
		fmt.Printf("  - Division: '%s'\n", emp.Division)

		response := helper.ToEmployeeResponseFromDB(&emp)
		if response == nil {
			fmt.Printf("  ❌ Response is NIL!\n")
			continue
		}

		fmt.Printf("  ✓ Response created: %s - %s\n", response.EmploymentStatus, response.JobLevel)
	}

	// Test ToEmployeeResponses
	fmt.Println("\n\nTesting ToEmployeeResponses...")
	responses := helper.ToEmployeeResponses(employees)
	fmt.Printf("✓ Created %d responses from %d employees\n", len(responses), len(employees))

	fmt.Println("\n✓ All tests passed!")
}
