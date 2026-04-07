package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"hris/internal/holiday/dto"
	"hris/internal/holiday/entity"
	"hris/internal/holiday/repository"
)

type HolidayService interface {
	Create(ctx context.Context, req *dto.CreateHolidayRequest) (*dto.HolidayResponse, error)
	GetByID(ctx context.Context, id string) (*dto.HolidayResponse, error)
	GetAllByYear(ctx context.Context, year int) ([]dto.HolidayResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateHolidayRequest) (*dto.HolidayResponse, error)
	Delete(ctx context.Context, id string) error
}

type holidayService struct {
	holidayRepo repository.HolidayRepository
}

func NewHolidayService(holidayRepo repository.HolidayRepository) HolidayService {
	return &holidayService{
		holidayRepo: holidayRepo,
	}
}

func (s *holidayService) Create(ctx context.Context, req *dto.CreateHolidayRequest) (*dto.HolidayResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Check if holiday already exists for this date
	if isHoliday, _ := s.holidayRepo.IsHoliday(ctx, date); isHoliday {
		return nil, errors.New("holiday already exists for this date")
	}

	now := time.Now()
	holiday := &entity.Holiday{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Date:        date,
		Type:        req.Type,
		IsRecurring: req.IsRecurring,
		Year:        req.Year,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.holidayRepo.Create(ctx, holiday); err != nil {
		return nil, err
	}

	return s.toResponse(holiday), nil
}

func (s *holidayService) GetByID(ctx context.Context, id string) (*dto.HolidayResponse, error) {
	holiday, err := s.holidayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(holiday), nil
}

func (s *holidayService) GetAllByYear(ctx context.Context, year int) ([]dto.HolidayResponse, error) {
	holidays, err := s.holidayRepo.FindAll(ctx, year)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.HolidayResponse, len(holidays))
	for i, holiday := range holidays {
		responses[i] = *s.toResponse(holiday)
	}

	return responses, nil
}

func (s *holidayService) Update(ctx context.Context, id string, req *dto.UpdateHolidayRequest) (*dto.HolidayResponse, error) {
	holiday, err := s.holidayRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		holiday.Name = *req.Name
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		holiday.Date = date
	}
	if req.Type != nil {
		holiday.Type = *req.Type
	}
	if req.IsRecurring != nil {
		holiday.IsRecurring = *req.IsRecurring
	}
	if req.Year != nil {
		holiday.Year = req.Year
	}
	if req.Description != nil {
		holiday.Description = req.Description
	}
	holiday.UpdatedAt = time.Now()

	if err := s.holidayRepo.Update(ctx, holiday); err != nil {
		return nil, err
	}

	return s.toResponse(holiday), nil
}

func (s *holidayService) Delete(ctx context.Context, id string) error {
	return s.holidayRepo.Delete(ctx, id)
}

func (s *holidayService) toResponse(holiday *entity.Holiday) *dto.HolidayResponse {
	return &dto.HolidayResponse{
		ID:          holiday.ID,
		Name:        holiday.Name,
		Date:        holiday.Date.Format("2006-01-02"),
		Type:        holiday.Type,
		IsRecurring: holiday.IsRecurring,
		Year:        holiday.Year,
		Description: holiday.Description,
		CreatedAt:   holiday.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   holiday.UpdatedAt.Format(time.RFC3339),
	}
}
