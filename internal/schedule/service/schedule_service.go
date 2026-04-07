package service

import (
	"context"
	"errors"

	"hris/internal/schedule/dto"
	"hris/internal/schedule/helper"
)

var (
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrScheduleInUse    = errors.New("schedule is being used by employees")
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, req *dto.CreateScheduleRequest, companyID string) (*dto.ScheduleResponse, error)
	GetScheduleByID(ctx context.Context, id string, companyID string) (*dto.ScheduleResponse, error)
	GetAllSchedules(ctx context.Context, page, perPage int, path string, companyID string) (*helper.Pagination[*dto.ScheduleResponse], error)
	UpdateSchedule(ctx context.Context, id string, companyID string, req *dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error)
	DeleteSchedule(ctx context.Context, id string, companyID string) error
}
