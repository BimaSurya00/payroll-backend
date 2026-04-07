package repository

import (
	"context"
	"time"

	"hris/internal/holiday/entity"
)

type HolidayRepository interface {
	Create(ctx context.Context, holiday *entity.Holiday) error
	FindByID(ctx context.Context, id string) (*entity.Holiday, error)
	FindAll(ctx context.Context, year int) ([]*entity.Holiday, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error)
	Update(ctx context.Context, holiday *entity.Holiday) error
	Delete(ctx context.Context, id string) error
	IsHoliday(ctx context.Context, date time.Time) (bool, error)
}
