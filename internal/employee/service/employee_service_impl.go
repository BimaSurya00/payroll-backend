package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hris/internal/employee/dto"
	"hris/internal/employee/helper"
	"hris/internal/employee/repository"
	userEntity "hris/internal/user/entity"
	userRepository "hris/internal/user/repository"
	sharedHelper "hris/shared/helper"
)

var (
	ErrEmployeeNotFound     = errors.New("employee not found")
	ErrUserCreationFailed   = errors.New("failed to create user")
	ErrEmployeeCreateFailed = errors.New("failed to create employee profile")
)

type employeeService struct {
	employeeRepo repository.EmployeeRepository
	userRepo     userRepository.UserRepository
}

func NewEmployeeService(employeeRepo repository.EmployeeRepository, userRepo userRepository.UserRepository) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
	}
}

func (s *employeeService) CreateEmployee(ctx context.Context, req *dto.CreateEmployeeRequest, companyID string) (*dto.EmployeeResponse, error) {
	// 1. Validate Request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 2. Create User with company_id
	hashedPassword, err := sharedHelper.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := userEntity.NewUser(req.Name, req.Email, hashedPassword, "USER", companyID)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUserCreationFailed, err)
	}

	// 3. Create Employee (Postgres)
	joinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		// Compensating Action: Rollback User creation
		_ = s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("invalid join date format: %w", err)
	}

	var scheduleID *uuid.UUID
	if req.ScheduleID != "" {
		parsedScheduleID, err := uuid.Parse(req.ScheduleID)
		if err != nil {
			// Compensating Action: Rollback User creation
			_ = s.userRepo.Delete(ctx, user.ID)
			return nil, fmt.Errorf("invalid schedule ID format: %w", err)
		}
		scheduleID = &parsedScheduleID
	}

	var departmentID *uuid.UUID
	if req.DepartmentID != "" {
		parsedDeptID, err := uuid.Parse(req.DepartmentID)
		if err != nil {
			// Compensating Action: Rollback User creation
			_ = s.userRepo.Delete(ctx, user.ID)
			return nil, fmt.Errorf("invalid department ID format: %w", err)
		}
		departmentID = &parsedDeptID
	}

	userID, err := uuid.Parse(user.ID)
	if err != nil {
		// Compensating Action: Rollback User creation
		_ = s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	companyUUID, err := uuid.Parse(companyID)
	if err != nil {
		// Compensating Action: Rollback User creation
		_ = s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("invalid company ID format: %w", err)
	}

	employee := &repository.Employee{
		ID:                uuid.New(),
		CompanyID:         companyUUID,
		UserID:            userID,
		Position:          req.Position,
		PhoneNumber:       req.PhoneNumber,
		Address:           req.Address,
		SalaryBase:        req.SalaryBase,
		JoinDate:          joinDate,
		BankName:          req.BankName,
		BankAccountNumber: req.BankAccountNumber,
		BankAccountHolder: req.BankAccountHolder,
		ScheduleID:        scheduleID,
		DepartmentID:      departmentID,
		EmploymentStatus:  req.EmploymentStatus,
		JobLevel:          req.JobLevel,
		Gender:            req.Gender,
		Division:          req.Division,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		// Compensating Action: Rollback User creation
		_ = s.userRepo.Delete(ctx, user.ID)
		return nil, fmt.Errorf("%w: %v", ErrEmployeeCreateFailed, err)
	}

	// 4. Return Success - fetch created employee data
	employeeData, err := s.employeeRepo.FindByID(ctx, employee.ID)
	if err != nil {
		// Employee was created but we can't fetch it - return basic response
		return helper.ToEmployeeResponse(employee, user), nil
	}

	// Use ToEmployeeResponseFromDB to properly convert EmployeeWithUser
	// This includes DepartmentName from the JOIN query
	return helper.ToEmployeeResponseFromDB(employeeData), nil
}

