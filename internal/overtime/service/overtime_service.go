package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	employeerepo "hris/internal/employee/repository"
	"hris/internal/overtime/dto"
	"hris/internal/overtime/entity"
	"hris/internal/overtime/repository"
	userRepo "hris/internal/user/repository"
)

var (
	ErrOvertimeRequestNotFound  = errors.New("overtime request not found")
	ErrInvalidTimeRange         = errors.New("invalid time range")
	ErrDuplicateOvertimeRequest = errors.New("overtime request for this date already exists")
	ErrOvertimeTooLong          = errors.New("overtime hours exceed maximum allowed")
	ErrInvalidRequestStatus     = errors.New("invalid request status for this action")
)

type OvertimeService interface {
	// Overtime Requests
	CreateOvertimeRequest(ctx context.Context, userID string, req *dto.CreateOvertimeRequestRequest) (*dto.OvertimeRequestResponse, error)
	GetMyOvertimeRequests(ctx context.Context, userID string, page, perPage int) ([]dto.OvertimeRequestResponse, int64, error)
	GetAllOvertimeRequests(ctx context.Context, page, perPage int, status, employeeID string, companyID string) ([]dto.OvertimeRequestResponse, int64, error)
	GetOvertimeRequestByID(ctx context.Context, id string, companyID string) (*dto.OvertimeRequestResponse, error)
	GetPendingOvertimeRequests(ctx context.Context, companyID string) ([]dto.OvertimeRequestResponse, error)
	ApproveOvertimeRequest(ctx context.Context, id string, approverID string, companyID string) (*dto.OvertimeRequestResponse, error)
	RejectOvertimeRequest(ctx context.Context, id string, approverID string, companyID string, rejectionReason string) (*dto.OvertimeRequestResponse, error)

	// Overtime Attendance
	ClockIn(ctx context.Context, requestID string, notes string) (*dto.OvertimeAttendanceResponse, error)
	ClockOut(ctx context.Context, requestID string, notes string) (*dto.OvertimeAttendanceResponse, error)

	// Overtime Policies
	GetActivePolicies(ctx context.Context) ([]dto.OvertimePolicyResponse, error)

	// Overtime Calculation
	CalculateOvertimePay(ctx context.Context, employeeID string, startDate, endDate string) (*dto.OvertimeCalculationResponse, error)
}

type overtimeService struct {
	overtimeRequestRepo    repository.OvertimeRequestRepository
	overtimePolicyRepo     repository.OvertimePolicyRepository
	overtimeAttendanceRepo repository.OvertimeAttendanceRepository
	employeeRepo           employeerepo.EmployeeRepository
	userRepo               userRepo.UserRepository
}

func NewOvertimeService(
	requestRepo repository.OvertimeRequestRepository,
	policyRepo repository.OvertimePolicyRepository,
	attendanceRepo repository.OvertimeAttendanceRepository,
	employeeRepo employeerepo.EmployeeRepository,
	userRepo userRepo.UserRepository,
) OvertimeService {
	return &overtimeService{
		overtimeRequestRepo:    requestRepo,
		overtimePolicyRepo:     policyRepo,
		overtimeAttendanceRepo: attendanceRepo,
		employeeRepo:           employeeRepo,
		userRepo:               userRepo,
	}
}

