package dto

type CreateEmployeeRequest struct {
	// User Account
	Name     string `json:"name" validate:"required,min=3,max=100,trimmed_string"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,password_strength"`

	// Employee Details
	Position           string  `json:"position" validate:"required,min=3,max=100,trimmed_string"`
	PhoneNumber        string  `json:"phoneNumber" validate:"omitempty,max=20"`
	Address            string  `json:"address" validate:"omitempty,max=500"`
	SalaryBase         float64 `json:"salaryBase" validate:"required,gt=0"`
	JoinDate           string  `json:"joinDate" validate:"required,datetime=2006-01-02"`
	BankName           string  `json:"bankName" validate:"omitempty,max=100"`
	BankAccountNumber  string  `json:"bankAccountNumber" validate:"omitempty,max=50"`
	BankAccountHolder  string  `json:"bankAccountHolder" validate:"omitempty,max=100"`
	ScheduleID         string  `json:"scheduleId" validate:"omitempty,uuid"`

	// New Fields
	EmploymentStatus   string  `json:"employmentStatus" validate:"required,oneof=PERMANENT CONTRACT PROBATION"`
	JobLevel           string  `json:"jobLevel" validate:"required,oneof=CEO MANAGER SUPERVISOR STAFF"`
	Gender             string  `json:"gender" validate:"required,oneof=MALE FEMALE"`
	Division           string  `json:"division" validate:"omitempty,max=100"` // Deprecated, use departmentId
	DepartmentID       string  `json:"departmentId" validate:"required,uuid"`
}
