package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"hris/internal/attendance/entity"
)

type AttendanceFilter struct {
	EmployeeID *uuid.UUID
	ScheduleID *uuid.UUID
	Status     *string
	DateFrom   *time.Time
	DateTo     *time.Time
	CompanyID  *string
}

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.Attendance) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error)
	FindByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*entity.Attendance, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, skip, limit int64) ([]*entity.Attendance, error)
	CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error)
	FindAll(ctx context.Context, filter AttendanceFilter, skip, limit int64) ([]*entity.Attendance, error)
	CountAll(ctx context.Context, filter AttendanceFilter) (int64, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	CountLateByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error)
	CountAbsentByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error)
	GetAttendanceSummaryByPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*AttendanceSummary, error)
	CountByStatusAndDate(ctx context.Context, status string, date time.Time) (int64, error)
	GetMonthlySummaryAll(ctx context.Context, startDate, endDate time.Time) ([]AttendanceMonthlySummary, error)
	GetMonthlySummaryByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) ([]AttendanceMonthlySummary, error)
}

type AttendanceMonthlySummary struct {
	EmployeeID   uuid.UUID
	TotalPresent int
	TotalLate    int
	TotalAbsent  int
	TotalLeave   int
	TotalDays    int
}

type AttendanceSummary struct {
	TotalPresent int
	TotalLate    int
	TotalAbsent  int
	TotalLeave   int
	TotalDays    int
}
