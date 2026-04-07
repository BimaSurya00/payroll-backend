package dto

import (
	"time"

	"hris/internal/user/entity"
)

type UserResponse struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"companyId"`
	CompanyName     string    `json:"companyName"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	IsActive        bool      `json:"isActive"`
	ProfileImageUrl string    `json:"profileImageUrl"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func ToUserResponse(user *entity.User, companyName ...string) *UserResponse {
	resp := &UserResponse{
		ID:              user.ID,
		CompanyID:       user.CompanyID,
		Name:            user.Name,
		Email:           user.Email,
		Role:            user.Role,
		IsActive:        user.IsActive,
		ProfileImageUrl: user.ProfileImageUrl,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}

	if len(companyName) > 0 {
		resp.CompanyName = companyName[0]
	}

	return resp
}

func ToUserResponses(users []*entity.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
