package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	leaveattendance "hris/internal/attendance/entity"
	attendanceRepo "hris/internal/attendance/repository"
	auditservice "hris/internal/audit/service"
	employeerepo "hris/internal/employee/repository"
	holidayrepo "hris/internal/holiday/repository"
	"hris/internal/leave/dto"
	"hris/internal/leave/entity"
	leaverepo "hris/internal/leave/repository"
	userEntity "hris/internal/user/entity"
	userRepo "hris/internal/user/repository"
	sharedHelper "hris/shared/helper"
)

var (
	ErrLeaveRequestNotFound = errors.New("leave request not found")
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrInsufficientBalance  = errors.New("insufficient leave balance")
	ErrOverlappingLeave     = errors.New("overlapping leave request exists")
	ErrInvalidRequestStatus = errors.New("invalid request status for this action")
	ErrEmployeeNotFound     = errors.New("employee not found")
)

type LeaveService interface {
	CreateLeaveRequest(ctx context.Context, employeeID string, req *dto.CreateLeaveRequestRequest) (*dto.LeaveRequestResponse, error)
	GetMyLeaveRequests(ctx context.Context, employeeID string, page, perPage int) ([]dto.LeaveRequestResponse, int64, error)
	GetAllLeaveRequests(ctx context.Context, page, perPage int, companyID string) ([]dto.LeaveRequestResponse, int64, error)
	GetLeaveRequestByID(ctx context.Context, id string, companyID string) (*dto.LeaveRequestResponse, error)
	GetPendingLeaveRequests(ctx context.Context, companyID string) ([]dto.LeaveRequestResponse, error)
	ApproveLeaveRequest(ctx context.Context, id, approverID string, companyID string, req *dto.ApproveLeaveRequest) (*dto.LeaveRequestResponse, error)
	RejectLeaveRequest(ctx context.Context, id, approverID string, companyID string, req *dto.RejectLeaveRequest) error
	GetMyLeaveBalances(ctx context.Context, employeeID string, year int) (*dto.LeaveBalanceResponse, error)
}

type leaveService struct {
	leaveTypeRepo    leaverepo.LeaveTypeRepository
	leaveBalanceRepo leaverepo.LeaveBalanceRepository
	leaveRequestRepo leaverepo.LeaveRequestRepository
	employeeRepo     employeerepo.EmployeeRepository
	userRepo         userRepo.UserRepository
	attendanceRepo   attendanceRepo.AttendanceRepository
	holidayRepo      holidayrepo.HolidayRepository
	auditService     auditservice.AuditService
	pool             *pgxpool.Pool
}

func NewLeaveService(
	leaveTypeRepo leaverepo.LeaveTypeRepository,
	leaveBalanceRepo leaverepo.LeaveBalanceRepository,
	leaveRequestRepo leaverepo.LeaveRequestRepository,
	employeeRepo employeerepo.EmployeeRepository,
	userRepo userRepo.UserRepository,
	attendanceRepo attendanceRepo.AttendanceRepository,
	holidayRepo holidayrepo.HolidayRepository,
	auditService auditservice.AuditService,
	pool *pgxpool.Pool,
) LeaveService {
	return &leaveService{
		leaveTypeRepo:    leaveTypeRepo,
		leaveBalanceRepo: leaveBalanceRepo,
		leaveRequestRepo: leaveRequestRepo,
		employeeRepo:     employeeRepo,
		userRepo:         userRepo,
		attendanceRepo:   attendanceRepo,
		holidayRepo:      holidayRepo,
		auditService:     auditService,
		pool:             pool,
	}
}

