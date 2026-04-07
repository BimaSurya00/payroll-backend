package repository

import (
	"context"

	"github.com/google/uuid"
	"hris/internal/schedule/entity"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error)
	FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.Schedule, error)
	FindAll(ctx context.Context, skip, limit int64) ([]*entity.Schedule, error)
	FindAllByCompany(ctx context.Context, companyID string, skip, limit int64) ([]*entity.Schedule, error)
	Count(ctx context.Context) (int64, error)
	CountByCompany(ctx context.Context, companyID string) (int64, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByCompany(ctx context.Context, id uuid.UUID, companyID string) error
	FindByEmployeeCount(ctx context.Context, employeeID uuid.UUID) (int64, error)
}
