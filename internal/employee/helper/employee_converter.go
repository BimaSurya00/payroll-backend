package helper

import (
	"time"

	"hris/internal/employee/dto"
	"hris/internal/employee/repository"
	userEntity "hris/internal/user/entity"
)

// ToEmployeeResponse converts Employee entity + User entity to EmployeeResponse DTO
func ToEmployeeResponse(employee *repository.Employee, user *userEntity.User) *dto.EmployeeResponse {
	// Defensive: Check if employee is nil
	if employee == nil {
		return nil
	}

	var scheduleID *string
	if employee.ScheduleID != nil {
		sid := employee.ScheduleID.String()
		scheduleID = &sid
	}

	// Handle empty strings with defaults
	employmentStatus := employee.EmploymentStatus
	if employmentStatus == "" {
		employmentStatus = "PROBATION"
	}

	jobLevel := employee.JobLevel
	if jobLevel == "" {
		jobLevel = "STAFF"
	}

	division := employee.Division
	if division == "" {
		division = "GENERAL"
	}

	// Handle nil user - DEFENSIVE CHECK NEEDED!
	userName := ""
	userEmail := ""
	if user != nil {
		userName = user.Name
		userEmail = user.Email
	}

	// Handle DepartmentID
	var deptID *string
	if employee.DepartmentID != nil {
		did := employee.DepartmentID.String()
		deptID = &did
	}

	return &dto.EmployeeResponse{
		ID:                employee.ID.String(),
		UserID:            employee.UserID.String(),
		UserName:          userName,
		UserEmail:         userEmail,
		Position:          employee.Position,
		PhoneNumber:       employee.PhoneNumber,
		Address:           employee.Address,
		SalaryBase:        employee.SalaryBase,
		JoinDate:          employee.JoinDate.Format("2006-01-02"),
		BankName:          employee.BankName,
		BankAccountNumber: employee.BankAccountNumber,
		BankAccountHolder: employee.BankAccountHolder,
		ScheduleID:        scheduleID,
		Schedule:          nil, // No schedule detail when using Employee entity
		EmploymentStatus:  employmentStatus,
		JobLevel:          jobLevel,
		Gender:            employee.Gender,
		Division:          division,
		DepartmentID:      deptID,
		DepartmentName:    "", // Not available when using Employee entity
		CreatedAt:         employee.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         employee.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEmployeeResponseWithSchedule converts EmployeeWithUser + User to EmployeeResponse with schedule detail
func ToEmployeeResponseWithSchedule(employee *repository.EmployeeWithUser, user *userEntity.User) *dto.EmployeeResponse {
	var scheduleID *string
	var scheduleDetail *dto.ScheduleDetail

	// Defensive: Check if employee is nil
	if employee == nil {
		return nil
	}

	// Only build schedule detail if all required fields are not nil
	if employee.ScheduleID != nil && employee.ScheduleName != nil && employee.ScheduleTimeIn != nil && employee.ScheduleTimeOut != nil {
		sid := employee.ScheduleID.String()
		scheduleID = &sid

		// Include schedule detail if available
		allowedLate := 0
		if employee.ScheduleAllowedLateMinutes != nil {
			allowedLate = *employee.ScheduleAllowedLateMinutes
		}

		scheduleDetail = &dto.ScheduleDetail{
			ID:                 sid,
			Name:               *employee.ScheduleName,
			TimeIn:             *employee.ScheduleTimeIn,
			TimeOut:            *employee.ScheduleTimeOut,
			AllowedLateMinutes: allowedLate,
		}
	}

	userName := ""
	userEmail := ""
	if user != nil {
		userName = user.Name
		userEmail = user.Email
	}

	// Handle empty strings with defaults
	employmentStatus := employee.EmploymentStatus
	if employmentStatus == "" {
		employmentStatus = "PROBATION"
	}

	jobLevel := employee.JobLevel
	if jobLevel == "" {
		jobLevel = "STAFF"
	}

	division := employee.Division
	if division == "" {
		division = "GENERAL"
	}

	// Handle DepartmentID
	var deptID *string
	if employee.DepartmentID != nil {
		did := employee.DepartmentID.String()
		deptID = &did
	}

	// Handle DepartmentName
	deptName := ""
	if employee.DepartmentName != nil {
		deptName = *employee.DepartmentName
	}

	return &dto.EmployeeResponse{
		ID:                employee.ID.String(),
		UserID:            employee.UserID.String(),
		UserName:          userName,
		UserEmail:         userEmail,
		Position:          employee.Position,
		PhoneNumber:       employee.PhoneNumber,
		Address:           employee.Address,
		SalaryBase:        employee.SalaryBase,
		JoinDate:          employee.JoinDate.Format("2006-01-02"),
		BankName:          employee.BankName,
		BankAccountNumber: employee.BankAccountNumber,
		BankAccountHolder: employee.BankAccountHolder,
		ScheduleID:        scheduleID,
		Schedule:          scheduleDetail,
		EmploymentStatus:  employmentStatus,
		JobLevel:          jobLevel,
		Gender:            employee.Gender,
		Division:          division,
		DepartmentID:      deptID,
		DepartmentName:    deptName,
		CreatedAt:         employee.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         employee.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEmployeeResponseFromDB converts EmployeeWithUser from database join to EmployeeResponse DTO
func ToEmployeeResponseFromDB(employee *repository.EmployeeWithUser) *dto.EmployeeResponse {
	var scheduleID *string
	var scheduleDetail *dto.ScheduleDetail

	// Defensive: Check if employee is nil
	if employee == nil {
		return nil
	}

	// Only build schedule detail if both ScheduleID and ScheduleName are not nil
	if employee.ScheduleID != nil && employee.ScheduleName != nil && employee.ScheduleTimeIn != nil && employee.ScheduleTimeOut != nil {
		sid := employee.ScheduleID.String()
		scheduleID = &sid

		// Include schedule detail if available
		allowedLate := 0
		if employee.ScheduleAllowedLateMinutes != nil {
			allowedLate = *employee.ScheduleAllowedLateMinutes
		}

		scheduleDetail = &dto.ScheduleDetail{
			ID:                 sid,
			Name:               *employee.ScheduleName,
			TimeIn:             *employee.ScheduleTimeIn,
			TimeOut:            *employee.ScheduleTimeOut,
			AllowedLateMinutes: allowedLate,
		}
	}

	// Handle empty employment_status (default to PROBATION if empty)
	employmentStatus := employee.EmploymentStatus
	if employmentStatus == "" {
		employmentStatus = "PROBATION"
	}

	// Handle empty job_level (default to STAFF if empty)
	jobLevel := employee.JobLevel
	if jobLevel == "" {
		jobLevel = "STAFF"
	}

	// Handle empty division (default to GENERAL if empty)
	division := employee.Division
	if division == "" {
		division = "GENERAL"
	}

	// Handle nullable UserName and UserEmail
	userName := ""
	if employee.UserName != nil {
		userName = *employee.UserName
	}

	userEmail := ""
	if employee.UserEmail != nil {
		userEmail = *employee.UserEmail
	}

	// Handle DepartmentID
	var deptID *string
	if employee.DepartmentID != nil {
		did := employee.DepartmentID.String()
		deptID = &did
	}

	// Handle DepartmentName
	deptName := ""
	if employee.DepartmentName != nil {
		deptName = *employee.DepartmentName
	}

	return &dto.EmployeeResponse{
		ID:                employee.ID.String(),
		UserID:            employee.UserID.String(),
		UserName:          userName,
		UserEmail:         userEmail,
		Position:          employee.Position,
		PhoneNumber:       employee.PhoneNumber,
		Address:           employee.Address,
		SalaryBase:        employee.SalaryBase,
		JoinDate:          employee.JoinDate.Format("2006-01-02"),
		BankName:          employee.BankName,
		BankAccountNumber: employee.BankAccountNumber,
		BankAccountHolder: employee.BankAccountHolder,
		ScheduleID:        scheduleID,
		Schedule:          scheduleDetail,
		EmploymentStatus:  employmentStatus,
		JobLevel:          jobLevel,
		Gender:            employee.Gender,
		Division:          division,
		DepartmentID:      deptID,
		DepartmentName:    deptName,
		CreatedAt:         employee.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         employee.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEmployeeResponses converts array of EmployeeWithUser to array of EmployeeResponse pointers
func ToEmployeeResponses(employees []repository.EmployeeWithUser) []*dto.EmployeeResponse {
	responses := make([]*dto.EmployeeResponse, 0, len(employees))
	for _, employee := range employees {
		response := ToEmployeeResponseFromDB(&employee)
		if response != nil {
			responses = append(responses, response)
		}
	}
	return responses
}
