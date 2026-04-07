package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hris/internal/schedule/dto"
	"hris/internal/schedule/entity"
	"hris/internal/schedule/helper"
	"hris/internal/schedule/repository"
)

type scheduleService struct {
	repo repository.ScheduleRepository
}

func NewScheduleService(repo repository.ScheduleRepository) ScheduleService {
	return &scheduleService{repo: repo}
}

func (s *scheduleService) CreateSchedule(ctx context.Context, req *dto.CreateScheduleRequest, companyID string) (*dto.ScheduleResponse, error) {
	scheduleID := uuid.New().String()

	schedule := &entity.Schedule{
		ID:                  scheduleID,
		CompanyID:           companyID,
		Name:                req.Name,
		TimeIn:              req.TimeIn,
		TimeOut:             req.TimeOut,
		AllowedLateMinutes:  req.AllowedLateMinutes,
		OfficeLat:           req.OfficeLat,
		OfficeLong:          req.OfficeLong,
		AllowedRadiusMeters: req.AllowedRadiusMeters,
		Description:         req.Description,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return helper.ToScheduleResponse(schedule), nil
}

func (s *scheduleService) GetScheduleByID(ctx context.Context, id string, companyID string) (*dto.ScheduleResponse, error) {
	scheduleID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule ID format: %w", err)
	}

	schedule, err := s.repo.FindByIDAndCompany(ctx, scheduleID, companyID)
	if err != nil {
		if errors.Is(err, repository.ErrScheduleNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return helper.ToScheduleResponse(schedule), nil
}

func (s *scheduleService) GetAllSchedules(ctx context.Context, page, perPage int, path string, companyID string) (*helper.Pagination[*dto.ScheduleResponse], error) {
	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	schedules, err := s.repo.FindAllByCompany(ctx, companyID, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}

	total, err := s.repo.CountByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to count schedules: %w", err)
	}

	scheduleResponses := helper.ToScheduleResponses(schedules)
	pagination := helper.NewPagination(scheduleResponses, page, perPage, total, path)

	return pagination, nil
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id string, companyID string, req *dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	if !req.HasUpdates() {
		return nil, fmt.Errorf("no fields to update")
	}

	scheduleID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule ID format: %w", err)
	}

	// Check if schedule exists and belongs to company
	_, err = s.repo.FindByIDAndCompany(ctx, scheduleID, companyID)
	if err != nil {
		if errors.Is(err, repository.ErrScheduleNotFound) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.TimeIn != nil {
		updates["time_in"] = *req.TimeIn
	}
	if req.TimeOut != nil {
		updates["time_out"] = *req.TimeOut
	}
	if req.AllowedLateMinutes != nil {
		updates["allowed_late_minutes"] = *req.AllowedLateMinutes
	}
	if req.OfficeLat != nil {
		updates["office_lat"] = *req.OfficeLat
	}
	if req.OfficeLong != nil {
		updates["office_long"] = *req.OfficeLong
	}
	if req.AllowedRadiusMeters != nil {
		updates["allowed_radius_meters"] = *req.AllowedRadiusMeters
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if err := s.repo.Update(ctx, scheduleID, updates); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	// Fetch updated schedule
	updatedSchedule, err := s.repo.FindByID(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated schedule: %w", err)
	}

	return helper.ToScheduleResponse(updatedSchedule), nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, id string, companyID string) error {
	scheduleID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid schedule ID format: %w", err)
	}

	// Check if schedule exists and belongs to company
	_, err = s.repo.FindByIDAndCompany(ctx, scheduleID, companyID)
	if err != nil {
		if errors.Is(err, repository.ErrScheduleNotFound) {
			return ErrScheduleNotFound
		}
		return fmt.Errorf("failed to get schedule: %w", err)
	}

	// Check if schedule is being used by any employee
	count, err := s.repo.FindByEmployeeCount(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("failed to check schedule usage: %w", err)
	}

	if count > 0 {
		return ErrScheduleInUse
	}

	if err := s.repo.DeleteByCompany(ctx, scheduleID, companyID); err != nil {
		if errors.Is(err, repository.ErrScheduleNotFound) {
			return ErrScheduleNotFound
		}
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	return nil
}