func (s *leaveService) CreateLeaveRequest(ctx context.Context, employeeID string, req *dto.CreateLeaveRequestRequest) (*dto.LeaveRequestResponse, error) {
	// Parse and validate dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, ErrInvalidDateRange
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, ErrInvalidDateRange
	}

	if startDate.After(endDate) {
		return nil, ErrInvalidDateRange
	}

	// Get employee by userID (employeeID parameter is actually userID from JWT)
	userUUID, err := uuid.Parse(employeeID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Find employee by user_id
	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeerepo.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}

	// Fetch user from MongoDB to get name
	user, err := s.userRepo.FindByID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get leave type
	leaveTypeUUID, err := uuid.Parse(req.LeaveTypeID)
	if err != nil {
		return nil, fmt.Errorf("invalid leave type ID: %w", err)
	}

	leaveType, err := s.leaveTypeRepo.FindByID(ctx, leaveTypeUUID)
	if err != nil {
		return nil, fmt.Errorf("leave type not found: %w", err)
	}

	if !leaveType.IsActive {
		return nil, errors.New("leave type is not active")
	}

	// Fetch holidays in date range
	holidays, _ := s.holidayRepo.FindByDateRange(ctx, startDate, endDate)
	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h.Date.Format("2006-01-02")] = true
	}

	// Calculate total days (working days only, excluding weekends and holidays)
	totalDays := sharedHelper.CountWorkingDaysExcluding(startDate, endDate, holidayMap)

	// Validate: at least 1 working day
	if totalDays == 0 {
		return nil, errors.New("selected dates contain no working days")
	}

	// Check for overlapping leave requests
	overlapping, err := s.leaveRequestRepo.FindByEmployeeAndDateRange(ctx, employee.ID, startDate, endDate)
	if err == nil && len(overlapping) > 0 {
		return nil, ErrOverlappingLeave
	}

	// Check leave balance
	currentYear := time.Now().Year()
	if startDate.Year() != currentYear {
		return nil, errors.New("cannot request leave for different year")
	}

	balance, err := s.leaveBalanceRepo.FindByEmployeeTypeAndYear(ctx, employee.ID, leaveTypeUUID, currentYear)
	if err != nil {
		if errors.Is(err, leaverepo.ErrLeaveBalanceNotFound) {
			// Initialize balance if not exists
			newBalance := &entity.LeaveBalance{
				ID:          uuid.New().String(),
				EmployeeID:  employee.ID.String(),
				LeaveTypeID: leaveTypeUUID.String(),
				Year:        currentYear,
				Balance:     leaveType.MaxDaysPerYear,
				Used:        0,
				Pending:     0,
			}
			if err := s.leaveBalanceRepo.Create(ctx, newBalance); err != nil {
				return nil, fmt.Errorf("failed to initialize balance: %w", err)
			}
			balance = newBalance
		} else {
			return nil, fmt.Errorf("failed to check balance: %w", err)
		}
	}

	available := balance.Balance - balance.Used - balance.Pending
	if available < totalDays {
		return nil, ErrInsufficientBalance
	}

	// Add to pending
	if err := s.leaveBalanceRepo.AddToPending(ctx, employee.ID, leaveTypeUUID, currentYear, totalDays); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Create leave request
	leaveRequest := &entity.LeaveRequest{
		ID:               uuid.New().String(),
		EmployeeID:       employee.ID.String(),
		LeaveTypeID:      leaveTypeUUID.String(),
		StartDate:        startDate,
		EndDate:          endDate,
		TotalDays:        totalDays,
		Reason:           req.Reason,
		AttachmentURL:    stringPtr(req.AttachmentURL),
		EmergencyContact: stringPtr(req.EmergencyContact),
		Status:           "PENDING",
		CompanyID:        employee.CompanyID.String(),
	}

	if err := s.leaveRequestRepo.Create(ctx, leaveRequest); err != nil {
		return nil, fmt.Errorf("failed to create leave request: %w", err)
	}

	return s.toLeaveRequestResponse(leaveRequest, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid), nil
}

func (s *leaveService) GetMyLeaveRequests(ctx context.Context, userID string, page, perPage int) ([]dto.LeaveRequestResponse, int64, error) {
	// Find employee by userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeerepo.ErrEmployeeNotFound) {
			return nil, 0, ErrEmployeeNotFound
		}
		return nil, 0, fmt.Errorf("employee not found: %w", err)
	}

	offset := (page - 1) * perPage
	requests, err := s.leaveRequestRepo.FindByEmployeeID(ctx, employee.ID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	// PERBAIKAN: Gunakan count query yang benar
	total, err := s.leaveRequestRepo.CountByEmployeeID(ctx, employee.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count leave requests: %w", err)
	}

	// Get user from MongoDB
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("user not found: %w", err)
	}

	responses := make([]dto.LeaveRequestResponse, len(requests))
	for i, req := range requests {
		leaveType, _ := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(req.LeaveTypeID))
		responses[i] = *s.toLeaveRequestResponse(&req, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid)
	}

	return responses, total, nil
}

