package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account stored in PostgreSQL
type User struct {
	ID              string    `json:"id" db:"id"`
	CompanyID       string    `json:"companyId" db:"company_id"`
	Name            string    `json:"name" db:"name"`
	Email           string    `json:"email" db:"email"`
	Password        string    `json:"-" db:"password"`
	Role            string    `json:"role" db:"role"`
	IsActive        bool      `json:"isActive" db:"is_active"`
	ProfileImageUrl string    `json:"profileImageUrl" db:"profile_image_url"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

func NewUser(name, email, password, role, companyID string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New().String(),
		CompanyID: companyID,
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
