package service

import (
	"context"

	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/helper"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUsers(ctx context.Context, page, perPage int, path string) (*helper.Pagination[*dto.UserResponse], error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}
