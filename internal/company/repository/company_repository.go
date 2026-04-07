package repository

import (
	"context"

	"hris/internal/company/entity"
)

// CompanyRepository defines the interface for company data access
type CompanyRepository interface {
	Create(ctx context.Context, company *entity.Company) error
	FindByID(ctx context.Context, id string) (*entity.Company, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Company, error)
	FindAll(ctx context.Context, page, perPage int) ([]*entity.Company, int64, error)
	Update(ctx context.Context, company *entity.Company) error
	Delete(ctx context.Context, id string) error
}
