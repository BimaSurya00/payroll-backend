package dto

import "github.com/google/uuid"

type UpdatePayrollStatusRequest struct {
	Status  string    `json:"status" validate:"required,oneof=DRAFT APPROVED PAID CANCELLED"`
	Notes   string    `json:"notes" validate:"omitempty,max=500"`
	UserID  uuid.UUID `json:"-"` // Populated by handler from JWT claims
}