func (s *leaveService) GetAllLeaveRequests(ctx context.Context, page, perPage int, companyID string) ([]dto.LeaveRequestResponse, int64, error) {
	offset := (page - 1) * perPage

	// Fetch leave requests for the company
	requests, err := s.leaveRequestRepo.FindAllByCompany(ctx, companyID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch leave requests: %w", err)
	}

	// Count total for the company
	total, err := s.leaveRequestRepo.CountAllByCompany(ctx, companyID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count leave requests: %w", err)
	}

	if len(requests) == 0 {
		return []dto.LeaveRequestResponse{}, total, nil
	}

	// Collect unique IDs for batch fetching
	employeeIDs := make(map[uuid.UUID]bool)
	leaveTypeIDs := make(map[uuid.UUID]bool)
	for _, req := range requests {
		if empID, err := uuid.Parse(req.EmployeeID); err == nil {
			employeeIDs[empID] = true
		}
		if ltID, err := uuid.Parse(req.LeaveTypeID); err == nil {
			leaveTypeIDs[ltID] = true
		}
	}

	// Batch fetch employees FIRST to get UserIDs
	empIDSlice := make([]uuid.UUID, 0, len(employeeIDs))
	for id := range employeeIDs {
		empIDSlice = append(empIDSlice, id)
	}
	employees, empErr := s.employeeRepo.FindByIDs(ctx, empIDSlice)
	if empErr != nil {
		return nil, 0, fmt.Errorf("failed to fetch employees: %w", empErr)
	}
	if len(employees) == 0 {
		return []dto.LeaveRequestResponse{}, total, nil
	}

	empMap := make(map[string]*employeerepo.Employee)
	// Build mapping: EmployeeID -> UserID untuk lookup user
	empToUserMap := make(map[string]string)
	userIDs := make(map[string]bool)
	for _, e := range employees {
		empMap[e.ID.String()] = e
		empToUserMap[e.ID.String()] = e.UserID.String()
		userIDs[e.UserID.String()] = true
	}

	// Batch fetch users by UserID (not EmployeeID)
	userIDSlice := make([]string, 0, len(userIDs))
	for id := range userIDs {
		userIDSlice = append(userIDSlice, id)
	}
	users, userErr := s.userRepo.FindByIDs(ctx, userIDSlice)
	if userErr != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", userErr)
	}
	userMap := make(map[string]*userEntity.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// Batch fetch leave types
	ltIDSlice := make([]uuid.UUID, 0, len(leaveTypeIDs))
	for id := range leaveTypeIDs {
		ltIDSlice = append(ltIDSlice, id)
	}
	leaveTypes, ltErr := s.leaveTypeRepo.FindByIDs(ctx, ltIDSlice)
	if ltErr != nil {
		return nil, 0, fmt.Errorf("failed to fetch leave types: %w", ltErr)
	}
	ltMap := make(map[string]*entity.LeaveType)
	for _, lt := range leaveTypes {
		ltMap[lt.ID] = lt
	}

	// Build responses
	responses := make([]dto.LeaveRequestResponse, 0, len(requests))
	for _, req := range requests {
		emp, empOk := empMap[req.EmployeeID]
		lt, ltOk := ltMap[req.LeaveTypeID]

		if !empOk || !ltOk {
			continue
		}

		// Get UserID dari employee, lalu lookup user
		userID := empToUserMap[req.EmployeeID]
		user, userOk := userMap[userID]

		if !userOk {
			continue
		}

		responses = append(responses, *s.toLeaveRequestResponse(
			&req, user.Name, user.Email, emp.Position, lt.Name, lt.Code, lt.IsPaid,
		))
	}

	return responses, total, nil
}

