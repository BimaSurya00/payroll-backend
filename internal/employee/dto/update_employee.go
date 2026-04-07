package dto

type UpdateEmployeeRequest struct {
	Position           *string  `json:"position,omitempty" validate:"omitempty,min=3,max=100,trimmed_string"`
	PhoneNumber        *string  `json:"phoneNumber,omitempty" validate:"omitempty,max=20"`
	Address            *string  `json:"address,omitempty" validate:"omitempty,max=500"`
	SalaryBase         *float64 `json:"salaryBase,omitempty" validate:"omitempty,gt=0"`
	BankName           *string  `json:"bankName,omitempty" validate:"omitempty,max=100"`
	BankAccountNumber  *string  `json:"bankAccountNumber,omitempty" validate:"omitempty,max=50"`
	BankAccountHolder  *string  `json:"bankAccountHolder,omitempty" validate:"omitempty,max=100"`
	ScheduleID         *string  `json:"scheduleId,omitempty" validate:"omitempty,uuid"`

	// New Fields
	EmploymentStatus   *string  `json:"employmentStatus,omitempty" validate:"omitempty,oneof=PERMANENT CONTRACT PROBATION"`
	JobLevel           *string  `json:"jobLevel,omitempty" validate:"omitempty,oneof=CEO MANAGER SUPERVISOR STAFF"`
	Gender             *string  `json:"gender,omitempty" validate:"omitempty,oneof=MALE FEMALE"`
	Division           *string  `json:"division,omitempty" validate:"omitempty,max=100"` // Deprecated
	DepartmentID       *string  `json:"departmentId,omitempty" validate:"omitempty,uuid"`
}

func (u *UpdateEmployeeRequest) HasUpdates() bool {
	return u.Position != nil || u.PhoneNumber != nil || u.Address != nil ||
		u.SalaryBase != nil || u.BankName != nil || u.BankAccountNumber != nil ||
		u.BankAccountHolder != nil || u.ScheduleID != nil ||
		u.EmploymentStatus != nil || u.JobLevel != nil || u.Gender != nil || 
		u.Division != nil || u.DepartmentID != nil
}
