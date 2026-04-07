package service

import (
	"context"

	"hris/internal/payroll/dto"
	"hris/internal/payroll/helper"
)

type PayrollService interface {
	GenerateBulk(ctx context.Context, req *dto.GeneratePayrollRequest, companyID string) (*dto.GeneratePayrollResponse, error)
	GetByID(ctx context.Context, id string, companyID string) (*dto.PayrollResponse, error)
	GetAll(ctx context.Context, page, perPage int, path string, companyID string) (*helper.PayrollPagination, error)
	UpdateStatus(ctx context.Context, id string, companyID string, req *dto.UpdatePayrollStatusRequest) error
	ExportCSV(ctx context.Context, month, year int, companyID string) ([]byte, string, error)
	GetMyPayrolls(ctx context.Context, userID string, page, perPage int, path string) (*helper.PayrollPagination, error)
	GetMyPayrollByID(ctx context.Context, userID string, payrollID string) (*dto.PayrollResponse, error)
}
