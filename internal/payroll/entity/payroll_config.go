package entity

import "time"

type PayrollConfig struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Code            string    `json:"code" db:"code"`
	Type            string    `json:"type" db:"type"`            // EARNING | DEDUCTION
	Amount          float64   `json:"amount" db:"amount"`
	CalculationType string    `json:"calculationType" db:"calculation_type"` // FIXED | PER_DAY | PERCENTAGE
	IsActive        bool      `json:"isActive" db:"is_active"`
	Description     string    `json:"description" db:"description"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}
