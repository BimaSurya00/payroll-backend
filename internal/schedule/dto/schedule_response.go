package dto

type ScheduleResponse struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	TimeIn             string `json:"timeIn"`
	TimeOut            string `json:"timeOut"`
	AllowedLateMinutes int    `json:"allowedLateMinutes"`
	Description        string `json:"description"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}
