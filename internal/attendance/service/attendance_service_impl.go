package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hris/internal/attendance/dto"
	attendanceEntity "hris/internal/attendance/entity"
	"hris/internal/attendance/helper"
	"hris/internal/attendance/repository"
	employeeRepo "hris/internal/employee/repository"
	scheduleEntity "hris/internal/schedule/entity"
	scheduleRepo "hris/internal/schedule/repository"
	sharedHelper "hris/shared/helper"
)

type attendanceService struct {
	attendanceRepo repository.AttendanceRepository
	employeeRepo   employeeRepo.EmployeeRepository
	scheduleRepo   scheduleRepo.ScheduleRepository
}

func NewAttendanceService(
	attendanceRepo repository.AttendanceRepository,
	employeeRepo employeeRepo.EmployeeRepository,
	scheduleRepo scheduleRepo.ScheduleRepository,
) AttendanceService {
	return &attendanceService{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
		scheduleRepo:   scheduleRepo,
	}
}

func (s *attendanceService) ClockIn(ctx context.Context, userID string, req *dto.ClockInRequest) (*dto.ClockInResponse, error) {
	// 1. Get employee by userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeeRepo.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// 2. Get schedule
	var schedule *scheduleEntity.Schedule
	if employee.ScheduleID != nil {
		schedule, err = s.scheduleRepo.FindByID(ctx, *employee.ScheduleID)
		if err != nil {
			if errors.Is(err, scheduleRepo.ErrScheduleNotFound) {
				return nil, ErrScheduleNotFound
			}
			return nil, fmt.Errorf("failed to get schedule: %w", err)
		}
	} else {
		return nil, ErrScheduleNotFound
	}

	// 3. Check if already clocked in today
	now := sharedHelper.Now()
	today := sharedHelper.Today()

	_, err = s.attendanceRepo.FindByEmployeeAndDate(ctx, employee.ID, today)
	if err == nil {
		return nil, ErrAlreadyClockedIn
	}
	if !errors.Is(err, repository.ErrAttendanceNotFound) {
		return nil, fmt.Errorf("failed to check today's attendance: %w", err)
	}

	// 4. Validate GPS location
	distance := helper.CalculateDistance(req.Lat, req.Long, schedule.OfficeLat, schedule.OfficeLong)
	if !helper.IsWithinRadius(req.Lat, req.Long, schedule.OfficeLat, schedule.OfficeLong, float64(schedule.AllowedRadiusMeters)) {
		return nil, ErrOutOfOfficeRange
	}

	// 5. Determine status (PRESENT/LATE)
	status := s.determineStatus(now, schedule.TimeIn, schedule.AllowedLateMinutes)

	// 6. Create attendance
	attendanceID := uuid.New()
	attendance := &attendanceEntity.Attendance{
		ID:          attendanceID.String(),
		EmployeeID:  employee.ID.String(),
		ScheduleID:  &schedule.ID,
		Date:        today,
		ClockInTime: &now,
		ClockInLat:  &req.Lat,
		ClockInLong: &req.Long,
		Status:      status,
		Notes:       "",
		CreatedAt:   now,
	}

	if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to create attendance: %w", err)
	}

	return &dto.ClockInResponse{
		AttendanceID: attendance.ID,
		EmployeeID:   employee.ID.String(),
		ClockInTime:  now.Format(time.RFC3339),
		Status:       status,
		Distance:     distance,
		ScheduleName: schedule.Name,
	}, nil
}

