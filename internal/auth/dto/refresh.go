package dto

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	UserID       string `json:"user_id" validate:"required"` // Required for security
}