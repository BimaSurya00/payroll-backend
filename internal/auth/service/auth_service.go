package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/dto"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/helper"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/repository"
	userDto "github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/entity"
	userRepo "github.com/itsahyarr/go-fiber-boilerplate/internal/user/repository"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
	sharedHelper "github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID, tokenID string) error
	LogoutAll(ctx context.Context, userID string) error
}

type authService struct {
	userRepo  userRepo.UserRepository
	tokenRepo repository.TokenRepository
	jwtHelper *helper.JWTHelper
	jwtConfig *config.JWTConfig
}

func NewAuthService(
	userRepo userRepo.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtConfig *config.JWTConfig,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtHelper: helper.NewJWTHelper(jwtConfig),
		jwtConfig: jwtConfig,
	}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := sharedHelper.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := entity.NewUser(req.Name, req.Email, hashedPassword, constants.RoleUser)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token pair
	tokenPair, refreshTokenID, err := s.jwtHelper.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in KeyDB with expiration
	expiresAt := time.Now().Add(s.jwtConfig.RefreshExpiry)
	if err := s.tokenRepo.SetRefreshToken(ctx, user.ID, refreshTokenID, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.AuthResponse{
		User:         userDto.ToUserResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.AccessTokenExpiry,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !sharedHelper.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate token pair
	tokenPair, refreshTokenID, err := s.jwtHelper.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token in KeyDB
	expiresAt := time.Now().Add(s.jwtConfig.RefreshExpiry)
	if err := s.tokenRepo.SetRefreshToken(ctx, user.ID, refreshTokenID, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &dto.AuthResponse{
		User:         userDto.ToUserResponse(user),
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.AccessTokenExpiry,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate refresh token format (should be UUID)
	if req.RefreshToken == "" {
		return nil, errors.New("invalid refresh token")
	}

	// Extract user ID from token ID (we need to get it from KeyDB)
	// Since we don't know the user ID, we need to pass it in the request
	// Alternative: Store user_id in the refresh token itself or use a different approach

	// For security, we require both user_id and token_id
	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	// Verify refresh token exists and is valid
	tokenData, err := s.tokenRepo.GetRefreshToken(ctx, req.UserID, req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Verify token hasn't expired (double-check)
	if time.Now().After(tokenData.ExpiresAt) {
		_ = s.tokenRepo.DeleteRefreshToken(ctx, req.UserID, req.RefreshToken)
		return nil, errors.New("refresh token has expired")
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, tokenData.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// ===== REFRESH TOKEN ROTATION (Best Security Practice) =====
	// 1. Delete the old refresh token
	if err := s.tokenRepo.DeleteRefreshToken(ctx, user.ID, req.RefreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// 2. Generate new token pair
	newTokenPair, newRefreshTokenID, err := s.jwtHelper.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// 3. Store new refresh token
	newExpiresAt := time.Now().Add(s.jwtConfig.RefreshExpiry)
	if err := s.tokenRepo.SetRefreshToken(ctx, user.ID, newRefreshTokenID, newExpiresAt); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &dto.AuthResponse{
		User:         userDto.ToUserResponse(user),
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		ExpiresAt:    newTokenPair.AccessTokenExpiry,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID, tokenID string) error {
	// Delete specific refresh token
	return s.tokenRepo.DeleteRefreshToken(ctx, userID, tokenID)
}

func (s *authService) LogoutAll(ctx context.Context, userID string) error {
	// Delete all refresh tokens for the user (logout from all devices)
	return s.tokenRepo.DeleteAllUserRefreshTokens(ctx, userID)
}
