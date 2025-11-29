package dto

import "github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"

type AuthResponse struct {
	User         *dto.UserResponse `json:"user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresAt    int64             `json:"expires_at"`
	TokenType    string            `json:"token_type"`
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