func (s *leaveService) GetLeaveRequestByID(ctx context.Context, id string, companyID string) (*dto.LeaveRequestResponse, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	request, err := s.leaveRequestRepo.FindByIDAndCompany(ctx, requestUUID, companyID)
	if err != nil {
		return nil, ErrLeaveRequestNotFound
	}

	// Get user from MongoDB
	user, err := s.userRepo.FindByID(ctx, request.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get employee for position
	employee, _ := s.employeeRepo.FindByID(ctx, uuid.MustParse(request.EmployeeID))

	leaveType, _ := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(request.LeaveTypeID))

	return s.toLeaveRequestResponse(request, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid), nil
}

func (s *leaveService) GetPendingLeaveRequests(ctx context.Context, companyID string) ([]dto.LeaveRequestResponse, error) {
	requests, err := s.leaveRequestRepo.FindPendingByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return []dto.LeaveRequestResponse{}, nil
	}

	// Collect unique IDs
	userIDs := make(map[string]bool)
	employeeIDs := make(map[uuid.UUID]bool)
	leaveTypeIDs := make(map[uuid.UUID]bool)
	for _, req := range requests {
		userIDs[req.EmployeeID] = true
		if empID, err := uuid.Parse(req.EmployeeID); err == nil {
			employeeIDs[empID] = true
		}
		if ltID, err := uuid.Parse(req.LeaveTypeID); err == nil {
			leaveTypeIDs[ltID] = true
		}
	}

	// Batch fetch users from MongoDB
	userIDSlice := make([]string, 0, len(userIDs))
	for id := range userIDs {
		userIDSlice = append(userIDSlice, id)
	}
	users, _ := s.userRepo.FindByIDs(ctx, userIDSlice)
	userMap := make(map[string]*userEntity.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	// Batch fetch employees
	empIDSlice := make([]uuid.UUID, 0, len(employeeIDs))
	for id := range employeeIDs {
		empIDSlice = append(empIDSlice, id)
	}
	employees, _ := s.employeeRepo.FindByIDs(ctx, empIDSlice)
	empMap := make(map[string]*employeerepo.Employee)
	for _, e := range employees {
		empMap[e.ID.String()] = e
	}

	// Batch fetch leave types
	ltIDSlice := make([]uuid.UUID, 0, len(leaveTypeIDs))
	for id := range leaveTypeIDs {
		ltIDSlice = append(ltIDSlice, id)
	}
	leaveTypes, _ := s.leaveTypeRepo.FindByIDs(ctx, ltIDSlice)
	ltMap := make(map[string]*entity.LeaveType)
	for _, lt := range leaveTypes {
		ltMap[lt.ID] = lt
	}

	// Build responses
	responses := make([]dto.LeaveRequestResponse, 0, len(requests))
	for _, req := range requests {
		user, userOk := userMap[req.EmployeeID]
		emp, empOk := empMap[req.EmployeeID]
		lt, ltOk := ltMap[req.LeaveTypeID]

		if !userOk || !empOk || !ltOk {
			continue
		}

		responses = append(responses, *s.toLeaveRequestResponse(
			&req, user.Name, user.Email, emp.Position,
			lt.Name, lt.Code, lt.IsPaid,
		))
	}

	return responses, nil
}

