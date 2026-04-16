package entity

import "time"

// Company represents a tenant company in the SaaS system
type Company struct {
	ID                  string    `json:"id" db:"id"`
	Name                string    `json:"name" db:"name"`
	Slug                string    `json:"slug" db:"slug"`
	IsActive            bool      `json:"isActive" db:"is_active"`
	Plan                string    `json:"plan" db:"plan"`
	MaxEmployees        int       `json:"maxEmployees" db:"max_employees"`
	OfficeLat           *float64  `json:"officeLat" db:"office_lat"`
	OfficeLong          *float64  `json:"officeLong" db:"office_long"`
	AllowedRadiusMeters *int      `json:"allowedRadiusMeters" db:"allowed_radius_meters"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}
