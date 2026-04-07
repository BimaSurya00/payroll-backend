package dto

// CreateCompanyRequest represents the request body for creating a company
type CreateCompanyRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=200"`
	Slug         string `json:"slug" validate:"required,min=2,max=100,slug"`
	Plan         string `json:"plan" validate:"omitempty,oneof=free starter pro enterprise"`
	MaxEmployees int    `json:"maxEmployees" validate:"omitempty,min=1"`
}

// UpdateCompanyRequest represents the request body for updating a company
type UpdateCompanyRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=2,max=200"`
	Slug         *string `json:"slug" validate:"omitempty,min=2,max=100,slug"`
	IsActive     *bool   `json:"isActive"`
	Plan         *string `json:"plan" validate:"omitempty,oneof=free starter pro enterprise"`
	MaxEmployees *int    `json:"maxEmployees" validate:"omitempty,min=1"`
}

// CompanyResponse represents the API response for a company
type CompanyResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	IsActive     bool   `json:"isActive"`
	Plan         string `json:"plan"`
	MaxEmployees int    `json:"maxEmployees"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