func (s *leaveService) ApproveLeaveRequest(ctx context.Context, id, approverID string, companyID string, req *dto.ApproveLeaveRequest) (*dto.LeaveRequestResponse, error) {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid request ID: %w", err)
	}

	// Find approver's employee record (approverID is UserID from JWT)
	// If not found (e.g., admin without employee record), use NULL
	approverUserUUID, _ := uuid.Parse(approverID)
	approverEmployee, _ := s.employeeRepo.FindByUserID(ctx, approverUserUUID)

	// Use pointer for approverUUID (can be nil for admin without employee record)
	var approverUUID *uuid.UUID
	if approverEmployee != nil {
		approverUUID = &approverEmployee.ID
	}

	request, err := s.leaveRequestRepo.FindByIDAndCompany(ctx, requestUUID, companyID)
	if err != nil {
		return nil, ErrLeaveRequestNotFound
	}

	if request.Status != "PENDING" {
		return nil, ErrInvalidRequestStatus
	}

	// Get leave type
	leaveType, err := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(request.LeaveTypeID))
	if err != nil {
		return nil, fmt.Errorf("leave type not found: %w", err)
	}

	// === BEGIN TRANSACTION ===
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Move balance from pending to used within transaction
	if leaveType.IsPaid {
		currentYear := request.StartDate.Year()
		employeeUUID := uuid.MustParse(request.EmployeeID)
		leaveTypeUUID := uuid.MustParse(request.LeaveTypeID)

		_, err = tx.Exec(ctx,
			`UPDATE leave_balances
			 SET pending = pending - $1, used = used + $1, updated_at = NOW()
			 WHERE employee_id = $2 AND leave_type_id = $3 AND year = $4`,
			request.TotalDays, employeeUUID, leaveTypeUUID, currentYear,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update balance: %w", err)
		}
	}

	// Update request status using transaction
	now := time.Now()
	_, err = tx.Exec(ctx,
		`UPDATE leave_requests
		 SET status = $1, approved_by = $2, approved_at = $3, updated_at = NOW()
		 WHERE id = $4`,
		"APPROVED", approverUUID, now, requestUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update request: %w", err)
	}

	// === COMMIT TRANSACTION ===
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Create attendance records for leave days (non-critical, outside transaction)
	s.createLeaveAttendances(ctx, uuid.MustParse(request.EmployeeID), request.StartDate, request.EndDate, requestUUID)

	// Audit log — leave request approved
	_ = s.auditService.Log(ctx, auditservice.AuditEntry{
		UserID:       approverID,
		Action:       "APPROVE",
		ResourceType: "leave_request",
		ResourceID:   id,
		NewData: map[string]interface{}{
			"status": "APPROVED",
			"notes":  req.ApprovalNote,
		},
	})

	employee, err := s.employeeRepo.FindByID(ctx, uuid.MustParse(request.EmployeeID))
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, employee.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	updatedRequest, err := s.leaveRequestRepo.FindByID(ctx, requestUUID)
	if err != nil {
		return nil, err
	}

	return s.toLeaveRequestResponse(updatedRequest, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid), nil
}

func (s *leaveService) RejectLeaveRequest(ctx context.Context, id, approverID string, companyID string, req *dto.RejectLeaveRequest) error {
	requestUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid request ID: %w", err)
	}

	request, err := s.leaveRequestRepo.FindByIDAndCompany(ctx, requestUUID, companyID)
	if err != nil {
		return ErrLeaveRequestNotFound
	}

	if request.Status != "PENDING" {
		return ErrInvalidRequestStatus
	}

	// Get leave type
	leaveType, err := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(request.LeaveTypeID))
	if err != nil {
		return fmt.Errorf("leave type not found: %w", err)
	}

	// === BEGIN TRANSACTION ===
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Return balance from pending within transaction
	if leaveType.IsPaid {
		currentYear := request.StartDate.Year()
		employeeUUID := uuid.MustParse(request.EmployeeID)
		leaveTypeUUID := uuid.MustParse(request.LeaveTypeID)

		_, err = tx.Exec(ctx,
			`UPDATE leave_balances
			 SET pending = pending - $1, updated_at = NOW()
			 WHERE employee_id = $2 AND leave_type_id = $3 AND year = $4`,
			request.TotalDays, employeeUUID, leaveTypeUUID, currentYear,
		)
		if err != nil {
			return fmt.Errorf("failed to return balance: %w", err)
		}
	}

	// Find approver's employee record (approverID is UserID from JWT)
	// If not found (e.g., admin without employee record), use NULL
	approverUserUUID := uuid.MustParse(approverID)
	approverEmployee, empErr := s.employeeRepo.FindByUserID(ctx, approverUserUUID)

	var approverUUID *uuid.UUID
	if empErr == nil {
		approverUUID = &approverEmployee.ID
	}

	_, err = tx.Exec(ctx,
		`UPDATE leave_requests
		 SET status = $1, approved_by = $2, rejection_reason = $3, updated_at = NOW()
		 WHERE id = $4`,
		"REJECTED", approverUUID, req.RejectionReason, requestUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update request: %w", err)
	}

	// === COMMIT TRANSACTION ===
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Audit log — leave request rejected
	_ = s.auditService.Log(ctx, auditservice.AuditEntry{
		UserID:       approverID,
		Action:       "REJECT",
		ResourceType: "leave_request",
		ResourceID:   id,
		NewData: map[string]interface{}{
			"status": "REJECTED",
			"reason": req.RejectionReason,
		},
	})

	return nil
}

