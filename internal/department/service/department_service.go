package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"hris/internal/department/dto"
	"hris/internal/department/entity"
	"hris/internal/department/repository"
	employeerepo "hris/internal/employee/repository"
)

type DepartmentService interface {
	Create(ctx context.Context, req *dto.CreateDepartmentRequest, companyID string) (*dto.DepartmentResponse, error)
	GetByID(ctx context.Context, id string, companyID string) (*dto.DepartmentResponse, error)
	GetAll(ctx context.Context, companyID string) ([]dto.DepartmentResponse, error)
	Update(ctx context.Context, id string, companyID string, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	Delete(ctx context.Context, id string, companyID string) error
}

type departmentService struct {
	departmentRepo repository.DepartmentRepository
	employeeRepo   employeerepo.EmployeeRepository
}

func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
	employeeRepo employeerepo.EmployeeRepository,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

func (s *departmentService) Create(ctx context.Context, req *dto.CreateDepartmentRequest, companyID string) (*dto.DepartmentResponse, error) {
	// Check if code already exists within the company
	if _, err := s.departmentRepo.FindByCodeAndCompany(ctx, req.Code, companyID); err == nil {
		return nil, errors.New("department code already exists")
	}

	// Validate head employee if provided (check employee belongs to same company)
	if req.HeadEmployeeID != nil && *req.HeadEmployeeID != "" {
		headEmployeeUUID, err := uuid.Parse(*req.HeadEmployeeID)
		if err != nil {
			return nil, errors.New("invalid head employee ID format")
		}
		if _, err := s.employeeRepo.FindByIDAndCompany(ctx, headEmployeeUUID, companyID); err != nil {
			return nil, errors.New("head employee not found")
		}
	}

	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	department := &entity.Department{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		HeadEmployeeID: req.HeadEmployeeID,
		IsActive:       isActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return nil, err
	}

	return s.toResponse(ctx, department, nil), nil
}

func (s *departmentService) GetByID(ctx context.Context, id string, companyID string) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		return nil, err
	}

	deptUUID, _ := uuid.Parse(department.ID)
	counts, err := s.employeeRepo.CountByDepartmentIDs(ctx, []uuid.UUID{deptUUID}, companyID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(ctx, department, counts), nil
}

func (s *departmentService) GetAll(ctx context.Context, companyID string) ([]dto.DepartmentResponse, error) {
	departments, err := s.departmentRepo.FindAllByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	deptIDs := make([]uuid.UUID, len(departments))
	for i, dept := range departments {
		deptIDs[i], _ = uuid.Parse(dept.ID)
	}

	counts, err := s.employeeRepo.CountByDepartmentIDs(ctx, deptIDs, companyID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DepartmentResponse, len(departments))
	for i, dept := range departments {
		responses[i] = *s.toResponse(ctx, &dept, counts)
	}

	return responses, nil
}

func (s *departmentService) Update(ctx context.Context, id string, companyID string, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	department, err := s.departmentRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		return nil, err
	}

	// Check if new code already exists within the company (if code is being changed)
	if req.Code != nil && *req.Code != department.Code {
		if _, err := s.departmentRepo.FindByCodeAndCompany(ctx, *req.Code, companyID); err == nil {
			return nil, errors.New("department code already exists")
		}
		department.Code = *req.Code
	}

	// Validate head employee if provided (check employee belongs to same company)
	if req.HeadEmployeeID != nil {
		if *req.HeadEmployeeID == "" {
			department.HeadEmployeeID = nil
		} else {
			if _, err := s.employeeRepo.FindByIDAndCompany(ctx, uuid.MustParse(*req.HeadEmployeeID), companyID); err != nil {
				return nil, errors.New("head employee not found")
			}
			department.HeadEmployeeID = req.HeadEmployeeID
		}
	}

	if req.Name != nil {
		department.Name = *req.Name
	}
	if req.Description != nil {
		department.Description = req.Description
	}
	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}
	department.UpdatedAt = time.Now()

	if err := s.departmentRepo.Update(ctx, department); err != nil {
		return nil, err
	}

	return s.toResponse(ctx, department, nil), nil
}

func (s *departmentService) Delete(ctx context.Context, id string, companyID string) error {
	// Check if department exists and belongs to company
	_, err := s.departmentRepo.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		return err
	}

	// TODO: Check if department has employees before deleting
	// For now, allow deletion

	return s.departmentRepo.DeleteByCompany(ctx, id, companyID)
}

func (s *departmentService) toResponse(ctx context.Context, dept *entity.Department, employeeCounts map[uuid.UUID]int) *dto.DepartmentResponse {
	response := &dto.DepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Code:        dept.Code,
		Description: dept.Description,
		IsActive:    dept.IsActive,
		CreatedAt:   dept.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   dept.UpdatedAt.Format(time.RFC3339),
	}

	if employeeCounts != nil {
		if deptUUID, err := uuid.Parse(dept.ID); err == nil {
			response.EmployeeCount = employeeCounts[deptUUID]
		}
	}

	response.HeadEmployeeID = dept.HeadEmployeeID

	if dept.HeadEmployeeID != nil && *dept.HeadEmployeeID != "" {
		if headEmployeeUUID, err := uuid.Parse(*dept.HeadEmployeeID); err == nil {
			if emp, err := s.employeeRepo.FindByID(ctx, headEmployeeUUID); err == nil {
				response.HeadEmployeeName = &emp.Position
			}
		}
	}

	return response
}
