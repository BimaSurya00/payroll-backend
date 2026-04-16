package dto

type CreateScheduleRequest struct {
	Name               string `json:"name" validate:"required,min=3,max=100,trimmed_string"`
	TimeIn             string `json:"timeIn" validate:"required"`
	TimeOut            string `json:"timeOut" validate:"required"`
	AllowedLateMinutes int    `json:"allowedLateMinutes" validate:"min=0,max=60"`
	Description        string `json:"description" validate:"omitempty,max=500"`
}
