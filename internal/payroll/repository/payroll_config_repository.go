package repository

import (
	"context"

	"hris/internal/payroll/entity"
)

type PayrollConfigRepository interface {
	Create(ctx context.Context, config *entity.PayrollConfig) error
	FindByID(ctx context.Context, id string) (*entity.PayrollConfig, error)
	FindByCode(ctx context.Context, code string) (*entity.PayrollConfig, error)
	FindAll(ctx context.Context) ([]*entity.PayrollConfig, error)
	FindActiveByType(ctx context.Context, configType string) ([]*entity.PayrollConfig, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
