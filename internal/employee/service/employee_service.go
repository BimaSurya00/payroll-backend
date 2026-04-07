package service

import (
	"context"

	"hris/internal/employee/dto"
	"hris/internal/employee/helper"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req *dto.CreateEmployeeRequest, companyID string) (*dto.EmployeeResponse, error)
	GetAllEmployees(ctx context.Context, page, perPage int, path string, search string, companyID string) (*helper.Pagination[*dto.EmployeeResponse], error)
	GetEmployeeByID(ctx context.Context, id string, companyID string) (*dto.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, id string, companyID string, req *dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, id string, companyID string) error

	// Self-service endpoints
	GetMyProfile(ctx context.Context, userID string) (*dto.EmployeeResponse, error)
	UpdateMyProfile(ctx context.Context, userID string, req *dto.UpdateMyProfileRequest) (*dto.EmployeeResponse, error)
}
