package dto

type CreateScheduleRequest struct {
	Name                string  `json:"name" validate:"required,min=3,max=100,trimmed_string"`
	TimeIn              string  `json:"timeIn" validate:"required"`
	TimeOut             string  `json:"timeOut" validate:"required"`
	AllowedLateMinutes  int     `json:"allowedLateMinutes" validate:"min=0,max=60"`
	OfficeLat           float64 `json:"officeLat" validate:"required,gte=-90,lte=90"`
	OfficeLong          float64 `json:"officeLong" validate:"required,gte=-180,lte=180"`
	AllowedRadiusMeters int     `json:"allowedRadiusMeters" validate:"min=10,max=1000"`
	Description         string  `json:"description" validate:"omitempty,max=500"`
}
