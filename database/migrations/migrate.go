package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Database connection string
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/hris_db?sslmode=disable"
	}

	conn, err := pgx.Connect(nil, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(nil)

	fmt.Println("Running migration 000005_add_employee_fields...")

	// Read migration file
	migrationSQL := `
-- Add employment_status column
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS employment_status VARCHAR(20) DEFAULT 'PROBATION'
CHECK (employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION'));

-- Add job_level column
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS job_level VARCHAR(20) DEFAULT 'STAFF'
CHECK (job_level IN ('CEO', 'MANAGER', 'SUPERVISOR', 'STAFF'));

-- Add gender column
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS gender VARCHAR(10)
CHECK (gender IN ('MALE', 'FEMALE'));

-- Add division column
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS division VARCHAR(100) DEFAULT 'GENERAL';

-- Add comments for documentation
COMMENT ON COLUMN employees.employment_status IS 'Employment status: PERMANENT, CONTRACT, PROBATION';
COMMENT ON COLUMN employees.job_level IS 'Job level: CEO, MANAGER, SUPERVISOR, STAFF';
COMMENT ON COLUMN employees.gender IS 'Gender: MALE, FEMALE';
COMMENT ON COLUMN employees.division IS 'Department/Division: IT, HR, FINANCE, MARKETING, OPERATIONS, GENERAL';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_employees_employment_status ON employees(employment_status);
CREATE INDEX IF NOT EXISTS idx_employees_job_level ON employees(job_level);
CREATE INDEX IF NOT EXISTS idx_employees_division ON employees(division);
CREATE INDEX IF NOT EXISTS idx_employees_gender ON employees(gender);

-- Update existing records to have default values
UPDATE employees
SET
    employment_status = 'PERMANENT',
    job_level = 'STAFF',
    division = 'GENERAL'
WHERE employment_status IS NULL
   OR job_level IS NULL
   OR division IS NULL;
`

	// Execute migration
	_, err = conn.Exec(nil, migrationSQL)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	fmt.Println("✓ Migration completed successfully!")
	fmt.Println("✓ Added employment_status, job_level, gender, and division columns")
	fmt.Println("✓ Created indexes for better query performance")
	fmt.Println("✓ Updated existing records with default values")
}
