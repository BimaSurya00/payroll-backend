package dto

type AttendanceResponse struct {
	ID           string  `json:"id"`
	EmployeeID   string  `json:"employeeId"`
	EmployeeName string  `json:"employeeName"`
	Date         string  `json:"date"`
	ClockInTime  *string `json:"clockInTime,omitempty"`
	ClockOutTime *string `json:"clockOutTime,omitempty"`
	Status       string  `json:"status"`
	Notes        string  `json:"notes"`
	ScheduleName *string `json:"scheduleName,omitempty"`
}
