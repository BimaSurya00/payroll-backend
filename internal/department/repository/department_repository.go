package repository

import (
	"context"

	"hris/internal/department/entity"
)

type DepartmentRepository interface {
	Create(ctx context.Context, department *entity.Department) error
	FindByID(ctx context.Context, id string) (*entity.Department, error)
	FindByIDAndCompany(ctx context.Context, id string, companyID string) (*entity.Department, error)
	FindByCode(ctx context.Context, code string) (*entity.Department, error)
	FindByCodeAndCompany(ctx context.Context, code string, companyID string) (*entity.Department, error)
	FindAll(ctx context.Context) ([]entity.Department, error)
	FindAllByCompany(ctx context.Context, companyID string) ([]entity.Department, error)
	FindActive(ctx context.Context) ([]entity.Department, error)
	FindActiveByCompany(ctx context.Context, companyID string) ([]entity.Department, error)
	Update(ctx context.Context, department *entity.Department) error
	Delete(ctx context.Context, id string) error
	DeleteByCompany(ctx context.Context, id string, companyID string) error
}
