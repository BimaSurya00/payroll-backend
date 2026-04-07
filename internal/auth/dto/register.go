package dto

type RegisterRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=100,trimmed_string"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,password_strength"`
	CompanyID string `json:"companyId" validate:"required,uuid"`
}
