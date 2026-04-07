package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hris/database"
	"hris/internal/auth/entity"
	"hris/shared/constants"
)

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID, tokenID string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, userID, tokenID string) (*entity.RefreshTokenData, error)
	DeleteRefreshToken(ctx context.Context, userID, tokenID string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID string) error
	RefreshTokenExists(ctx context.Context, userID, tokenID string) (bool, error)
}

type tokenRepository struct {
	keydb *database.KeyDB
}

func NewTokenRepository(keydb *database.KeyDB) TokenRepository {
	return &tokenRepository{keydb: keydb}
}

func (r *tokenRepository) SetRefreshToken(ctx context.Context, userID, tokenID string, expiresAt time.Time) error {
	key := r.buildRefreshTokenKey(userID, tokenID)

	tokenData := entity.RefreshTokenData{
		TokenID:   tokenID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(tokenData)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("token expiration time is in the past")
	}

	return r.keydb.Set(ctx, key, data, ttl)
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, userID, tokenID string) (*entity.RefreshTokenData, error) {
	key := r.buildRefreshTokenKey(userID, tokenID)

	data, err := r.keydb.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	}

	var tokenData entity.RefreshTokenData
	if err := json.Unmarshal([]byte(data), &tokenData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	// Double-check expiration
	if time.Now().After(tokenData.ExpiresAt) {
		_ = r.DeleteRefreshToken(ctx, userID, tokenID)
		return nil, fmt.Errorf("refresh token has expired")
	}

	return &tokenData, nil
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, userID, tokenID string) error {
	key := r.buildRefreshTokenKey(userID, tokenID)
	return r.keydb.Del(ctx, key)
}

func (r *tokenRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf(constants.CacheKeyRefreshTokenPrefix, userID, "*")

	keys, err := r.keydb.GetKeys(ctx, pattern)
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	if len(keys) == 0 {
		return nil // No tokens to delete
	}

	return r.keydb.Del(ctx, keys...)
}

func (r *tokenRepository) RefreshTokenExists(ctx context.Context, userID, tokenID string) (bool, error) {
	key := r.buildRefreshTokenKey(userID, tokenID)
	return r.keydb.Exists(ctx, key)
}

func (r *tokenRepository) buildRefreshTokenKey(userID, tokenID string) string {
	return fmt.Sprintf(constants.CacheKeyRefreshTokenPrefix, userID, tokenID)
}
