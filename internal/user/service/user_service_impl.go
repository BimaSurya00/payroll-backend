package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	auditservice "hris/internal/audit/service"
	employeeRepo "hris/internal/employee/repository"
	scheduleRepo "hris/internal/schedule/repository"
	"hris/internal/user/dto"
	"hris/internal/user/entity"
	"hris/internal/user/helper"
	userRepository "hris/internal/user/repository"
	sharedHelper "hris/shared/helper"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type service struct {
	repo         userRepository.UserRepository
	employeeRepo employeeRepo.EmployeeRepository
	scheduleRepo scheduleRepo.ScheduleRepository
	auditService auditservice.AuditService
}

func NewUserService(repo userRepository.UserRepository) UserService {
	return &service{
		repo:         repo,
		employeeRepo: nil,
		scheduleRepo: nil,
	}
}

func NewUserServiceWithEmployee(
	repo userRepository.UserRepository,
	employeeRepo employeeRepo.EmployeeRepository,
	scheduleRepo scheduleRepo.ScheduleRepository,
) UserService {
	return &service{
		repo:         repo,
		employeeRepo: employeeRepo,
		scheduleRepo: scheduleRepo,
	}
}

func NewUserServiceWithAudit(repo userRepository.UserRepository, auditService auditservice.AuditService) UserService {
	return &service{
		repo:         repo,
		employeeRepo: nil,
		scheduleRepo: nil,
		auditService: auditService,
	}
}

func NewUserServiceWithEmployeeAndAudit(
	repo userRepository.UserRepository,
	employeeRepo employeeRepo.EmployeeRepository,
	scheduleRepo scheduleRepo.ScheduleRepository,
	auditService auditservice.AuditService,
) UserService {
	return &service{
		repo:         repo,
		employeeRepo: employeeRepo,
		scheduleRepo: scheduleRepo,
		auditService: auditService,
	}
}

func (s *service) CreateUser(ctx context.Context, req *dto.CreateUserRequest, companyID string) (*dto.UserResponse, error) {
	// Check if email already exists within the same company
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, userRepository.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing email: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := sharedHelper.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity with company_id from admin context
	user := entity.NewUser(req.Name, req.Email, hashedPassword, req.Role, companyID)

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Audit log — user created
	if s.auditService != nil {
		_ = s.auditService.Log(ctx, auditservice.AuditEntry{
			Action:       "CREATE",
			ResourceType: "user",
			ResourceID:   user.ID,
			NewData: map[string]interface{}{
				"email":     user.Email,
				"name":      user.Name,
				"role":      user.Role,
				"companyId": user.CompanyID,
			},
		})
	}

	// Auto-create employee for USER role if employee repo is available
	if req.Role == "USER" && s.employeeRepo != nil && s.scheduleRepo != nil {
		// Get default schedule (first schedule)
		schedules, err := s.scheduleRepo.FindAll(ctx, 0, 1)
		if err == nil && len(schedules) > 0 {
			defaultSchedule := schedules[0]

			userID, err := uuid.Parse(user.ID)
			if err == nil {
				scheduleID, _ := uuid.Parse(defaultSchedule.ID)
				employee := &employeeRepo.Employee{
					ID:         uuid.New(),
					UserID:     userID,
					Position:   "Staff",
					SalaryBase: 0,
					JoinDate:   time.Now(),
					ScheduleID: &scheduleID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				// Create employee (ignore error, user already created)
				_ = s.employeeRepo.Create(ctx, employee)
			}
		}
	}

	return dto.ToUserResponse(user), nil
}

func (s *service) GetUsers(ctx context.Context, page, perPage int, path string, companyID string, userRole string) (*helper.Pagination[*dto.UserResponse], error) {
	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	var users []*entity.User
	var total int64
	var err error

	// SUPER_USER can see all users, others only see their company
	if userRole == "SUPER_USER" {
		users, err = s.repo.FindAll(ctx, skip, limit)
		if err != nil {
			return nil, err
		}
		total, err = s.repo.Count(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// Filter by company_id for non-super users
		users, err = s.repo.FindAllByCompany(ctx, companyID, skip, limit)
		if err != nil {
			return nil, err
		}
		total, err = s.repo.CountByCompany(ctx, companyID)
		if err != nil {
			return nil, err
		}
	}

	userResponses := dto.ToUserResponses(users)
	pagination := helper.NewPagination(userResponses, page, perPage, total, path)

	return pagination, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, userRepository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *service) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if !req.HasUpdates() {
		return nil, fmt.Errorf("no fields to update")
	}

	// Fetch current user
	oldUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, userRepository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Apply updates to user entity
	user := *oldUser // copy to avoid modifying original

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Email != nil {
		existingUser, err := s.repo.FindByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, userRepository.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to check existing email: %w", err)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := sharedHelper.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.CompanyID != nil {
		user.CompanyID = *req.CompanyID
	}

	// Update user — repo handles updated_at
	if err := s.repo.Update(ctx, &user); err != nil {
		return nil, err
	}

	// Audit log — user updated
	if s.auditService != nil {
		_ = s.auditService.Log(ctx, auditservice.AuditEntry{
			Action:       "UPDATE",
			ResourceType: "user",
			ResourceID:   id,
			OldData: map[string]interface{}{
				"name":  oldUser.Name,
				"email": oldUser.Email,
				"role":  oldUser.Role,
			},
			NewData: map[string]interface{}{
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	}

	// Fetch updated user (to get server-side updated_at)
	updatedUser, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, userRepository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dto.ToUserResponse(updatedUser), nil
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, userRepository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Audit log — user deleted
	if s.auditService != nil {
		_ = s.auditService.Log(ctx, auditservice.AuditEntry{
			Action:       "DELETE",
			ResourceType: "user",
			ResourceID:   id,
		})
	}

	return nil
}
