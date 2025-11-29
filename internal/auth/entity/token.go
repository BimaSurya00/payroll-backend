package entity

import "time"

type TokenPair struct {
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	AccessTokenExpiry int64  `json:"access_token_expiry"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
}

type RefreshTokenData struct {
	TokenID   string    `json:"token_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}