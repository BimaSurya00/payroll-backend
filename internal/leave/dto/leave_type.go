package dto

type CreateLeaveTypeRequest struct {
	Name             string `json:"name" validate:"required,min=3,max=100"`
	Code             string `json:"code" validate:"required,min=2,max=50,uppercase"`
	Description      string `json:"description" validate:"omitempty,max=500"`
	MaxDaysPerYear   int    `json:"maxDaysPerYear" validate:"omitempty,min=0,max=365"`
	DefaultDays      int    `json:"defaultDays" validate:"omitempty,min=0,max=365"`
	IsPaid           bool   `json:"isPaid"`
	RequiresApproval bool   `json:"requiresApproval"`
	IsActive         bool   `json:"isActive"`
	Color            string `json:"color" validate:"omitempty,max=7"`
}

type UpdateLeaveTypeRequest struct {
	Name             *string `json:"name" validate:"omitempty,min=3,max=100"`
	Description      *string `json:"description" validate:"omitempty,max=500"`
	MaxDaysPerYear   *int    `json:"maxDaysPerYear" validate:"omitempty,min=1,max=365"`
	IsPaid           *bool   `json:"isPaid"`
	RequiresApproval *bool   `json:"requiresApproval"`
	IsActive         *bool   `json:"isActive"`
}

type LeaveTypeResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Code             string `json:"code"`
	Description      string `json:"description"`
	MaxDaysPerYear   int    `json:"maxDaysPerYear"`
	DefaultDays      int    `json:"defaultDays"`
	IsPaid           bool   `json:"isPaid"`
	RequiresApproval bool   `json:"requiresApproval"`
	IsActive         bool   `json:"isActive"`
	Color            string `json:"color"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}