func (s *employeeService) GetAllEmployees(ctx context.Context, page, perPage int, path string, search string, companyID string) (*helper.Pagination[*dto.EmployeeResponse], error) {
	employees, total, err := s.employeeRepo.FindAllByCompany(ctx, companyID, page, perPage, search)
	if err != nil {
		return nil, err
	}

	// Fetch user data for each employee
	employeeResponses := make([]*dto.EmployeeResponse, 0, len(employees))
	for _, emp := range employees {
		// Get user using the user_id
		user, err := s.userRepo.FindByID(ctx, emp.UserID.String())
		if err != nil {
			// If user not found, still include employee but with empty user data
			employeeWithoutUser := &repository.Employee{
				ID:                emp.ID,
				UserID:            emp.UserID,
				Position:          emp.Position,
				PhoneNumber:       emp.PhoneNumber,
				Address:           emp.Address,
				SalaryBase:        emp.SalaryBase,
				JoinDate:          emp.JoinDate,
				BankName:          emp.BankName,
				BankAccountNumber: emp.BankAccountNumber,
				BankAccountHolder: emp.BankAccountHolder,
				ScheduleID:        emp.ScheduleID,
				EmploymentStatus:  emp.EmploymentStatus,
				JobLevel:          emp.JobLevel,
				Gender:            emp.Gender,
				Division:          emp.Division,
				CreatedAt:         emp.CreatedAt,
				UpdatedAt:         emp.UpdatedAt,
			}
			employeeResponses = append(employeeResponses, helper.ToEmployeeResponse(employeeWithoutUser, nil))
		} else {
			employeeWithUser := &repository.Employee{
				ID:                emp.ID,
				UserID:            emp.UserID,
				Position:          emp.Position,
				PhoneNumber:       emp.PhoneNumber,
				Address:           emp.Address,
				SalaryBase:        emp.SalaryBase,
				JoinDate:          emp.JoinDate,
				BankName:          emp.BankName,
				BankAccountNumber: emp.BankAccountNumber,
				BankAccountHolder: emp.BankAccountHolder,
				ScheduleID:        emp.ScheduleID,
				EmploymentStatus:  emp.EmploymentStatus,
				JobLevel:          emp.JobLevel,
				Gender:            emp.Gender,
				Division:          emp.Division,
				CreatedAt:         emp.CreatedAt,
				UpdatedAt:         emp.UpdatedAt,
			}
			employeeResponses = append(employeeResponses, helper.ToEmployeeResponse(employeeWithUser, user))
		}
	}

	pagination := helper.NewPagination(employeeResponses, page, perPage, total, path)

	return pagination, nil
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id string, companyID string) (*dto.EmployeeResponse, error) {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID format: %w", err)
	}

	employeeData, err := s.employeeRepo.FindByIDAndCompany(ctx, employeeID, companyID)
	if err != nil {
		if err.Error() == "employee not found" {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Fetch user data
	user, err := s.userRepo.FindByID(ctx, employeeData.UserID.String())
	if err != nil {
		// Return employee without user data if user not found (but with schedule)
		return helper.ToEmployeeResponseWithSchedule(employeeData, nil), nil
	}

	// Return employee with user data and schedule detail
	return helper.ToEmployeeResponseWithSchedule(employeeData, user), nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, id string, companyID string, req *dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error) {
	if !req.HasUpdates() {
		return nil, fmt.Errorf("no fields to update")
	}

	employeeID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID format: %w", err)
	}

	// Get existing employee with company verification
	existing, err := s.employeeRepo.FindByIDAndCompany(ctx, employeeID, companyID)
	if err != nil {
		if err.Error() == "employee not found" {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Build update entity
	employee := &repository.Employee{
		ID:                employeeID,
		UserID:            existing.UserID,
		Position:          existing.Position,
		PhoneNumber:       existing.PhoneNumber,
		Address:           existing.Address,
		SalaryBase:        existing.SalaryBase,
		JoinDate:          existing.JoinDate,
		BankName:          existing.BankName,
		BankAccountNumber: existing.BankAccountNumber,
		BankAccountHolder: existing.BankAccountHolder,
		ScheduleID:        existing.ScheduleID,
		CreatedAt:         existing.CreatedAt,
		UpdatedAt:         time.Now(),
	}

	// Set defaults for new fields if they're empty
	if existing.EmploymentStatus == "" {
		employee.EmploymentStatus = "PROBATION"
	} else {
		employee.EmploymentStatus = existing.EmploymentStatus
	}

	if existing.JobLevel == "" {
		employee.JobLevel = "STAFF"
	} else {
		employee.JobLevel = existing.JobLevel
	}

	employee.Gender = existing.Gender
	if existing.Division == "" {
		employee.Division = "GENERAL"
	} else {
		employee.Division = existing.Division
	}

	// Update fields if provided
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.PhoneNumber != nil {
		employee.PhoneNumber = *req.PhoneNumber
	}
	if req.Address != nil {
		employee.Address = *req.Address
	}
	if req.SalaryBase != nil {
		employee.SalaryBase = *req.SalaryBase
	}
	if req.BankName != nil {
		employee.BankName = *req.BankName
	}
	if req.BankAccountNumber != nil {
		employee.BankAccountNumber = *req.BankAccountNumber
	}
	if req.BankAccountHolder != nil {
		employee.BankAccountHolder = *req.BankAccountHolder
	}

	// Parse schedule ID if provided
	if req.ScheduleID != nil {
		parsedScheduleID, err := uuid.Parse(*req.ScheduleID)
		if err != nil {
			return nil, fmt.Errorf("invalid schedule ID format: %w", err)
		}
		employee.ScheduleID = &parsedScheduleID
	}

	// Update new fields if provided
	if req.EmploymentStatus != nil {
		employee.EmploymentStatus = *req.EmploymentStatus
	}
	if req.JobLevel != nil {
		employee.JobLevel = *req.JobLevel
	}
	if req.Gender != nil {
		employee.Gender = *req.Gender
	}
	if req.Division != nil {
		employee.Division = *req.Division
	}

	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, err
	}

	// Fetch updated employee and user data
	updatedData, err := s.employeeRepo.FindByID(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	// Fetch user data
	user, err := s.userRepo.FindByID(ctx, updatedData.UserID.String())
	if err != nil {
		// Return employee without user data if user not found
		updatedEmployee := &repository.Employee{
			ID:                updatedData.ID,
			UserID:            updatedData.UserID,
			Position:          updatedData.Position,
			PhoneNumber:       updatedData.PhoneNumber,
			Address:           updatedData.Address,
			SalaryBase:        updatedData.SalaryBase,
			JoinDate:          updatedData.JoinDate,
			BankName:          updatedData.BankName,
			BankAccountNumber: updatedData.BankAccountNumber,
			BankAccountHolder: updatedData.BankAccountHolder,
			ScheduleID:        updatedData.ScheduleID,
			EmploymentStatus:  updatedData.EmploymentStatus,
			JobLevel:          updatedData.JobLevel,
			Gender:            updatedData.Gender,
			Division:          updatedData.Division,
			CreatedAt:         updatedData.CreatedAt,
			UpdatedAt:         updatedData.UpdatedAt,
		}
		return helper.ToEmployeeResponse(updatedEmployee, nil), nil
	}

	updatedEmployee := &repository.Employee{
		ID:                updatedData.ID,
		UserID:            updatedData.UserID,
		Position:          updatedData.Position,
		PhoneNumber:       updatedData.PhoneNumber,
		Address:           updatedData.Address,
		SalaryBase:        updatedData.SalaryBase,
		JoinDate:          updatedData.JoinDate,
		BankName:          updatedData.BankName,
		BankAccountNumber: updatedData.BankAccountNumber,
		BankAccountHolder: updatedData.BankAccountHolder,
		ScheduleID:        updatedData.ScheduleID,
		EmploymentStatus:  updatedData.EmploymentStatus,
		JobLevel:          updatedData.JobLevel,
		Gender:            updatedData.Gender,
		Division:          updatedData.Division,
		CreatedAt:         updatedData.CreatedAt,
		UpdatedAt:         updatedData.UpdatedAt,
	}
	return helper.ToEmployeeResponse(updatedEmployee, user), nil
}

func (s *employeeService) DeleteEmployee(ctx context.Context, id string, companyID string) error {
	employeeID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee ID format: %w", err)
	}

	// Get employee first to get user ID and verify company ownership
	employee, err := s.employeeRepo.FindByIDAndCompany(ctx, employeeID, companyID)
	if err != nil {
		if err.Error() == "employee not found" {
			return ErrEmployeeNotFound
		}
		return fmt.Errorf("failed to get employee: %w", err)
	}

	// Delete employee (Postgres) - use DeleteByCompany for extra safety
	if err := s.employeeRepo.DeleteByCompany(ctx, employeeID, companyID); err != nil {
		return err
	}

	// Delete associated user
	if err := s.userRepo.Delete(ctx, employee.UserID.String()); err != nil {
		// Log error but don't fail - employee is already deleted
		// In production, you might want to implement a cleanup job
		return fmt.Errorf("employee deleted but failed to delete associated user: %w", err)
	}

	return nil
}

func (s *employeeService) GetMyProfile(ctx context.Context, userID string) (*dto.EmployeeResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	employeeWithUser, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}

	// Get user data
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return helper.ToEmployeeResponseWithSchedule(employeeWithUser, user), nil
}

