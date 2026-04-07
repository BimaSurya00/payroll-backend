package entity

import "time"

type TokenPair struct {
	AccessToken       string `json:"accessToken"`
	RefreshToken      string `json:"refreshToken"`
	AccessTokenExpiry int64  `json:"accessTokenExpiry"`
}

type TokenClaims struct {
	UserID    string `json:"userId"`
	Role      string `json:"role"`
	Type      string `json:"type"`
	CompanyID string `json:"companyId"`
}

type RefreshTokenData struct {
	TokenID   string    `json:"tokenId"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}
