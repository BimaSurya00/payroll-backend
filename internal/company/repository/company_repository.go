package repository

import (
	"context"

	"hris/internal/company/dto"
	"hris/internal/company/entity"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *entity.Company) error
	FindByID(ctx context.Context, id string) (*entity.Company, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Company, error)
	FindAll(ctx context.Context, page, perPage int) ([]*entity.Company, int64, error)
	FindAllWithStats(ctx context.Context, page, perPage int) ([]*dto.CompanyListItem, int64, error)
	GetStats(ctx context.Context, companyID string) (*dto.CompanyStatsResponse, error)
	Update(ctx context.Context, company *entity.Company) error
	Delete(ctx context.Context, id string) error
}
