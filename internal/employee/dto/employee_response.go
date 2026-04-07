package dto

type ScheduleDetail struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	TimeIn             string `json:"timeIn"`
	TimeOut            string `json:"timeOut"`
	AllowedLateMinutes int    `json:"allowedLateMinutes"`
}

type EmployeeResponse struct {
	ID                string          `json:"id"`
	UserID            string          `json:"userId"`
	UserName          string          `json:"userName"`
	UserEmail         string          `json:"userEmail"`
	Position          string          `json:"position"`
	PhoneNumber       string          `json:"phoneNumber"`
	Address           string          `json:"address"`
	SalaryBase        float64         `json:"salaryBase"`
	JoinDate          string          `json:"joinDate"`
	BankName          string          `json:"bankName"`
	BankAccountNumber string          `json:"bankAccountNumber"`
	BankAccountHolder string          `json:"bankAccountHolder"`
	ScheduleID        *string         `json:"scheduleId"`
	Schedule          *ScheduleDetail `json:"schedule,omitempty"`
	EmploymentStatus  string          `json:"employmentStatus"`
	JobLevel          string          `json:"jobLevel"`
	Gender            string          `json:"gender"`
	Division          string          `json:"division"`
	DepartmentID      *string         `json:"departmentId,omitempty"`
	DepartmentName    string          `json:"departmentName"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
}