func (s *overtimeService) CreateOvertimeRequest(ctx context.Context, userID string, req *dto.CreateOvertimeRequestRequest) (*dto.OvertimeRequestResponse, error) {
	// Find employee by userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse and validate date
	overtimeDate, err := time.Parse("2006-01-02", req.OvertimeDate)
	if err != nil {
		return nil, ErrInvalidTimeRange
	}

	// Parse and validate times
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, ErrInvalidTimeRange
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, ErrInvalidTimeRange
	}

	// Calculate total hours
	duration := endTime.Sub(startTime)
	totalHours := duration.Hours()
	if totalHours <= 0 {
		return nil, ErrInvalidTimeRange
	}

	// Check for duplicate request
	_, err = s.overtimeRequestRepo.FindByEmployeeIDAndDate(ctx, employee.ID, overtimeDate)
	if err == nil {
		return nil, ErrDuplicateOvertimeRequest
	}

	// Get overtime policy
	policyUUID, err := uuid.Parse(req.OvertimePolicyID)
	if err != nil {
		return nil, fmt.Errorf("invalid policy ID: %w", err)
	}

	policy, err := s.overtimePolicyRepo.FindByID(ctx, policyUUID)
	if err != nil {
		return nil, fmt.Errorf("overtime policy not found: %w", err)
	}

	if !policy.IsActive {
		return nil, errors.New("overtime policy is not active")
	}

	// Validate against policy limits
	maxHours := floatValue(policy.MaxOvertimeHoursPerDay)
	if totalHours > maxHours {
		return nil, ErrOvertimeTooLong
	}

	// Create overtime request
	overtimeRequest := &entity.OvertimeRequest{
		ID:               uuid.New().String(),
		EmployeeID:       employee.ID.String(),
		OvertimeDate:     overtimeDate,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		TotalHours:       totalHours,
		Reason:           req.Reason,
		OvertimePolicyID: req.OvertimePolicyID,
		Status:           "PENDING",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.overtimeRequestRepo.Create(ctx, overtimeRequest); err != nil {
		return nil, fmt.Errorf("failed to create overtime request: %w", err)
	}

	return s.toOvertimeRequestResponse(overtimeRequest, user.Name, user.Email, employee.Position, policy), nil
}

func (s *overtimeService) GetMyOvertimeRequests(ctx context.Context, userID string, page, perPage int) ([]dto.OvertimeRequestResponse, int64, error) {
	// Find employee by userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("employee not found: %w", err)
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("user not found: %w", err)
	}

	requests, err := s.overtimeRequestRepo.FindByEmployeeID(ctx, employee.ID, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.OvertimeRequestResponse, len(requests))
	for i, req := range requests {
		// Parse policy UUID safely
		var policy *entity.OvertimePolicy
		if policyUUID, err := uuid.Parse(req.OvertimePolicyID); err == nil {
			policy, _ = s.overtimePolicyRepo.FindByID(ctx, policyUUID)
		}
		responses[i] = *s.toOvertimeRequestResponse(&req, user.Name, user.Email, employee.Position, policy)
	}

	total := int64(len(responses))

	return responses, total, nil
}

func (s *overtimeService) GetAllOvertimeRequests(ctx context.Context, page, perPage int, status, employeeID string, companyID string) ([]dto.OvertimeRequestResponse, int64, error) {
	// Filter by company
	var employeeFilter uuid.UUID = uuid.Nil

	if employeeID != "" {
		employeeFilter, _ = uuid.Parse(employeeID)
	}

	requests, err := s.overtimeRequestRepo.FindAllByCompany(ctx, companyID, page, perPage, status, employeeFilter)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.OvertimeRequestResponse, len(requests))
	for i, req := range requests {
		// Get employee and user info
		employeeUUID, _ := uuid.Parse(req.EmployeeID)
		employee, _ := s.employeeRepo.FindByID(ctx, employeeUUID)

		var userName, userEmail, position string
		if employee != nil {
			user, _ := s.userRepo.FindByID(ctx, employee.UserID.String())
			if user != nil {
				userName = user.Name
				userEmail = user.Email
			}
			position = employee.Position
		}

		// Parse policy UUID safely
		var policy *entity.OvertimePolicy
		if policyUUID, err := uuid.Parse(req.OvertimePolicyID); err == nil {
			policy, _ = s.overtimePolicyRepo.FindByID(ctx, policyUUID)
		}
		responses[i] = *s.toOvertimeRequestResponse(&req, userName, userEmail, position, policy)
	}

	total := int64(len(responses))

	return responses, total, nil
}

