package helper

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/entity"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
)

type JWTHelper struct {
	config *config.JWTConfig
}

func NewJWTHelper(cfg *config.JWTConfig) *JWTHelper {
	return &JWTHelper{config: cfg}
}

func (h *JWTHelper) GenerateTokenPair(userID, role string) (*entity.TokenPair, string, error) {
	// Generate access token (JWT)
	accessToken, accessExpiry, err := h.generateAccessToken(userID, role)
	if err != nil {
		return nil, "", err
	}

	// Generate refresh token (secure random UUID v7)
	refreshTokenID := uuid.New().String()

	return &entity.TokenPair{
		AccessToken:       accessToken,
		RefreshToken:      refreshTokenID,
		AccessTokenExpiry: accessExpiry.Unix(),
	}, refreshTokenID, nil
}

func (h *JWTHelper) generateAccessToken(userID, role string) (string, time.Time, error) {
	expiryTime := time.Now().Add(h.config.AccessExpiry)

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    constants.TokenTypeAccess,
		"exp":     expiryTime.Unix(),
		"iat":     time.Now().Unix(),
		"jti":     uuid.New().String(), // JWT ID for tracking
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiryTime, nil
}

func (h *JWTHelper) ValidateAccessToken(tokenString string) (*entity.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != constants.TokenTypeAccess {
		return nil, fmt.Errorf("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing user_id in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("missing role in token")
	}

	return &entity.TokenClaims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
	}, nil
}

// GenerateRefreshTokenID generates a new UUID for refresh token
func (h *JWTHelper) GenerateRefreshTokenID() string {
	return uuid.New().String()
}
