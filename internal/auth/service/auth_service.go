package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hris/config"
	"hris/internal/auth/dto"
	"hris/internal/auth/helper"
	"hris/internal/auth/repository"
	userDto "hris/internal/user/dto"
	"hris/internal/user/entity"
	userRepo "hris/internal/user/repository"
	"hris/shared/constants"
	sharedHelper "hris/shared/helper"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrAccountDeactivated  = errors.New("account is deactivated")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidOldPassword  = errors.New("old password is incorrect")
	ErrSamePassword        = errors.New("new password must be different from old password")
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID, tokenID string) error
	LogoutAll(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error
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
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := sharedHelper.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with company_id from request
	user := entity.NewUser(req.Name, req.Email, hashedPassword, constants.RoleUser, req.CompanyID)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token pair with company_id
	tokenPair, refreshTokenID, err := s.jwtHelper.GenerateTokenPair(user.ID, user.Role, user.CompanyID)
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
		return nil, ErrInvalidCredentials
	}

	// Check password
	if !sharedHelper.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrAccountDeactivated
	}

	// Generate token pair with user's company_id
	tokenPair, refreshTokenID, err := s.jwtHelper.GenerateTokenPair(user.ID, user.Role, user.CompanyID)
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
		return nil, ErrInvalidRefreshToken
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

	// 2. Generate new token pair with user's company_id
	newTokenPair, newRefreshTokenID, err := s.jwtHelper.GenerateTokenPair(user.ID, user.Role, user.CompanyID)
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

func (s *authService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	// 1. Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 2. Verify old password
	if !sharedHelper.CheckPassword(user.Password, req.OldPassword) {
		return ErrInvalidOldPassword
	}

	// 3. Check new password is different
	if sharedHelper.CheckPassword(user.Password, req.NewPassword) {
		return ErrSamePassword
	}

	// 4. Hash new password
	hashedPassword, err := sharedHelper.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Update password using entity Update method
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// 6. Logout semua device (invalidate semua refresh token) — security best practice
	if err := s.tokenRepo.DeleteAllUserRefreshTokens(ctx, userID); err != nil {
		// Log tapi jangan gagalkan — password sudah berubah
		fmt.Printf("Warning: failed to revoke tokens after password change: %v\n", err)
	}

	return nil
}