func (s *overtimeService) GetOvertimeRequestByID(ctx context.Context, id string, companyID string) (*dto.OvertimeRequestResponse, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	request, err := s.overtimeRequestRepo.FindByIDAndCompany(ctx, requestUUID, companyID)
	if err != nil {
		return nil, ErrOvertimeRequestNotFound
	}

	// Get employee first
	employeeUUID, _ := uuid.Parse(request.EmployeeID)
	employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Get user using employee.UserID
	user, err := s.userRepo.FindByID(ctx, employee.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse policy UUID safely
	var policy *entity.OvertimePolicy
	if policyUUID, err := uuid.Parse(request.OvertimePolicyID); err == nil {
		policy, _ = s.overtimePolicyRepo.FindByID(ctx, policyUUID)
	}

	return s.toOvertimeRequestResponse(request, user.Name, user.Email, employee.Position, policy), nil
}

func (s *overtimeService) GetPendingOvertimeRequests(ctx context.Context, companyID string) ([]dto.OvertimeRequestResponse, error) {
	requests, err := s.overtimeRequestRepo.FindPendingByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OvertimeRequestResponse, 0, len(requests))
	for _, req := range requests {
		// Get employee first
		employeeUUID, err := uuid.Parse(req.EmployeeID)
		if err != nil {
			continue
		}

		employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
		if err != nil {
			continue
		}

		// Get user using employee.UserID
		user, err := s.userRepo.FindByID(ctx, employee.UserID.String())
		if err != nil {
			continue
		}

		// Parse policy UUID safely, skip if invalid
		policyUUID, err := uuid.Parse(req.OvertimePolicyID)
		var policy *entity.OvertimePolicy
		if err == nil {
			policy, _ = s.overtimePolicyRepo.FindByID(ctx, policyUUID)
		}

		responses = append(responses, *s.toOvertimeRequestResponse(&req, user.Name, user.Email, employee.Position, policy))
	}

	return responses, nil
}

