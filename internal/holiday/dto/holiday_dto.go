package dto

type CreateHolidayRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=200"`
	Date        string `json:"date" validate:"required,datetime=2006-01-02"`
	Type        string `json:"type" validate:"required,oneof=NATIONAL COMPANY OPTIONAL"`
	IsRecurring bool   `json:"isRecurring"`
	Year        *int   `json:"year,omitempty" validate:"omitempty,min=2000"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type UpdateHolidayRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Date        *string `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Type        *string `json:"type,omitempty" validate:"omitempty,oneof=NATIONAL COMPANY OPTIONAL"`
	IsRecurring *bool   `json:"isRecurring,omitempty"`
	Year        *int    `json:"year,omitempty" validate:"omitempty,min=2000"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type HolidayResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	IsRecurring bool    `json:"isRecurring"`
	Year        *int    `json:"year,omitempty"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}
