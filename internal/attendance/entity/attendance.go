package entity

import "time"

type Attendance struct {
	ID             string     `json:"id" db:"id"`
	CompanyID      string     `json:"companyId" db:"company_id"`
	EmployeeID     string     `json:"employeeId" db:"employee_id"`
	ScheduleID     *string    `json:"scheduleId" db:"schedule_id"`
	Date           time.Time  `json:"date" db:"date"`
	ClockInTime    *time.Time `json:"clockInTime,omitempty" db:"clock_in_time"`
	ClockOutTime   *time.Time `json:"clockOutTime,omitempty" db:"clock_out_time"`
	ClockInLat     *float64   `json:"clockInLat,omitempty" db:"clock_in_lat"`
	ClockInLong    *float64   `json:"clockInLong,omitempty" db:"clock_in_long"`
	ClockOutLat    *float64   `json:"clockOutLat,omitempty" db:"clock_out_lat"`
	ClockOutLong   *float64   `json:"clockOutLong,omitempty" db:"clock_out_long"`
	Status         string     `json:"status" db:"status"` // PRESENT, LATE, ABSENT, LEAVE
	Notes          string     `json:"notes" db:"notes"`
	CorrectedBy    *string    `json:"correctedBy,omitempty" db:"corrected_by"`
	CorrectedAt    *time.Time `json:"correctedAt,omitempty" db:"corrected_at"`
	CorrectionNote *string    `json:"correctionNote,omitempty" db:"correction_note"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
}