func (s *overtimeService) ApproveOvertimeRequest(ctx context.Context, id string, approverID string, companyID string) (*dto.OvertimeRequestResponse, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	request, err := s.overtimeRequestRepo.FindByIDAndCompany(ctx, requestUUID, companyID)
	if err != nil {
		return nil, ErrOvertimeRequestNotFound
	}

	if request.Status != "PENDING" {
		return nil, ErrInvalidRequestStatus
	}

	approverUUID, err := uuid.Parse(approverID)
	if err != nil {
		return nil, fmt.Errorf("invalid approver ID: %w", err)
	}

	// Update request status
	if err := s.overtimeRequestRepo.UpdateStatus(ctx, requestUUID, "APPROVED", &approverUUID, nil); err != nil {
		return nil, fmt.Errorf("failed to approve request: %w", err)
	}

	// Fetch updated request
	updatedRequest, err := s.overtimeRequestRepo.FindByID(ctx, requestUUID)
	if err != nil {
		return nil, err
	}

	// Get employee first
	employeeUUID, _ := uuid.Parse(updatedRequest.EmployeeID)
	employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Get user using employee.UserID
	user, err := s.userRepo.FindByID(ctx, employee.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse policy UUID safely
	var policy *entity.OvertimePolicy
	if policyUUID, err := uuid.Parse(updatedRequest.OvertimePolicyID); err == nil {
		policy, _ = s.overtimePolicyRepo.FindByID(ctx, policyUUID)
	}

	return s.toOvertimeRequestResponse(updatedRequest, user.Name, user.Email, employee.Position, policy), nil
}

func (s *overtimeService) RejectOvertimeRequest(ctx context.Context, id string, approverID string, companyID string, rejectionReason string) (*dto.OvertimeRequestResponse, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	request, err := s.overtimeRequestRepo.FindByIDAndCompany(ctx, requestUUID, companyID)
	if err != nil {
		return nil, ErrOvertimeRequestNotFound
	}

	if request.Status != "PENDING" {
		return nil, ErrInvalidRequestStatus
	}

	// Update request status
	if err := s.overtimeRequestRepo.UpdateStatus(ctx, requestUUID, "REJECTED", nil, &rejectionReason); err != nil {
		return nil, fmt.Errorf("failed to reject request: %w", err)
	}

	// Fetch updated request
	updatedRequest, err := s.overtimeRequestRepo.FindByID(ctx, requestUUID)
	if err != nil {
		return nil, err
	}

	// Get employee first
	employeeUUID, _ := uuid.Parse(updatedRequest.EmployeeID)
	employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Get user using employee.UserID
	user, err := s.userRepo.FindByID(ctx, employee.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse policy UUID safely
	var policy *entity.OvertimePolicy
	if policyUUID, err := uuid.Parse(updatedRequest.OvertimePolicyID); err == nil {
		policy, _ = s.overtimePolicyRepo.FindByID(ctx, policyUUID)
	}

	return s.toOvertimeRequestResponse(updatedRequest, user.Name, user.Email, employee.Position, policy), nil
}

func (s *overtimeService) ClockIn(ctx context.Context, requestID string, notes string) (*dto.OvertimeAttendanceResponse, error) {
	requestUUID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	request, err := s.overtimeRequestRepo.FindByID(ctx, requestUUID)
	if err != nil {
		return nil, ErrOvertimeRequestNotFound
	}

	if request.Status != "APPROVED" {
		return nil, errors.New("cannot clock in for unapproved request")
	}

	// Check if attendance already exists
	attendance, err := s.overtimeAttendanceRepo.FindByRequestID(ctx, requestUUID)
	if err == nil && attendance.ClockInTime != nil {
		return nil, errors.New("already clocked in")
	}

	// Create new attendance record
	now := time.Now()
	newAttendance := &entity.OvertimeAttendance{
		ID:                uuid.New().String(),
		OvertimeRequestID: requestID,
		EmployeeID:        request.EmployeeID,
		ClockInTime:       &now,
		Notes:             notes,
		CreatedAt:         now,
	}

	if err := s.overtimeAttendanceRepo.Create(ctx, newAttendance); err != nil {
		return nil, fmt.Errorf("failed to clock in: %w", err)
	}

	// Get employee and user info
	employeeUUID, _ := uuid.Parse(request.EmployeeID)
	employee, _ := s.employeeRepo.FindByID(ctx, employeeUUID)
	user, _ := s.userRepo.FindByID(ctx, employee.UserID.String())

	return s.toOvertimeAttendanceResponse(newAttendance, user.Name), nil
}

func (s *overtimeService) ClockOut(ctx context.Context, requestID string, notes string) (*dto.OvertimeAttendanceResponse, error) {
	requestUUID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	// Get attendance
	attendance, err := s.overtimeAttendanceRepo.FindByRequestID(ctx, requestUUID)
	if err != nil {
		return nil, errors.New("no clock in found")
	}

	if attendance.ClockOutTime != nil {
		return nil, errors.New("already clocked out")
	}

	// Calculate actual hours
	now := time.Now()
	actualHours := now.Sub(*attendance.ClockInTime).Hours()

	// Update clock out
	if err := s.overtimeAttendanceRepo.UpdateClockOut(ctx, uuid.MustParse(attendance.ID), now, actualHours, notes); err != nil {
		return nil, fmt.Errorf("failed to clock out: %w", err)
	}

	// Get updated attendance
	updatedAttendance, err := s.overtimeAttendanceRepo.FindByRequestID(ctx, requestUUID)
	if err != nil {
		return nil, err
	}

	// Get request to get employee ID
	request, _ := s.overtimeRequestRepo.FindByID(ctx, requestUUID)
	employeeUUID, _ := uuid.Parse(request.EmployeeID)
	employee, _ := s.employeeRepo.FindByID(ctx, employeeUUID)
	user, _ := s.userRepo.FindByID(ctx, employee.UserID.String())

	return s.toOvertimeAttendanceResponse(updatedAttendance, user.Name), nil
}

func (s *overtimeService) GetActivePolicies(ctx context.Context) ([]dto.OvertimePolicyResponse, error) {
	policies, err := s.overtimePolicyRepo.FindActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.OvertimePolicyResponse, len(policies))
	for i, policy := range policies {
		responses[i] = *s.toOvertimePolicyResponse(&policy)
	}

	return responses, nil
}

func (s *overtimeService) CalculateOvertimePay(ctx context.Context, employeeID string, startDate, endDate string) (*dto.OvertimeCalculationResponse, error) {
	// Find employee
	employeeUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Get approved overtime requests in date range
	requests, err := s.overtimeRequestRepo.FindByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}

	// Filter only approved and calculate total
	var totalHours float64
	var rateMultiplier float64
	var hourlyRate float64

	for _, req := range requests {
		if req.EmployeeID == employeeID && req.Status == "APPROVED" {
			// Parse request ID safely for attendance lookup
			var attendance *entity.OvertimeAttendance
			if reqID, err := uuid.Parse(req.ID); err == nil {
				attendance, _ = s.overtimeAttendanceRepo.FindByRequestID(ctx, reqID)
			}

			if attendance != nil && attendance.ActualHours > 0 {
				totalHours += attendance.ActualHours
			} else {
				totalHours += req.TotalHours
			}

			// Get policy for rate calculation (use first one)
			if rateMultiplier == 0 {
				if policyUUID, err := uuid.Parse(req.OvertimePolicyID); err == nil {
					policy, _ := s.overtimePolicyRepo.FindByID(ctx, policyUUID)
					if policy != nil && policy.RateMultiplier != nil {
						rateMultiplier = *policy.RateMultiplier
					}
				}
			}
		}
	}

	// Calculate hourly rate from salary base (assuming 173 working hours per month)
	hourlyRate = employee.SalaryBase / 173

	// Calculate overtime pay
	overtimePay := totalHours * hourlyRate * rateMultiplier

	return &dto.OvertimeCalculationResponse{
		EmployeeID:     employeeID,
		EmployeeName:   user.Name,
		TotalHours:     math.Round(totalHours*100) / 100,
		RateType:       "MULTIPLIER",
		RateMultiplier: rateMultiplier,
		HourlyRate:     math.Round(hourlyRate),
		OvertimePay:    math.Round(overtimePay),
	}, nil
}

