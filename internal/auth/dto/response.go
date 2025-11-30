package dto

import "github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"

type AuthResponse struct {
	User         *dto.UserResponse `json:"user"`
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	ExpiresAt    int64             `json:"expiresAt"`
	TokenType    string            `json:"tokenType"`
}

func NewAuthResponse(user *dto.UserResponse, accessToken, refreshToken string, expiresAt int64) *AuthResponse {
	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}
}