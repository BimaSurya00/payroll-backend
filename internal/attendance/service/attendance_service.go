package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"hris/internal/attendance/dto"
)

var (
	ErrAlreadyClockedIn = errors.New("already clocked in today")
	ErrNotClockedIn     = errors.New("not clocked in yet")
	ErrOutOfOfficeRange = errors.New("out of office range")
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrCompanyNotFound  = errors.New("company not found")
	ErrOfficeNotSet     = errors.New("office location not configured for company")
)

type GetAllAttendanceFilter struct {
	EmployeeID *uuid.UUID
	ScheduleID *uuid.UUID
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
}

type AttendanceService interface {
	ClockIn(ctx context.Context, userID string, req *dto.ClockInRequest) (*dto.ClockInResponse, error)
	ClockOut(ctx context.Context, userID string, req *dto.ClockOutRequest) (*dto.ClockOutResponse, error)
	GetHistory(ctx context.Context, userID string, page, perPage int, path string) (*Pagination[*dto.AttendanceResponse], error)
	GetAllAttendances(ctx context.Context, filter GetAllAttendanceFilter, page, perPage int, path string, companyID string) (*Pagination[*dto.AttendanceResponse], error)
	GetMonthlyReport(ctx context.Context, month, year int, companyID string) (*dto.MonthlyAttendanceReport, error)
	GetMyMonthlySummary(ctx context.Context, userID string, month, year int) (*dto.MyAttendanceSummary, error)
	CreateCorrection(ctx context.Context, adminID string, req *dto.CreateCorrectionRequest, companyID string) (*dto.AttendanceResponse, error)
	UpdateCorrection(ctx context.Context, adminID, attendanceID string, req *dto.UpdateCorrectionRequest, companyID string) (*dto.AttendanceResponse, error)
}
