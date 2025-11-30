package dto

import (
	"time"

	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/entity"
)

type UserResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	IsActive        bool      `json:"isActive"`
	ProfileImageUrl string    `json:"profileImageUrl"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func ToUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		Role:            user.Role,
		IsActive:        user.IsActive,
		ProfileImageUrl: user.ProfileImageUrl,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func ToUserResponses(users []*entity.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}