func (s *attendanceService) ClockOut(ctx context.Context, userID string, req *dto.ClockOutRequest) (*dto.ClockOutResponse, error) {
	// 1. Get employee
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeeRepo.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// 2. Get today's attendance
	now := sharedHelper.Now()
	today := sharedHelper.Today()

	attendance, err := s.attendanceRepo.FindByEmployeeAndDate(ctx, employee.ID, today)
	if err != nil {
		if errors.Is(err, repository.ErrAttendanceNotFound) {
			return nil, ErrNotClockedIn
		}
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}

	// 3. Validate GPS location (optional - get schedule)
	var distance float64
	if attendance.ScheduleID != nil {
		scheduleUUID, _ := uuid.Parse(*attendance.ScheduleID)
		schedule, err := s.scheduleRepo.FindByID(ctx, scheduleUUID)
		if err == nil {
			distance = helper.CalculateDistance(req.Lat, req.Long, schedule.OfficeLat, schedule.OfficeLong)
		}
	}

	// 4. Update clock out
	updates := map[string]interface{}{
		"clock_out_time": now,
		"clock_out_lat":  req.Lat,
		"clock_out_long": req.Long,
	}

	attendanceUUID, _ := uuid.Parse(attendance.ID)
	if err := s.attendanceRepo.Update(ctx, attendanceUUID, updates); err != nil {
		return nil, fmt.Errorf("failed to update attendance: %w", err)
	}

	return &dto.ClockOutResponse{
		AttendanceID: attendance.ID,
		ClockOutTime: now.Format(time.RFC3339),
		Distance:     distance,
	}, nil
}

