package dto

type CreateLeaveRequestRequest struct {
	LeaveTypeID      string  `json:"leaveTypeId" validate:"required,uuid"`
	StartDate        string  `json:"startDate" validate:"required,datetime=2006-01-02"`
	EndDate          string  `json:"endDate" validate:"required,datetime=2006-01-02"`
	Reason           string  `json:"reason" validate:"required,min=10,max=1000"`
	AttachmentURL    string  `json:"attachmentUrl" validate:"omitempty,url"`
	EmergencyContact string  `json:"emergencyContact" validate:"omitempty,max=20"`
}

type UpdateLeaveRequestRequest struct {
	Reason *string `json:"reason" validate:"omitempty,min=10,max=1000"`
}

type ApproveLeaveRequest struct {
	ApprovalNote string `json:"approvalNote" validate:"omitempty,max=500"`
}

type RejectLeaveRequest struct {
	RejectionReason string `json:"rejectionReason" validate:"required,min=10,max=500"`
}

type LeaveRequestResponse struct {
	ID                string          `json:"id"`
	EmployeeID        string          `json:"employeeId"`
	EmployeeName      string          `json:"employeeName"`
	EmployeePosition  string          `json:"employeePosition"`
	LeaveType         LeaveTypeDetail `json:"leaveType"`
	StartDate         string          `json:"startDate"`
	EndDate           string          `json:"endDate"`
	TotalDays         int             `json:"totalDays"`
	Reason            string          `json:"reason"`
	AttachmentURL     string          `json:"attachmentUrl"`
	EmergencyContact  string          `json:"emergencyContact"`
	Status            string          `json:"status"`
	ApprovedBy        *string         `json:"approvedBy"`
	ApprovedByName    *string         `json:"approvedByName"`
	ApprovedAt        *string         `json:"approvedAt"`
	RejectionReason   string          `json:"rejectionReason"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
}

type LeaveTypeDetail struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	IsPaid bool   `json:"isPaid"`
}

type LeaveBalanceResponse struct {
	EmployeeID   string              `json:"employeeId"`
	EmployeeName string              `json:"employeeName"`
	Year         int                 `json:"year"`
	Balances     []LeaveBalanceItem  `json:"balances"`
}

type LeaveBalanceItem struct {
	LeaveTypeID   string `json:"leaveTypeId"`
	LeaveTypeName string `json:"leaveTypeName"`
	Balance       int    `json:"balance"`
	Used          int    `json:"used"`
	Pending       int    `json:"pending"`
	Available     int    `json:"available"`
}
