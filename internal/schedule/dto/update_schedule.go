package dto

type UpdateScheduleRequest struct {
	Name               *string `json:"name" validate:"omitempty,min=3,max=100,trimmed_string"`
	TimeIn             *string `json:"timeIn" validate:"omitempty"`
	TimeOut            *string `json:"timeOut" validate:"omitempty"`
	AllowedLateMinutes *int    `json:"allowedLateMinutes" validate:"omitempty,min=0,max=60"`
	Description        *string `json:"description" validate:"omitempty,max=500"`
}

func (r *UpdateScheduleRequest) HasUpdates() bool {
	return r.Name != nil || r.TimeIn != nil || r.TimeOut != nil ||
		r.AllowedLateMinutes != nil || r.Description != nil
}
