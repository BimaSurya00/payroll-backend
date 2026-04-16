package dto

type CreateCompanyRequest struct {
	Name                string   `json:"name" validate:"required,min=2,max=200"`
	Slug                string   `json:"slug" validate:"required,min=2,max=100,slug"`
	Plan                string   `json:"plan" validate:"omitempty,oneof=free starter pro enterprise"`
	MaxEmployees        int      `json:"maxEmployees" validate:"omitempty,min=1"`
	OfficeLat           *float64 `json:"officeLat" validate:"omitempty,gte=-90,lte=90"`
	OfficeLong          *float64 `json:"officeLong" validate:"omitempty,gte=-180,lte=180"`
	AllowedRadiusMeters *int     `json:"allowedRadiusMeters" validate:"omitempty,min=10,max=1000"`
}

type UpdateCompanyRequest struct {
	Name                *string  `json:"name" validate:"omitempty,min=2,max=200"`
	Slug                *string  `json:"slug" validate:"omitempty,min=2,max=100,slug"`
	IsActive            *bool    `json:"isActive"`
	Plan                *string  `json:"plan" validate:"omitempty,oneof=free starter pro enterprise"`
	MaxEmployees        *int     `json:"maxEmployees" validate:"omitempty,min=1"`
	OfficeLat           *float64 `json:"officeLat" validate:"omitempty,gte=-90,lte=90"`
	OfficeLong          *float64 `json:"officeLong" validate:"omitempty,gte=-180,lte=180"`
	AllowedRadiusMeters *int     `json:"allowedRadiusMeters" validate:"omitempty,min=10,max=1000"`
}

type CompanyResponse struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Slug                string   `json:"slug"`
	IsActive            bool     `json:"isActive"`
	Plan                string   `json:"plan"`
	MaxEmployees        int      `json:"maxEmployees"`
	OfficeLat           *float64 `json:"officeLat"`
	OfficeLong          *float64 `json:"officeLong"`
	AllowedRadiusMeters *int     `json:"allowedRadiusMeters"`
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
}
