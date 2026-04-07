package entity

import "time"

type Payroll struct {
	ID             string     `json:"id" db:"id"`
	CompanyID      string     `json:"companyId" db:"company_id"`
	EmployeeID     string     `json:"employeeId" db:"employee_id"`
	PeriodStart    time.Time  `json:"periodStart" db:"period_start"`
	PeriodEnd      time.Time  `json:"periodEnd" db:"period_end"`
	BaseSalary     float64    `json:"baseSalary" db:"base_salary"`
	TotalAllowance float64    `json:"totalAllowance" db:"total_allowance"`
	TotalDeduction float64    `json:"totalDeduction" db:"total_deduction"`
	NetSalary      float64    `json:"netSalary" db:"net_salary"`
	Status         string     `json:"status" db:"status"` // DRAFT, APPROVED, PAID
	ApprovedBy     *string    `json:"approvedBy,omitempty" db:"approved_by"`
	ApprovedAt     *time.Time `json:"approvedAt,omitempty" db:"approved_at"`
	PaidAt         *time.Time `json:"paidAt,omitempty" db:"paid_at"`
	CancelledAt    *time.Time `json:"cancelledAt,omitempty" db:"cancelled_at"`
	Notes          string     `json:"notes" db:"notes"`
	GeneratedAt    time.Time  `json:"generatedAt" db:"generated_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
}

type PayrollItem struct {
	ID        string    `json:"id" db:"id"`
	PayrollID string    `json:"payrollId" db:"payroll_id"`
	Name      string    `json:"name" db:"name"`
	Amount    float64   `json:"amount" db:"amount"`
	Type      string    `json:"type" db:"type"` // EARNING, DEDUCTION
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// PayrollWithItems combines payroll header with items
type PayrollWithItems struct {
	Payroll *Payroll       `json:"payroll"`
	Items   []*PayrollItem `json:"items"`
}