// Helper functions
func (s *overtimeService) toOvertimeRequestResponse(request *entity.OvertimeRequest, userName, userEmail, employeePosition string, policy *entity.OvertimePolicy) *dto.OvertimeRequestResponse {
	var approvedAt *string
	if request.ApprovedAt != nil {
		formatted := request.ApprovedAt.Format(time.RFC3339)
		approvedAt = &formatted
	}

	return &dto.OvertimeRequestResponse{
		ID:               request.ID,
		EmployeeID:       request.EmployeeID,
		EmployeeName:     userName,
		EmployeePosition: employeePosition,
		OvertimeDate:     request.OvertimeDate.Format("2006-01-02"),
		StartTime:        request.StartTime,
		EndTime:          request.EndTime,
		TotalHours:       request.TotalHours,
		Reason:           request.Reason,
		OvertimePolicy: dto.OvertimePolicyDetail{
			ID:             policy.ID,
			Name:           policy.Name,
			RateType:       policy.RateType,
			RateMultiplier: floatValue(policy.RateMultiplier),
			FixedAmount:    floatValue(policy.FixedAmount),
		},
		Status:          request.Status,
		ApprovedBy:      request.ApprovedBy,
		ApprovedByName:  nil, // TODO: Fetch approver name if needed
		ApprovedAt:      approvedAt,
		RejectionReason: stringValue(request.RejectionReason),
		CreatedAt:       request.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       request.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *overtimeService) toOvertimePolicyResponse(policy *entity.OvertimePolicy) *dto.OvertimePolicyResponse {
	return &dto.OvertimePolicyResponse{
		ID:                       policy.ID,
		Name:                     policy.Name,
		Description:              stringValue(policy.Description),
		RateType:                 policy.RateType,
		RateMultiplier:           floatValue(policy.RateMultiplier),
		FixedAmount:              floatValue(policy.FixedAmount),
		MinOvertimeMinutes:       policy.MinOvertimeMinutes,
		MaxOvertimeHoursPerDay:   floatValue(policy.MaxOvertimeHoursPerDay),
		MaxOvertimeHoursPerMonth: floatValue(policy.MaxOvertimeHoursPerMonth),
		RequiresApproval:         policy.RequiresApproval,
		IsActive:                 policy.IsActive,
		CreatedAt:                policy.CreatedAt.Format(time.RFC3339),
		UpdatedAt:                policy.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *overtimeService) toOvertimeAttendanceResponse(attendance *entity.OvertimeAttendance, employeeName string) *dto.OvertimeAttendanceResponse {
	var clockInTime, clockOutTime *string
	if attendance.ClockInTime != nil {
		formatted := attendance.ClockInTime.Format(time.RFC3339)
		clockInTime = &formatted
	}
	if attendance.ClockOutTime != nil {
		formatted := attendance.ClockOutTime.Format(time.RFC3339)
		clockOutTime = &formatted
	}

	return &dto.OvertimeAttendanceResponse{
		ID:                attendance.ID,
		OvertimeRequestID: attendance.OvertimeRequestID,
		EmployeeID:        attendance.EmployeeID,
		EmployeeName:      employeeName,
		ClockInTime:       clockInTime,
		ClockOutTime:      clockOutTime,
		ActualHours:       attendance.ActualHours,
		Notes:             attendance.Notes,
		CreatedAt:         attendance.CreatedAt.Format(time.RFC3339),
	}
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func floatValue(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
