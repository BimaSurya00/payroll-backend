package repository

import (
	"context"

	"hris/internal/user/entity"
)

// UserRepository defines the interface for user data access (PostgreSQL)
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, skip, limit int64) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	FindByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
	// Company-scoped queries for multi-tenancy
	FindAllByCompany(ctx context.Context, companyID string, skip, limit int64) ([]*entity.User, error)
	CountByCompany(ctx context.Context, companyID string) (int64, error)
}
