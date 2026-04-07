package entity

import "time"

type OvertimePolicy struct {
	ID                       string     `json:"id" db:"id"`
	Name                     string     `json:"name" db:"name"`
	Description              *string    `json:"description" db:"description"`
	RateType                 string     `json:"rateType" db:"rate_type"` // FIXED, MULTIPLIER
	RateMultiplier           *float64   `json:"rateMultiplier" db:"rate_multiplier"`
	FixedAmount              *float64   `json:"fixedAmount" db:"fixed_amount"`
	MinOvertimeMinutes       int        `json:"minOvertimeMinutes" db:"min_overtime_minutes"`
	MaxOvertimeHoursPerDay   *float64   `json:"maxOvertimeHoursPerDay" db:"max_overtime_hours_per_day"`
	MaxOvertimeHoursPerMonth *float64   `json:"maxOvertimeHoursPerMonth" db:"max_overtime_hours_per_month"`
	RequiresApproval         bool       `json:"requiresApproval" db:"requires_approval"`
	IsActive                 bool       `json:"isActive" db:"is_active"`
	CreatedAt                time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt                time.Time  `json:"updatedAt" db:"updated_at"`
	CompanyID                string     `json:"companyId" db:"company_id"`
}

type OvertimeRequest struct {
	ID               string     `json:"id" db:"id"`
	EmployeeID       string     `json:"employeeId" db:"employee_id"`
	OvertimeDate     time.Time  `json:"overtimeDate" db:"overtime_date"`
	StartTime        string     `json:"startTime" db:"start_time"`
	EndTime          string     `json:"endTime" db:"end_time"`
	TotalHours       float64    `json:"totalHours" db:"total_hours"`
	Reason           string     `json:"reason" db:"reason"`
	OvertimePolicyID string     `json:"overtimePolicyId" db:"overtime_policy_id"`
	Status           string     `json:"status" db:"status"` // PENDING, APPROVED, REJECTED, CANCELLED
	ApprovedBy       *string    `json:"approvedBy" db:"approved_by"`
	ApprovedAt       *time.Time `json:"approvedAt" db:"approved_at"`
	RejectionReason  *string    `json:"rejectionReason" db:"rejection_reason"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" db:"updated_at"`
	CompanyID        string     `json:"companyId" db:"company_id"`
}

type OvertimeAttendance struct {
	ID                 string     `json:"id" db:"id"`
	OvertimeRequestID  string     `json:"overtimeRequestId" db:"overtime_request_id"`
	EmployeeID         string     `json:"employeeId" db:"employee_id"`
	ClockInTime        *time.Time `json:"clockInTime" db:"clock_in_time"`
	ClockOutTime       *time.Time `json:"clockOutTime" db:"clock_out_time"`
	ActualHours        float64    `json:"actualHours" db:"actual_hours"`
	Notes              string     `json:"notes" db:"notes"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
}
