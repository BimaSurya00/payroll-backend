package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/entity"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/helper"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/repository"
	sharedHelper "github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type service struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := sharedHelper.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := entity.NewUser(req.Name, req.Email, hashedPassword, req.Role)

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *service) GetUsers(ctx context.Context, page, perPage int, path string) (*helper.Pagination[*dto.UserResponse], error) {
	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	users, err := s.repo.FindAll(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	userResponses := dto.ToUserResponses(users)
	pagination := helper.NewPagination(userResponses, page, perPage, total, path)

	return pagination, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *service) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if !req.HasUpdates() {
		return nil, fmt.Errorf("no fields to update")
	}

	// Check if user exists
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Build update map
	updates := bson.M{
		"updatedAt": time.Now(),
	}

	if req.Name != nil {
		updates["name"] = *req.Name
	}

	if req.Email != nil {
		existingUser, err := s.repo.FindByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailAlreadyExists
		}
		updates["email"] = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := sharedHelper.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updates["password"] = hashedPassword
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Update user
	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	// Fetch updated user
	updatedUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dto.ToUserResponse(updatedUser), nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
