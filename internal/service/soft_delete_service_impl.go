package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"hris/internal/employee/repository"
)

type SoftDeleteService interface {
	SoftDeleteEmployee(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
}

type softDeleteServiceImpl struct {
	employeeRepo repository.EmployeeRepository
}

func NewSoftDeleteService(employeeRepo repository.EmployeeRepository) SoftDeleteService {
	return &softDeleteServiceImpl{
		employeeRepo: employeeRepo,
	}
}

func (s *softDeleteServiceImpl) SoftDeleteEmployee(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	employee, err := s.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Convert EmployeeWithUser to Employee for update
	employeeUpdate := &repository.Employee{
		ID:                 employee.ID,
		UserID:             employee.UserID,
		FullName:           employee.FullName,
		Position:           employee.Position,
		PhoneNumber:        employee.PhoneNumber,
		Address:            employee.Address,
		SalaryBase:         employee.SalaryBase,
		JoinDate:           employee.JoinDate,
		BankName:           employee.BankName,
		BankAccountNumber:  employee.BankAccountNumber,
		BankAccountHolder:  employee.BankAccountHolder,
		ScheduleID:         employee.ScheduleID,
		EmploymentStatus:   employee.EmploymentStatus,
		JobLevel:           employee.JobLevel,
		Gender:             employee.Gender,
		Division:           employee.Division,
		DepartmentID:       employee.DepartmentID,
		CreatedAt:          employee.CreatedAt,
		UpdatedAt:          time.Now(),
	}

	// Apply updates
	for field, value := range updates {
		switch field {
		case "full_name":
			employeeUpdate.FullName = value.(string)
		case "position":
			employeeUpdate.Position = value.(string)
		case "phone_number":
			employeeUpdate.PhoneNumber = value.(string)
		case "address":
			employeeUpdate.Address = value.(string)
		case "salary_base":
			employeeUpdate.SalaryBase = value.(float64)
		case "join_date":
			employeeUpdate.JoinDate = value.(time.Time)
		case "bank_name":
			employeeUpdate.BankName = value.(string)
		case "bank_account_number":
			employeeUpdate.BankAccountNumber = value.(string)
		case "bank_account_holder":
			employeeUpdate.BankAccountHolder = value.(string)
		case "schedule_id":
			if value != nil {
				employeeUpdate.ScheduleID = value.(*uuid.UUID)
			}
		case "department_id":
			if value != nil {
				employeeUpdate.DepartmentID = value.(*uuid.UUID)
			}
		case "employment_status":
			employeeUpdate.EmploymentStatus = value.(string)
		case "job_level":
			employeeUpdate.JobLevel = value.(string)
		case "gender":
			employeeUpdate.Gender = value.(string)
		case "division":
			employeeUpdate.Division = value.(string)
		}
	}

	return s.employeeRepo.Update(ctx, employeeUpdate)
}
