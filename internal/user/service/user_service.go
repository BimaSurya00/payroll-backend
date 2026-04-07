package service

import (
	"context"

	"hris/internal/user/dto"
	"hris/internal/user/helper"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest, companyID string) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, page, perPage int, path string, companyID string, userRole string) (*helper.Pagination[*dto.UserResponse], error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}
