package service

import (
	"context"

	"hris/internal/company/dto"
	"hris/internal/company/entity"
)

// CompanyService defines the interface for company business logic
type CompanyService interface {
	Create(ctx context.Context, req *dto.CreateCompanyRequest) (*entity.Company, error)
	GetByID(ctx context.Context, id string) (*entity.Company, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Company, error)
	GetAll(ctx context.Context, page, perPage int) ([]*entity.Company, int64, error)
	Update(ctx context.Context, id string, req *dto.UpdateCompanyRequest) (*entity.Company, error)
	Delete(ctx context.Context, id string) error
}
