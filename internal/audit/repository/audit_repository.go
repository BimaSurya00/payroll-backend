package repository

import (
	"context"
	"time"

	"hris/internal/audit/entity"
)

type AuditRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	FindByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*entity.AuditLog, error)
	FindByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.AuditLog, error)
	FindAll(ctx context.Context, filter AuditFilter, limit, offset int) ([]*entity.AuditLog, error)
	Count(ctx context.Context, filter AuditFilter) (int64, error)
}

type AuditFilter struct {
	UserID       *string
	Action       *string
	ResourceType *string
	CompanyID    *string
	DateFrom     *time.Time
	DateTo       *time.Time
}
