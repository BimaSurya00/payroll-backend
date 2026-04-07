package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"hris/internal/company/dto"
	"hris/internal/company/entity"
	"hris/internal/company/repository"
)

var (
	ErrCompanyNotFound   = errors.New("company not found")
	ErrSlugAlreadyExists = errors.New("company slug already exists")
)

type companyServiceImpl struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository) CompanyService {
	return &companyServiceImpl{companyRepo: companyRepo}
}

func (s *companyServiceImpl) Create(ctx context.Context, req *dto.CreateCompanyRequest) (*entity.Company, error) {
	// Check slug uniqueness
	existing, err := s.companyRepo.FindBySlug(ctx, req.Slug)
	if err == nil && existing != nil {
		return nil, ErrSlugAlreadyExists
	}

	plan := req.Plan
	if plan == "" {
		plan = "free"
	}

	maxEmployees := req.MaxEmployees
	if maxEmployees == 0 {
		maxEmployees = 25
	}

	company := &entity.Company{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Slug:         req.Slug,
		IsActive:     true,
		Plan:         plan,
		MaxEmployees: maxEmployees,
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *companyServiceImpl) GetByID(ctx context.Context, id string) (*entity.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrCompanyNotFound
	}
	return company, nil
}

func (s *companyServiceImpl) GetBySlug(ctx context.Context, slug string) (*entity.Company, error) {
	company, err := s.companyRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, ErrCompanyNotFound
	}
	return company, nil
}

func (s *companyServiceImpl) GetAll(ctx context.Context, page, perPage int) ([]*entity.Company, int64, error) {
	return s.companyRepo.FindAll(ctx, page, perPage)
}

func (s *companyServiceImpl) Update(ctx context.Context, id string, req *dto.UpdateCompanyRequest) (*entity.Company, error) {
	company, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrCompanyNotFound
	}

	if req.Name != nil {
		company.Name = *req.Name
	}
	if req.Slug != nil {
		// Check slug uniqueness if changing
		existing, err := s.companyRepo.FindBySlug(ctx, *req.Slug)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrSlugAlreadyExists
		}
		company.Slug = *req.Slug
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}
	if req.Plan != nil {
		company.Plan = *req.Plan
	}
	if req.MaxEmployees != nil {
		company.MaxEmployees = *req.MaxEmployees
	}

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *companyServiceImpl) Delete(ctx context.Context, id string) error {
	_, err := s.companyRepo.FindByID(ctx, id)
	if err != nil {
		return ErrCompanyNotFound
	}
	return s.companyRepo.Delete(ctx, id)
}
