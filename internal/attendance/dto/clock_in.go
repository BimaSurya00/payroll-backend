package dto

type ClockInRequest struct {
	Lat  float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Long float64 `json:"long" validate:"required,gte=-180,lte=180"`
}

type ClockInResponse struct {
	AttendanceID string  `json:"attendanceId"`
	EmployeeID   string  `json:"employeeId"`
	ClockInTime  string  `json:"clockInTime"`
	Status       string  `json:"status"`
	Distance     float64 `json:"distance"` // Distance from office in meters
	ScheduleName string  `json:"scheduleName"`
}
