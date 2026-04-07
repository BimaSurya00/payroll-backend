package dto

type ClockOutRequest struct {
	Lat  float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Long float64 `json:"long" validate:"required,gte=-180,lte=180"`
}

type ClockOutResponse struct {
	AttendanceID string  `json:"attendanceId"`
	ClockOutTime string  `json:"clockOutTime"`
	Distance     float64 `json:"distance"`
}
