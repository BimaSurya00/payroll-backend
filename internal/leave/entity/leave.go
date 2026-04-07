package entity

import "time"

type LeaveType struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Code             string    `json:"code" db:"code"`
	Description      *string   `json:"description" db:"description"`
	MaxDaysPerYear   int       `json:"maxDaysPerYear" db:"max_days_per_year"`
	DefaultDays      int       `json:"defaultDays" db:"default_days"`
	IsPaid           bool      `json:"isPaid" db:"is_paid"`
	RequiresApproval bool      `json:"requiresApproval" db:"requires_approval"`
	IsActive         bool      `json:"isActive" db:"is_active"`
	Color            string    `json:"color" db:"color"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
	CompanyID        string    `json:"companyId" db:"company_id"`
}

type LeaveBalance struct {
	ID          string    `json:"id" db:"id"`
	EmployeeID  string    `json:"employeeId" db:"employee_id"`
	LeaveTypeID string    `json:"leaveTypeId" db:"leave_type_id"`
	Year        int       `json:"year" db:"year"`
	Balance     int       `json:"balance" db:"balance"`
	Used        int       `json:"used" db:"used"`
	Pending     int       `json:"pending" db:"pending"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	CompanyID   string    `json:"companyId" db:"company_id"`
}

// Available returns the available leave balance
func (lb *LeaveBalance) Available() int {
	return lb.Balance - lb.Used - lb.Pending
}

type LeaveRequest struct {
	ID               string     `json:"id" db:"id"`
	EmployeeID       string     `json:"employeeId" db:"employee_id"`
	LeaveTypeID      string     `json:"leaveTypeId" db:"leave_type_id"`
	StartDate        time.Time  `json:"startDate" db:"start_date"`
	EndDate          time.Time  `json:"endDate" db:"end_date"`
	TotalDays        int        `json:"totalDays" db:"total_days"`
	Reason           string     `json:"reason" db:"reason"`
	AttachmentURL    *string    `json:"attachmentUrl" db:"attachment_url"`
	EmergencyContact *string    `json:"emergencyContact" db:"emergency_contact"`
	Status           string     `json:"status" db:"status"` // PENDING, APPROVED, REJECTED, CANCELLED
	ApprovedBy       *string    `json:"approvedBy" db:"approved_by"`
	ApprovedAt       *time.Time `json:"approvedAt" db:"approved_at"`
	RejectionReason  *string    `json:"rejectionReason" db:"rejection_reason"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" db:"updated_at"`
	CompanyID        string     `json:"companyId" db:"company_id"`
}
