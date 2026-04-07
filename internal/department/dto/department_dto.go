package dto

type CreateDepartmentRequest struct {
	Name           string  `json:"name" validate:"required,min=2,max=100"`
	Code           string  `json:"code" validate:"required,min=2,max=20,uppercase"`
	Description    *string `json:"description,omitempty" validate:"omitempty,max=500"`
	HeadEmployeeID *string `json:"headEmployeeId,omitempty" validate:"omitempty,uuid"`
	IsActive       *bool   `json:"isActive,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name           *string `json:"name" validate:"omitempty,min=2,max=100"`
	Code           *string `json:"code" validate:"omitempty,min=2,max=20,uppercase"`
	Description    *string `json:"description,omitempty" validate:"omitempty,max=500"`
	HeadEmployeeID *string `json:"headEmployeeId,omitempty" validate:"omitempty,uuid"`
	IsActive       *bool   `json:"isActive,omitempty"`
}

type DepartmentResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Code             string  `json:"code"`
	Description      *string `json:"description,omitempty"`
	HeadEmployeeID   *string `json:"headEmployeeId,omitempty"`
	HeadEmployeeName *string `json:"headEmployeeName,omitempty"`
	IsActive         bool    `json:"isActive"`
	EmployeeCount    int     `json:"employeeCount"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}
