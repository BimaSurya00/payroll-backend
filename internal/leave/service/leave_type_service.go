package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"hris/internal/leave/dto"
	"hris/internal/leave/entity"
	leaverepo "hris/internal/leave/repository"
)

var (
	ErrLeaveTypeNotFound = errors.New("leave type not found")
)

type LeaveTypeService interface {
	GetActiveLeaveTypes(ctx context.Context) ([]dto.LeaveTypeResponse, error)
	GetAllLeaveTypes(ctx context.Context) ([]dto.LeaveTypeResponse, error)
	GetLeaveTypeByID(ctx context.Context, id string) (*dto.LeaveTypeResponse, error)
	CreateLeaveType(ctx context.Context, req *dto.CreateLeaveTypeRequest) (*dto.LeaveTypeResponse, error)
	UpdateLeaveType(ctx context.Context, id string, req *dto.UpdateLeaveTypeRequest) (*dto.LeaveTypeResponse, error)
	DeleteLeaveType(ctx context.Context, id string) error
}

type leaveTypeService struct {
	leaveTypeRepo leaverepo.LeaveTypeRepository
}

func NewLeaveTypeService(leaveTypeRepo leaverepo.LeaveTypeRepository) LeaveTypeService {
	return &leaveTypeService{
		leaveTypeRepo: leaveTypeRepo,
	}
}

func (s *leaveTypeService) GetActiveLeaveTypes(ctx context.Context) ([]dto.LeaveTypeResponse, error) {
	leaveTypes, err := s.leaveTypeRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LeaveTypeResponse, len(leaveTypes))
	for i, lt := range leaveTypes {
		responses[i] = s.toLeaveTypeResponse(&lt)
	}

	return responses, nil
}

func (s *leaveTypeService) GetAllLeaveTypes(ctx context.Context) ([]dto.LeaveTypeResponse, error) {
	leaveTypes, err := s.leaveTypeRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LeaveTypeResponse, len(leaveTypes))
	for i, lt := range leaveTypes {
		responses[i] = s.toLeaveTypeResponse(&lt)
	}

	return responses, nil
}

func (s *leaveTypeService) GetLeaveTypeByID(ctx context.Context, id string) (*dto.LeaveTypeResponse, error) {
	leaveTypeUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrLeaveTypeNotFound
	}

	leaveType, err := s.leaveTypeRepo.FindByID(ctx, leaveTypeUUID)
	if err != nil {
		if errors.Is(err, leaverepo.ErrLeaveTypeNotFound) {
			return nil, ErrLeaveTypeNotFound
		}
		return nil, err
	}

	response := s.toLeaveTypeResponse(leaveType)
	return &response, nil
}

func (s *leaveTypeService) CreateLeaveType(ctx context.Context, req *dto.CreateLeaveTypeRequest) (*dto.LeaveTypeResponse, error) {
	// Use defaultDays if provided, otherwise use maxDaysPerYear
	maxDays := req.MaxDaysPerYear
	if maxDays == 0 && req.DefaultDays > 0 {
		maxDays = req.DefaultDays
	}

	leaveType := &entity.LeaveType{
		ID:               uuid.New().String(),
		Name:             req.Name,
		Code:             req.Code,
		Description:      &req.Description,
		MaxDaysPerYear:   maxDays,
		DefaultDays:      req.DefaultDays,
		IsPaid:           req.IsPaid,
		RequiresApproval: req.RequiresApproval,
		IsActive:         req.IsActive,
		Color:            req.Color,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.leaveTypeRepo.Create(ctx, leaveType); err != nil {
		return nil, err
	}

	response := s.toLeaveTypeResponse(leaveType)
	return &response, nil
}

func (s *leaveTypeService) UpdateLeaveType(ctx context.Context, id string, req *dto.UpdateLeaveTypeRequest) (*dto.LeaveTypeResponse, error) {
	leaveTypeUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrLeaveTypeNotFound
	}

	leaveType, err := s.leaveTypeRepo.FindByID(ctx, leaveTypeUUID)
	if err != nil {
		if errors.Is(err, leaverepo.ErrLeaveTypeNotFound) {
			return nil, ErrLeaveTypeNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		leaveType.Name = *req.Name
	}
	if req.Description != nil {
		leaveType.Description = req.Description
	}
	if req.MaxDaysPerYear != nil {
		leaveType.MaxDaysPerYear = *req.MaxDaysPerYear
	}
	if req.IsPaid != nil {
		leaveType.IsPaid = *req.IsPaid
	}
	if req.RequiresApproval != nil {
		leaveType.RequiresApproval = *req.RequiresApproval
	}
	if req.IsActive != nil {
		leaveType.IsActive = *req.IsActive
	}

	if err := s.leaveTypeRepo.Update(ctx, leaveType); err != nil {
		return nil, err
	}

	response := s.toLeaveTypeResponse(leaveType)
	return &response, nil
}

func (s *leaveTypeService) DeleteLeaveType(ctx context.Context, id string) error {
	leaveTypeUUID, err := uuid.Parse(id)
	if err != nil {
		return ErrLeaveTypeNotFound
	}

	if err := s.leaveTypeRepo.Delete(ctx, leaveTypeUUID); err != nil {
		if errors.Is(err, leaverepo.ErrLeaveTypeNotFound) {
			return ErrLeaveTypeNotFound
		}
		return err
	}

	return nil
}

func (s *leaveTypeService) toLeaveTypeResponse(leaveType *entity.LeaveType) dto.LeaveTypeResponse {
	description := ""
	if leaveType.Description != nil {
		description = *leaveType.Description
	}

	return dto.LeaveTypeResponse{
		ID:               leaveType.ID,
		Name:             leaveType.Name,
		Code:             leaveType.Code,
		Description:      description,
		MaxDaysPerYear:   leaveType.MaxDaysPerYear,
		DefaultDays:      leaveType.DefaultDays,
		IsPaid:           leaveType.IsPaid,
		RequiresApproval: leaveType.RequiresApproval,
		IsActive:         leaveType.IsActive,
		Color:            leaveType.Color,
		CreatedAt:        leaveType.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        leaveType.UpdatedAt.Format(time.RFC3339),
	}
}