func (s *attendanceService) GetHistory(ctx context.Context, userID string, page, perPage int, path string) (*Pagination[*dto.AttendanceResponse], error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeeRepo.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	attendances, err := s.attendanceRepo.FindByEmployeeID(ctx, employee.ID, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	total, err := s.attendanceRepo.CountByEmployeeID(ctx, employee.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count attendances: %w", err)
	}

	attendanceResponses := make([]*dto.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		attendanceResponses[i] = helper.ToAttendanceResponse(attendance)
	}

	pagination := NewPagination(attendanceResponses, page, perPage, total, path)

	return pagination, nil
}

func (s *attendanceService) determineStatus(clockInTime time.Time, scheduleTimeIn string, allowedLateMinutes int) string {
	scheduleTime, err := time.Parse("15:04", scheduleTimeIn)
	if err != nil {
		return "PRESENT"
	}

	loc := sharedHelper.GetLocation()
	localClockIn := clockInTime.In(loc)

	scheduledTime := time.Date(localClockIn.Year(), localClockIn.Month(), localClockIn.Day(),
		scheduleTime.Hour(), scheduleTime.Minute(), 0, 0, loc)

	deadline := scheduledTime.Add(time.Duration(allowedLateMinutes) * time.Minute)

	if localClockIn.After(deadline) {
		return "LATE"
	}

	return "PRESENT"
}

func (s *attendanceService) GetAllAttendances(ctx context.Context, filter GetAllAttendanceFilter, page, perPage int, path string, companyID string) (*Pagination[*dto.AttendanceResponse], error) {
	// Convert service filter to repository filter with company ID
	repoFilter := repository.AttendanceFilter{
		EmployeeID: filter.EmployeeID,
		ScheduleID: filter.ScheduleID,
		Status:     filter.Status,
		DateFrom:   filter.DateFrom,
		DateTo:     filter.DateTo,
		CompanyID:  &companyID,
	}

	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	attendances, err := s.attendanceRepo.FindAll(ctx, repoFilter, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}

	total, err := s.attendanceRepo.CountAll(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count attendances: %w", err)
	}

	attendanceResponses := make([]*dto.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		attendanceResponses[i] = helper.ToAttendanceResponse(attendance)
	}

	pagination := NewPagination(attendanceResponses, page, perPage, total, path)

	return pagination, nil
}

func (s *attendanceService) GetMonthlyReport(ctx context.Context, month, year int, companyID string) (*dto.MonthlyAttendanceReport, error) {
	// Validate month and year
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2020 || year > 2100 {
		return nil, errors.New("invalid year")
	}

	// Calculate date range
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// Get summaries for employees in the company
	summaries, err := s.attendanceRepo.GetMonthlySummaryByCompany(ctx, companyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Build report items
	items := make([]dto.AttendanceReportItem, 0, len(summaries))
	var totalPresent, totalLate, totalAbsent, totalLeave, totalDays int
	var totalAttendanceRate float64

	for _, summary := range summaries {
		employee, err := s.employeeRepo.FindByID(ctx, summary.EmployeeID)
		if err != nil {
			continue // Skip if employee not found
		}

		attendanceRate := dto.CalculateAttendanceRate(summary.TotalPresent, summary.TotalLate, summary.TotalDays)

		item := dto.AttendanceReportItem{
			EmployeeID:     summary.EmployeeID.String(),
			EmployeeName:   employee.Position, // Using Position as name
			Position:       employee.Position,
			Division:       employee.Division,
			TotalPresent:   summary.TotalPresent,
			TotalLate:      summary.TotalLate,
			TotalAbsent:    summary.TotalAbsent,
			TotalLeave:     summary.TotalLeave,
			TotalDays:      summary.TotalDays,
			AttendanceRate: attendanceRate,
			Period:         startDate.Format("2006-01"),
		}
		items = append(items, item)

		// Aggregate totals
		totalPresent += summary.TotalPresent
		totalLate += summary.TotalLate
		totalAbsent += summary.TotalAbsent
		totalLeave += summary.TotalLeave
		totalDays += summary.TotalDays
	}

	if len(items) > 0 {
		totalAttendanceRate = dto.CalculateAttendanceRate(totalPresent, totalLate, totalDays)
	}

	// Build aggregate summary
	summary := dto.AttendanceReportItem{
		EmployeeID:     "",
		EmployeeName:   "All Employees",
		Position:       "",
		Division:       "",
		TotalPresent:   totalPresent,
		TotalLate:      totalLate,
		TotalAbsent:    totalAbsent,
		TotalLeave:     totalLeave,
		TotalDays:      totalDays,
		AttendanceRate: totalAttendanceRate,
		Period:         startDate.Format("2006-01"),
	}

	return &dto.MonthlyAttendanceReport{
		Period:         startDate.Format("2006-01"),
		Month:          month,
		Year:           year,
		TotalEmployees: len(items),
		Summary:        summary,
		Items:          items,
	}, nil
}

func (s *attendanceService) GetMyMonthlySummary(ctx context.Context, userID string, month, year int) (*dto.MyAttendanceSummary, error) {
	// Validate month and year
	if month < 1 || month > 12 {
		return nil, errors.New("invalid month")
	}
	if year < 2020 || year > 2100 {
		return nil, errors.New("invalid year")
	}

	// Get employee by user ID
	employee, err := s.employeeRepo.FindByUserID(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, ErrEmployeeNotFound
	}

	// Calculate date range
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// Get attendance summary
	summary, err := s.attendanceRepo.GetAttendanceSummaryByPeriod(ctx, employee.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	attendanceRate := dto.CalculateAttendanceRate(summary.TotalPresent, summary.TotalLate, summary.TotalDays)

	return &dto.MyAttendanceSummary{
		Period:         startDate.Format("2006-01"),
		TotalPresent:   summary.TotalPresent,
		TotalLate:      summary.TotalLate,
		TotalAbsent:    summary.TotalAbsent,
		TotalLeave:     summary.TotalLeave,
		TotalDays:      summary.TotalDays,
		AttendanceRate: attendanceRate,
	}, nil
}

func (s *attendanceService) CreateCorrection(ctx context.Context, adminID string, req *dto.CreateCorrectionRequest, companyID string) (*dto.AttendanceResponse, error) {
	// Parse employee ID
	employeeUUID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		return nil, errors.New("invalid employee ID")
	}

	// Verify employee belongs to the company
	_, err = s.employeeRepo.FindByIDAndCompany(ctx, employeeUUID, companyID)
	if err != nil {
		return nil, errors.New("employee not found or does not belong to your company")
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Check if attendance already exists for this employee on this date
	existing, _ := s.attendanceRepo.FindByEmployeeAndDate(ctx, employeeUUID, date)
	if existing != nil {
		return nil, errors.New("attendance already exists for this date, use update correction instead")
	}

	// Store company ID in attendance
	_ = companyID // Use companyID to avoid unused variable error

	// Parse clock in time
	clockInTime, err := time.Parse("15:04", req.ClockIn)
	if err != nil {
		return nil, errors.New("invalid clock in time format")
	}
	clockInDateTime := time.Date(date.Year(), date.Month(), date.Day(), clockInTime.Hour(), clockInTime.Minute(), 0, 0, time.UTC)

	// Parse clock out time if provided
	var clockOutDateTime *time.Time
	if req.ClockOut != "" {
		clockOutTime, err := time.Parse("15:04", req.ClockOut)
		if err != nil {
			return nil, errors.New("invalid clock out time format")
		}
		dt := time.Date(date.Year(), date.Month(), date.Day(), clockOutTime.Hour(), clockOutTime.Minute(), 0, 0, time.UTC)
		clockOutDateTime = &dt
	}

	// Create attendance with correction info
	now := time.Now()
	attendance := &attendanceEntity.Attendance{
		ID:             uuid.New().String(),
		EmployeeID:     employeeUUID.String(),
		Date:           date,
		ClockInTime:    &clockInDateTime,
		ClockOutTime:   clockOutDateTime,
		Status:         req.Status,
		Notes:          req.Notes,
		CorrectedBy:    &adminID,
		CorrectedAt:    &now,
		CorrectionNote: &req.Notes,
		CreatedAt:      now,
	}

	if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
		return nil, err
	}

	return helper.ToAttendanceResponse(attendance), nil
}

func (s *attendanceService) UpdateCorrection(ctx context.Context, adminID, attendanceID string, req *dto.UpdateCorrectionRequest, companyID string) (*dto.AttendanceResponse, error) {
	// Parse attendance ID
	attendanceUUID, err := uuid.Parse(attendanceID)
	if err != nil {
		return nil, errors.New("invalid attendance ID")
	}

	// Find existing attendance
	attendance, err := s.attendanceRepo.FindByID(ctx, attendanceUUID)
	if err != nil {
		return nil, repository.ErrAttendanceNotFound
	}

	// Verify the attendance belongs to an employee in the company
	employeeUUID, _ := uuid.Parse(attendance.EmployeeID)
	_, err = s.employeeRepo.FindByIDAndCompany(ctx, employeeUUID, companyID)
	if err != nil {
		return nil, errors.New("attendance not found or does not belong to your company")
	}

	// Update fields if provided
	updates := make(map[string]interface{})

	if req.ClockIn != nil {
		// Parse and update clock in time
		clockInTime, err := time.Parse("15:04", *req.ClockIn)
		if err != nil {
			return nil, errors.New("invalid clock in time format")
		}
		newClockInTime := time.Date(attendance.Date.Year(), attendance.Date.Month(), attendance.Date.Day(), clockInTime.Hour(), clockInTime.Minute(), 0, 0, time.UTC)
		updates["clock_in_time"] = newClockInTime
	}

	if req.ClockOut != nil {
		// Parse and update clock out time
		clockOutTime, err := time.Parse("15:04", *req.ClockOut)
		if err != nil {
			return nil, errors.New("invalid clock out time format")
		}
		newClockOutTime := time.Date(attendance.Date.Year(), attendance.Date.Month(), attendance.Date.Day(), clockOutTime.Hour(), clockOutTime.Minute(), 0, 0, time.UTC)
		updates["clock_out_time"] = newClockOutTime
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	// Always update correction info
	updates["corrected_by"] = adminID
	updates["corrected_at"] = time.Now()
	if req.Notes != nil {
		updates["correction_note"] = *req.Notes
	}

	// Update attendance
	if err := s.attendanceRepo.Update(ctx, attendanceUUID, updates); err != nil {
		return nil, err
	}

	// Fetch updated attendance
	updatedAttendance, err := s.attendanceRepo.FindByID(ctx, attendanceUUID)
	if err != nil {
		return nil, err
	}

	return helper.ToAttendanceResponse(updatedAttendance), nil
}
