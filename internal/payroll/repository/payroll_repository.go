package repository

import (
	"context"

	"github.com/google/uuid"
	"hris/internal/payroll/entity"
)

type PayrollRepository interface {
	CreateWithItems(ctx context.Context, payroll *entity.Payroll, items []*entity.PayrollItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Payroll, error)
	FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.Payroll, error)
	FindByIDWithItems(ctx context.Context, id uuid.UUID) (*entity.PayrollWithItems, error)
	FindByIDWithItemsAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.PayrollWithItems, error)
	FindAll(ctx context.Context, skip, limit int64) ([]*entity.Payroll, error)
	FindAllByCompany(ctx context.Context, companyID string, skip, limit int64) ([]*entity.Payroll, error)
	FindByPeriod(ctx context.Context, periodStart, periodEnd string) ([]*entity.Payroll, error)
	FindByPeriodAndCompany(ctx context.Context, companyID string, periodStart, periodEnd string) ([]*entity.Payroll, error)
	FindByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, periodStart, periodEnd string) (*entity.Payroll, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Count(ctx context.Context) (int64, error)
	CountByCompany(ctx context.Context, companyID string) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByCompany(ctx context.Context, id uuid.UUID, companyID string) error
	FindByEmployeeID(ctx context.Context, employeeID string, skip, limit int64) ([]*entity.Payroll, error)
	CountByEmployeeID(ctx context.Context, employeeID string) (int64, error)
}
