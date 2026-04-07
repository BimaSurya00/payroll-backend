package dto

type UpdateUserRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=3,max=100,trimmed_string"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Password  *string `json:"password,omitempty" validate:"omitempty,min=8,password_strength"`
	Role      *string `json:"role,omitempty" validate:"omitempty,oneof=SUPER_USER ADMIN USER"`
	IsActive  *bool   `json:"isActive,omitempty"`
	CompanyID *string `json:"companyId,omitempty" validate:"omitempty,uuid"`
}

func (u *UpdateUserRequest) HasUpdates() bool {
	return u.Name != nil || u.Email != nil || u.Password != nil || u.Role != nil || u.IsActive != nil || u.CompanyID != nil
}