func (s *employeeService) UpdateMyProfile(ctx context.Context, userID string, req *dto.UpdateMyProfileRequest) (*dto.EmployeeResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Fetch current employee
	employeeWithUser, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}

	// Build employee struct for update
	employee := &repository.Employee{
		ID:                employeeWithUser.ID,
		UserID:            employeeWithUser.UserID,
		Position:          employeeWithUser.Position,
		PhoneNumber:       employeeWithUser.PhoneNumber,
		SalaryBase:        employeeWithUser.SalaryBase,
		Address:           employeeWithUser.Address,
		BankName:          employeeWithUser.BankName,
		BankAccountNumber: employeeWithUser.BankAccountNumber,
		BankAccountHolder: employeeWithUser.BankAccountHolder,
		JoinDate:          employeeWithUser.JoinDate,
		EmploymentStatus:  employeeWithUser.EmploymentStatus,
		JobLevel:          employeeWithUser.JobLevel,
		Gender:            employeeWithUser.Gender,
		Division:          employeeWithUser.Division,
		ScheduleID:        employeeWithUser.ScheduleID,
		CreatedAt:         employeeWithUser.CreatedAt,
		UpdatedAt:         time.Now(),
	}

	// Update only allowed fields
	hasChanges := false
	if req.PhoneNumber != nil {
		employee.PhoneNumber = *req.PhoneNumber
		hasChanges = true
	}
	if req.Address != nil {
		employee.Address = *req.Address
		hasChanges = true
	}
	if req.BankName != nil {
		employee.BankName = *req.BankName
		hasChanges = true
	}
	if req.BankAccountNumber != nil {
		employee.BankAccountNumber = *req.BankAccountNumber
		hasChanges = true
	}
	if req.BankAccountHolder != nil {
		employee.BankAccountHolder = *req.BankAccountHolder
		hasChanges = true
	}

	if !hasChanges {
		return nil, errors.New("no fields to update")
	}

	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Return updated profile
	return s.GetMyProfile(ctx, userID)
}

func (s *employeeService) validateCreateRequest(req *dto.CreateEmployeeRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.Position == "" {
		return errors.New("position is required")
	}
	if req.SalaryBase <= 0 {
		return errors.New("salary base must be greater than 0")
	}
	if req.JoinDate == "" {
		return errors.New("join date is required")
	}
	return nil
}
