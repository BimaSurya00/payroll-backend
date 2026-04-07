package dto

type UpdateScheduleRequest struct {
	Name                *string  `json:"name" validate:"omitempty,min=3,max=100,trimmed_string"`
	TimeIn              *string  `json:"timeIn" validate:"omitempty"`
	TimeOut             *string  `json:"timeOut" validate:"omitempty"`
	AllowedLateMinutes  *int     `json:"allowedLateMinutes" validate:"omitempty,min=0,max=60"`
	OfficeLat           *float64 `json:"officeLat" validate:"omitempty,gte=-90,lte=90"`
	OfficeLong          *float64 `json:"officeLong" validate:"omitempty,gte=-180,lte=180"`
	AllowedRadiusMeters *int     `json:"allowedRadiusMeters" validate:"omitempty,min=10,max=1000"`
	Description         *string  `json:"description" validate:"omitempty,max=500"`
}

func (r *UpdateScheduleRequest) HasUpdates() bool {
	return r.Name != nil || r.TimeIn != nil || r.TimeOut != nil ||
		r.AllowedLateMinutes != nil || r.OfficeLat != nil ||
		r.OfficeLong != nil || r.AllowedRadiusMeters != nil ||
		r.Description != nil
}
