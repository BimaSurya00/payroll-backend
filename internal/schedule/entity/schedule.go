package entity

import "time"

type Schedule struct {
	ID                 string    `json:"id" db:"id"`
	CompanyID          string    `json:"companyId" db:"company_id"`
	Name               string    `json:"name" db:"name"`
	TimeIn             string    `json:"timeIn" db:"time_in"`
	TimeOut            string    `json:"timeOut" db:"time_out"`
	AllowedLateMinutes int       `json:"allowedLateMinutes" db:"allowed_late_minutes"`
	Description        string    `json:"description" db:"description"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`
}
