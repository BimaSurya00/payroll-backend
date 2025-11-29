package service

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/itsahyarr/go-fiber-boilerplate/config"
// 	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/dto"
// 	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/entity"
// 	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
// 	"github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
// )

// // Mock repositories
// type mockUserRepository struct {
// 	users map[string]*entity.User
// }

// func newMockUserRepository() *mockUserRepository {
// 	return &mockUserRepository{
// 		users: make(map[string]*entity.User),
// 	}
// }

// func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
// 	m.users[user.Email] = user
// 	return nil
// }

// func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
// 	if user, exists := m.users[email]; exists {
// 		return user, nil
// 	}
// 	return nil, nil
// }

// func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
// 	for _, user := range m.users {
// 		if user.ID == id {
// 			return user, nil
// 		}
// 	}
// 	return nil, nil
// }

// func (m *mockUserRepository) FindAll(ctx context.Context, skip, limit int64) ([]*entity.User, error) {
// 	return nil, nil
// }

// func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
// 	return 0, nil
// }

// func (m *mockUserRepository) Update(ctx context.Context, id string, updates any) error {
// 	return nil
// }

// func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
// 	return nil
// }

// type mockTokenRepository struct {
// 	tokens map[string]string
// }

// func newMockTokenRepository() *mockTokenRepository {
// 	return &mockTokenRepository{
// 		tokens: make(map[string]string),
// 	}
// }

// func (m *mockTokenRepository) StoreRefreshToken(ctx context.Context, userID, token string, expiry time.Duration) error {
// 	m.tokens[userID] = token
// 	return nil
// }

// func (m *mockTokenRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
// 	if token, exists := m.tokens[userID]; exists {
// 		return token, nil
// 	}
// 	return "", nil
// }

// func (m *mockTokenRepository) DeleteRefreshToken(ctx context.Context, userID string) error {
// 	delete(m.tokens, userID)
// 	return nil
// }

// func (m *mockTokenRepository) RefreshTokenExists(ctx context.Context, userID string) (bool, error) {
// 	_, exists := m.tokens[userID]
// 	return exists, nil
// }

// func TestRegister(t *testing.T) {
// 	userRepo := newMockUserRepository()
// 	tokenRepo := newMockTokenRepository()
// 	jwtConfig := &config.JWTConfig{
// 		Secret:        "test-secret",
// 		AccessExpiry:  15 * time.Minute,
// 		RefreshExpiry: 7 * 24 * time.Hour,
// 	}

// 	service := NewAuthService(userRepo, tokenRepo, jwtConfig)

// 	tests := []struct {
// 		name    string
// 		req     *dto.RegisterRequest
// 		wantErr bool
// 		errMsg  string
// 	}{
// 		{
// 			name: "Successful registration",
// 			req: &dto.RegisterRequest{
// 				Name:     "John Doe",
// 				Email:    "john@example.com",
// 				Password: "Password123!",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Duplicate email",
// 			req: &dto.RegisterRequest{
// 				Name:     "Jane Doe",
// 				Email:    "john@example.com",
// 				Password: "Password123!",
// 			},
// 			wantErr: true,
// 			errMsg:  "email already exists",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := service.Register(context.Background(), tt.req)

// 			if tt.wantErr {
// 				if err == nil {
// 					t.Errorf("Expected error but got none")
// 				}
// 				if err.Error() != tt.errMsg {
// 					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
// 				}
// 			} else {
// 				if err != nil {
// 					t.Errorf("Unexpected error: %v", err)
// 				}
// 				if result == nil {
// 					t.Error("Expected result but got nil")
// 				}
// 				if result != nil && result.AccessToken == "" {
// 					t.Error("Expected access token")
// 				}
// 				if result != nil && result.RefreshToken == "" {
// 					t.Error("Expected refresh token")
// 				}
// 			}
// 		})
// 	}
// }

// func TestLogin(t *testing.T) {
// 	userRepo := newMockUserRepository()
// 	tokenRepo := newMockTokenRepository()
// 	jwtConfig := &config.JWTConfig{
// 		Secret:        "test-secret",
// 		AccessExpiry:  15 * time.Minute,
// 		RefreshExpiry: 7 * 24 * time.Hour,
// 	}

// 	service := NewAuthService(userRepo, tokenRepo, jwtConfig)

// 	// Create a test user
// 	hashedPassword, _ := helper.HashPassword("Password123!")
// 	user := entity.NewUser("John Doe", "john@example.com", hashedPassword, constants.RoleUser)
// 	userRepo.users[user.Email] = user

// 	tests := []struct {
// 		name    string
// 		req     *dto.LoginRequest
// 		wantErr bool
// 		errMsg  string
// 	}{
// 		{
// 			name: "Successful login",
// 			req: &dto.LoginRequest{
// 				Email:    "john@example.com",
// 				Password: "Password123!",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Invalid email",
// 			req: &dto.LoginRequest{
// 				Email:    "nonexistent@example.com",
// 				Password: "Password123!",
// 			},
// 			wantErr: true,
// 			errMsg:  "invalid credentials",
// 		},
// 		{
// 			name: "Invalid password",
// 			req: &dto.LoginRequest{
// 				Email:    "john@example.com",
// 				Password: "WrongPassword!",
// 			},
// 			wantErr: true,
// 			errMsg:  "invalid credentials",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := service.Login(context.Background(), tt.req)

// 			if tt.wantErr {
// 				if err == nil {
// 					t.Errorf("Expected error but got none")
// 				}
// 				if err.Error() != tt.errMsg {
// 					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
// 				}
// 			} else {
// 				if err != nil {
// 					t.Errorf("Unexpected error: %v", err)
// 				}
// 				if result == nil {
// 					t.Error("Expected result but got nil")
// 				}
// 				if result != nil && result.AccessToken == "" {
// 					t.Error("Expected access token")
// 				}
// 			}
// 		})
// 	}
// }
