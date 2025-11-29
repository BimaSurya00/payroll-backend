package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100,trimmed_string"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_strength"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}