package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SeedEmployeeData seeds sample employee data with new fields
func SeedEmployeeData(conn *pgx.Conn) error {
	ctx := context.Background()
	fmt.Println("Seeding employee data with new fields...")

	// Sample employees with new fields
	employees := []struct {
		UserID             uuid.UUID
		Position           string
		PhoneNumber        string
		Address            string
		SalaryBase         float64
		BankName           string
		BankAccountNumber  string
		BankAccountHolder  string
		EmploymentStatus   string
		JobLevel           string
		Gender             string
		Division           string
	}{
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Position:           "Senior Software Engineer",
			PhoneNumber:        "+6281234567801",
			Address:            "Jl. Sudirman No. 1, Jakarta Selatan",
			SalaryBase:         18000000,
			BankName:           "BCA",
			BankAccountNumber:  "1234567801",
			BankAccountHolder:  "Ahmad Sudirman",
			EmploymentStatus:   "PERMANENT",
			JobLevel:           "STAFF",
			Gender:             "MALE",
			Division:           "Information Technology",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Position:           "HR Manager",
			PhoneNumber:        "+6281234567802",
			Address:            "Jl. Thamrin No. 2, Jakarta Pusat",
			SalaryBase:         20000000,
			BankName:           "Mandiri",
			BankAccountNumber:  "1234567802",
			BankAccountHolder:  "Siti Rahayu",
			EmploymentStatus:   "PERMANENT",
			JobLevel:           "MANAGER",
			Gender:             "FEMALE",
			Division:           "Human Resources",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Position:           "IT Supervisor",
			PhoneNumber:        "+6281234567803",
			Address:            "Jl. Gatot Subroto No. 3, Jakarta Selatan",
			SalaryBase:         17000000,
			BankName:           "BCA",
			BankAccountNumber:  "1234567803",
			BankAccountHolder:  "Budi Hartono",
			EmploymentStatus:   "PERMANENT",
			JobLevel:           "SUPERVISOR",
			Gender:             "MALE",
			Division:           "Information Technology",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Position:           "Marketing Staff",
			PhoneNumber:        "+6281234567804",
			Address:            "Jl. Rasuna Said No. 4, Jakarta Selatan",
			SalaryBase:         12000000,
			BankName:           "BNI",
			BankAccountNumber:  "1234567804",
			BankAccountHolder:  "Dewi Sartika",
			EmploymentStatus:   "PROBATION",
			JobLevel:           "STAFF",
			Gender:             "FEMALE",
			Division:           "Marketing",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Position:           "Finance Manager",
			PhoneNumber:        "+6281234567805",
			Address:            "Jl. Sudirman Kav 5, Jakarta Pusat",
			SalaryBase:         22000000,
			BankName:           "Mandiri",
			BankAccountNumber:  "1234567805",
			BankAccountHolder:  "Eko Prasetyo",
			EmploymentStatus:   "PERMANENT",
			JobLevel:           "MANAGER",
			Gender:             "MALE",
			Division:           "Finance",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
			Position:           "Contract Developer",
			PhoneNumber:        "+6281234567806",
			Address:            "Jl. Fatmawati No. 6, Jakarta Selatan",
			SalaryBase:         15000000,
			BankName:           "BCA",
			BankAccountNumber:  "1234567806",
			BankAccountHolder:  "Feri Irawan",
			EmploymentStatus:   "CONTRACT",
			JobLevel:           "STAFF",
			Gender:             "MALE",
			Division:           "Information Technology",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"),
			Position:           "Operations Supervisor",
			PhoneNumber:        "+6281234567807",
			Address:            "Jl. Panglima Polim No. 7, Jakarta Selatan",
			SalaryBase:         16000000,
			BankName:           "BRI",
			BankAccountNumber:  "1234567807",
			BankAccountHolder:  "Gita Pertiwi",
			EmploymentStatus:   "PERMANENT",
			JobLevel:           "SUPERVISOR",
			Gender:             "FEMALE",
			Division:           "Operations",
		},
		{
			UserID:             uuid.MustParse("550e8400-e29b-41d4-a716-446655440008"),
			Position:           "CEO",
			PhoneNumber:        "+6281234567808",
			Address:            "Jl. Sudirman Kav 50, Jakarta Pusat",
			SalaryBase:         50000000,
			BankName:           "BCA",
			BankAccountNumber:  "1234567808",
			BankAccountHolder:  "Hendra Wijaya",
			EmploymentStatus:   "PERMANENT",
			JobLevel:           "CEO",
			Gender:             "MALE",
			Division:           "General",
		},
	}

	// Insert employees
	for _, emp := range employees {
		query := `
			INSERT INTO employees (
				user_id, position, phone_number, address, salary_base,
				bank_name, bank_account_number, bank_account_holder,
				employment_status, job_level, gender, division
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (user_id) DO UPDATE SET
				position = EXCLUDED.position,
				phone_number = EXCLUDED.phone_number,
				address = EXCLUDED.address,
				salary_base = EXCLUDED.salary_base,
				bank_name = EXCLUDED.bank_name,
				bank_account_number = EXCLUDED.bank_account_number,
				bank_account_holder = EXCLUDED.bank_account_holder,
				employment_status = EXCLUDED.employment_status,
				job_level = EXCLUDED.job_level,
				gender = EXCLUDED.gender,
				division = EXCLUDED.division
		`

		_, err := conn.Exec(
			ctx,
			query,
			emp.UserID, emp.Position, emp.PhoneNumber, emp.Address, emp.SalaryBase,
			emp.BankName, emp.BankAccountNumber, emp.BankAccountHolder,
			emp.EmploymentStatus, emp.JobLevel, emp.Gender, emp.Division,
		)

		if err != nil {
			log.Printf("Error inserting employee %s: %v", emp.UserID, err)
			return err
		}
	}

	fmt.Println("✓ Employee data seeded successfully")
	return nil
}

// SeedDivisions seeds division reference data
func SeedDivisions(conn *pgx.Conn) error {
	ctx := context.Background()
	fmt.Println("Seeding divisions...")

	divisions := []struct {
		Name        string
		Code        string
		Description string
	}{
		{"Information Technology", "IT", "IT and Software Development"},
		{"Human Resources", "HR", "HR and Recruitment"},
		{"Finance", "FIN", "Finance and Accounting"},
		{"Marketing", "MKT", "Marketing and Sales"},
		{"Operations", "OPS", "Operations and Logistics"},
		{"General", "GEN", "General Administration"},
	}

	for _, div := range divisions {
		query := `
			INSERT INTO divisions (name, code, description)
			VALUES ($1, $2, $3)
			ON CONFLICT (name) DO NOTHING
		`

		_, err := conn.Exec(ctx, query, div.Name, div.Code, div.Description)
		if err != nil {
			log.Printf("Error inserting division %s: %v", div.Name, err)
			return err
		}
	}

	fmt.Println("✓ Divisions seeded successfully")
	return nil
}
