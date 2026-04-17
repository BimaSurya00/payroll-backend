package entity

import "time"

type Holiday struct {
	ID          string    `json:"id" db:"id"`
	CompanyID   string    `json:"companyId" db:"company_id"`
	Name        string    `json:"name" db:"name"`
	Date        time.Time `json:"date" db:"date"`
	Type        string    `json:"type" db:"type"` // NATIONAL, COMPANY, OPTIONAL
	IsRecurring bool      `json:"isRecurring" db:"is_recurring"`
	Year        *int      `json:"year,omitempty" db:"year"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
