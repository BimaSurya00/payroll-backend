package dto

type CreateOvertimePolicyRequest struct {
	Name                     string  `json:"name" validate:"required,min=3,max=100"`
	Description              string  `json:"description" validate:"omitempty,max=500"`
	RateType                 string  `json:"rateType" validate:"required,oneof=FIXED MULTIPLIER"`
	RateMultiplier           float64 `json:"rateMultiplier" validate:"required_if=RateType MULTIPLIER,gt=0"`
	FixedAmount              float64 `json:"fixedAmount" validate:"required_if=RateType FIXED,gt=0"`
	MinOvertimeMinutes       int     `json:"minOvertimeMinutes" validate:"required,min=1"`
	MaxOvertimeHoursPerDay   float64 `json:"maxOvertimeHoursPerDay" validate:"required,gt=0"`
	MaxOvertimeHoursPerMonth float64 `json:"maxOvertimeHoursPerMonth" validate:"required,gt=0"`
	RequiresApproval         bool    `json:"requiresApproval"`
}

type UpdateOvertimePolicyRequest struct {
	Name                     *string  `json:"name" validate:"omitempty,min=3,max=100"`
	Description              *string  `json:"description" validate:"omitempty,max=500"`
	RateType                 *string  `json:"rateType" validate:"omitempty,oneof=FIXED MULTIPLIER"`
	RateMultiplier           *float64 `json:"rateMultiplier" validate:"omitempty,gt=0"`
	FixedAmount              *float64 `json:"fixedAmount" validate:"omitempty,gt=0"`
	MinOvertimeMinutes       *int     `json:"minOvertimeMinutes" validate:"omitempty,min=1"`
	MaxOvertimeHoursPerDay   *float64 `json:"maxOvertimeHoursPerDay" validate:"omitempty,gt=0"`
	MaxOvertimeHoursPerMonth *float64 `json:"maxOvertimeHoursPerMonth" validate:"omitempty,gt=0"`
	RequiresApproval         *bool    `json:"requiresApproval"`
	IsActive                 *bool    `json:"isActive"`
}

type OvertimePolicyResponse struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	Description              string  `json:"description"`
	RateType                 string  `json:"rateType"`
	RateMultiplier           float64 `json:"rateMultiplier"`
	FixedAmount              float64 `json:"fixedAmount"`
	MinOvertimeMinutes       int     `json:"minOvertimeMinutes"`
	MaxOvertimeHoursPerDay   float64 `json:"maxOvertimeHoursPerDay"`
	MaxOvertimeHoursPerMonth float64 `json:"maxOvertimeHoursPerMonth"`
	RequiresApproval         bool    `json:"requiresApproval"`
	IsActive                 bool    `json:"isActive"`
	CreatedAt                string  `json:"createdAt"`
	UpdatedAt                string  `json:"updatedAt"`
}
