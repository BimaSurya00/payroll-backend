package entity

import "time"

type Department struct {
	ID             string    `json:"id" db:"id"`
	CompanyID      string    `json:"companyId" db:"company_id"`
	Name           string    `json:"name" db:"name"`
	Code           string    `json:"code" db:"code"`
	Description    *string   `json:"description,omitempty" db:"description"`
	HeadEmployeeID *string   `json:"headEmployeeId,omitempty" db:"head_employee_id"`
	IsActive       bool      `json:"isActive" db:"is_active"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}
