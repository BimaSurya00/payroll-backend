package dto

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
	UserID       string `json:"userId" validate:"required"` // Required for security
}