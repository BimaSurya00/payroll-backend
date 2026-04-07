package entity

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	CompanyID         uuid.UUID  `json:"companyId" db:"company_id"`
	UserID            uuid.UUID  `json:"userId" db:"user_id"`
	FullName          string     `json:"fullName" db:"full_name"`
	Position          string     `json:"position" db:"position"`
	PhoneNumber       string     `json:"phoneNumber" db:"phone_number"`
	Address           string     `json:"address" db:"address"`
	SalaryBase        float64    `json:"salaryBase" db:"salary_base"`
	JoinDate          time.Time  `json:"joinDate" db:"join_date"`
	BankName          string     `json:"bankName" db:"bank_name"`
	BankAccountNumber string     `json:"bankAccountNumber" db:"bank_account_number"`
	BankAccountHolder string     `json:"bankAccountHolder" db:"bank_account_holder"`
	ScheduleID        *uuid.UUID `json:"scheduleId" db:"schedule_id"`
	EmploymentStatus  string     `json:"employmentStatus" db:"employment_status"`
	JobLevel          string     `json:"jobLevel" db:"job_level"`
	Gender            string     `json:"gender" db:"gender"`
	Division          string     `json:"division" db:"division"`
	DepartmentID      *uuid.UUID `json:"departmentId" db:"department_id"`
	DeletedAt         *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`

	// Relations
	User interface{} `json:"user,omitempty" db:"-"` // To hold joined User data
}

type EmployeeWithUser struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	CompanyID         uuid.UUID  `json:"companyId" db:"company_id"`
	UserID            uuid.UUID  `json:"userId" db:"user_id"`
	FullName          string     `json:"fullName" db:"full_name"`
	UserName          string     `json:"userName" db:"user_name"`
	UserEmail         string     `json:"userEmail" db:"user_email"`
	Position          string     `json:"position" db:"position"`
	PhoneNumber       string     `json:"phoneNumber" db:"phone_number"`
	Address           string     `json:"address" db:"address"`
	SalaryBase        float64    `json:"salaryBase" db:"salary_base"`
	JoinDate          time.Time  `json:"joinDate" db:"join_date"`
	BankName          string     `json:"bankName" db:"bank_name"`
	BankAccountNumber string     `json:"bankAccountNumber" db:"bank_account_number"`
	BankAccountHolder string     `json:"bankAccountHolder" db:"bank_account_holder"`
	ScheduleID        *uuid.UUID `json:"scheduleId" db:"schedule_id"`
	EmploymentStatus  string     `json:"employmentStatus" db:"employment_status"`
	JobLevel          string     `json:"jobLevel" db:"job_level"`
	Gender            string     `json:"gender" db:"gender"`
	Division          string     `json:"division" db:"division"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" db:"updated_at"`
}
