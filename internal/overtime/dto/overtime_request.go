package dto

type CreateOvertimeRequestRequest struct {
	OvertimeDate     string  `json:"overtimeDate" validate:"required,datetime=2006-01-02"`
	StartTime        string  `json:"startTime" validate:"required"` // Format: HH:MM
	EndTime          string  `json:"endTime" validate:"required"`   // Format: HH:MM
	Reason           string  `json:"reason" validate:"required,min=10,max=1000"`
	OvertimePolicyID string  `json:"overtimePolicyId" validate:"required,uuid"`
	EmployeeID       string  `json:"employeeId" validate:"omitempty,uuid"` // For Admin proxy
}

type UpdateOvertimeRequestRequest struct {
	Reason *string `json:"reason" validate:"omitempty,min=10,max=1000"`
}

type ApproveOvertimeRequest struct {
	ApprovalNote string `json:"approvalNote" validate:"omitempty,max=500"`
}

type RejectOvertimeRequest struct {
	RejectionReason string `json:"rejectionReason" validate:"required,min=10,max=500"`
}

type OvertimeRequestResponse struct {
	ID                string              `json:"id"`
	EmployeeID        string              `json:"employeeId"`
	EmployeeName      string              `json:"employeeName"`
	EmployeePosition  string              `json:"employeePosition"`
	OvertimeDate      string              `json:"overtimeDate"`
	StartTime         string              `json:"startTime"`
	EndTime           string              `json:"endTime"`
	TotalHours        float64             `json:"totalHours"`
	Reason            string              `json:"reason"`
	OvertimePolicy    OvertimePolicyDetail `json:"overtimePolicy"`
	Status            string              `json:"status"`
	ApprovedBy        *string             `json:"approvedBy"`
	ApprovedByName    *string             `json:"approvedByName"`
	ApprovedAt        *string             `json:"approvedAt"`
	RejectionReason   string              `json:"rejectionReason"`
	CreatedAt         string              `json:"createdAt"`
	UpdatedAt         string              `json:"updatedAt"`
}

type OvertimePolicyDetail struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	RateType       string  `json:"rateType"`
	RateMultiplier float64 `json:"rateMultiplier"`
	FixedAmount    float64 `json:"fixedAmount"`
}

type ClockInRequest struct {
	Notes string `json:"notes" validate:"omitempty,max=500"`
}

type ClockOutRequest struct {
	Notes string `json:"notes" validate:"omitempty,max=500"`
}

type OvertimeAttendanceResponse struct {
	ID                string      `json:"id"`
	OvertimeRequestID string      `json:"overtimeRequestId"`
	EmployeeID        string      `json:"employeeId"`
	EmployeeName      string      `json:"employeeName"`
	ClockInTime       *string     `json:"clockInTime"`
	ClockOutTime      *string     `json:"clockOutTime"`
	ActualHours       float64     `json:"actualHours"`
	Notes             string      `json:"notes"`
	CreatedAt         string      `json:"createdAt"`
}

type OvertimeCalculationResponse struct {
	EmployeeID      string  `json:"employeeId"`
	EmployeeName    string  `json:"employeeName"`
	TotalHours      float64 `json:"totalHours"`
	RateType        string  `json:"rateType"`
	RateMultiplier  float64 `json:"rateMultiplier"`
	HourlyRate      float64 `json:"hourlyRate"`
	OvertimePay     float64 `json:"overtimePay"`
}
