package dto

type ScheduleResponse struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	TimeIn              string  `json:"timeIn"`
	TimeOut             string  `json:"timeOut"`
	AllowedLateMinutes  int     `json:"allowedLateMinutes"`
	OfficeLat           float64 `json:"officeLat"`
	OfficeLong          float64 `json:"officeLong"`
	AllowedRadiusMeters int     `json:"allowedRadiusMeters"`
	Description         string  `json:"description"`
	CreatedAt           string  `json:"createdAt"`
	UpdatedAt           string  `json:"updatedAt"`
}