func (s *leaveService) GetMyLeaveBalances(ctx context.Context, userID string, year int) (*dto.LeaveBalanceResponse, error) {
	// Find employee by userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeerepo.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}

	// Get user from MongoDB
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get balances
	balances, err := s.leaveBalanceRepo.FindByEmployeeAndYear(ctx, employee.ID, year)
	if err != nil {
		return nil, err
	}

	response := &dto.LeaveBalanceResponse{
		EmployeeID:   employee.ID.String(),
		EmployeeName: user.Name,
		Year:         year,
		Balances:     make([]dto.LeaveBalanceItem, 0),
	}

	for _, balance := range balances {
		leaveType, _ := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(balance.LeaveTypeID))

		response.Balances = append(response.Balances, dto.LeaveBalanceItem{
			LeaveTypeID:   balance.LeaveTypeID,
			LeaveTypeName: leaveType.Name,
			Balance:       balance.Balance,
			Used:          balance.Used,
			Pending:       balance.Pending,
			Available:     balance.Balance - balance.Used - balance.Pending,
		})
	}

	return response, nil
}

func (s *leaveService) createLeaveAttendances(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time, leaveRequestID uuid.UUID) {
	currentDate := startDate
	for !currentDate.After(endDate) {
		// Skip weekends
		weekday := currentDate.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		attendance := &leaveattendance.Attendance{
			ID:         uuid.New().String(),
			EmployeeID: employeeID.String(),
			Date:       currentDate,
			Status:     "LEAVE",
			Notes:      fmt.Sprintf("Leave Request ID: %s", leaveRequestID.String()),
		}

		// Log error if attendance creation fails instead of silently ignoring
		if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
			zap.L().Error("failed to create leave attendance",
				zap.String("employeeID", employeeID.String()),
				zap.Time("date", currentDate),
				zap.String("leaveRequestID", leaveRequestID.String()),
				zap.Error(err),
			)
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}
}

func (s *leaveService) toLeaveRequestResponse(request *entity.LeaveRequest, userName, userEmail, employeePosition, leaveTypeName, leaveTypeCode string, leaveTypeIsPaid bool) *dto.LeaveRequestResponse {
	var approvedAt *string
	if request.ApprovedAt != nil {
		formatted := request.ApprovedAt.Format(time.RFC3339)
		approvedAt = &formatted
	}

	return &dto.LeaveRequestResponse{
		ID:               request.ID,
		EmployeeID:       request.EmployeeID,
		EmployeeName:     userName,
		EmployeePosition: employeePosition,
		LeaveType: dto.LeaveTypeDetail{
			ID:     request.LeaveTypeID,
			Name:   leaveTypeName,
			Code:   leaveTypeCode,
			IsPaid: leaveTypeIsPaid,
		},
		StartDate:        request.StartDate.Format("2006-01-02"),
		EndDate:          request.EndDate.Format("2006-01-02"),
		TotalDays:        request.TotalDays,
		Reason:           request.Reason,
		AttachmentURL:    stringValue(request.AttachmentURL),
		EmergencyContact: stringValue(request.EmergencyContact),
		Status:           request.Status,
		ApprovedBy:       request.ApprovedBy,
		ApprovedAt:       approvedAt,
		RejectionReason:  stringValue(request.RejectionReason),
		CreatedAt:        request.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        request.UpdatedAt.Format(time.RFC3339),
	}
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function to convert *string to string
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
