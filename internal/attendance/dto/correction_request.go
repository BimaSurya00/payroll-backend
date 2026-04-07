package dto

type CreateCorrectionRequest struct {
	EmployeeID string `json:"employeeId" validate:"required,uuid"`
	Date       string `json:"date" validate:"required"`
	ClockIn    string `json:"clockIn" validate:"required"`
	ClockOut   string `json:"clockOut,omitempty"`
	Status     string `json:"status" validate:"required,oneof=PRESENT LATE ABSENT"`
	Notes      string `json:"notes" validate:"required,min=5"`
}

type UpdateCorrectionRequest struct {
	ClockIn  *string `json:"clockIn,omitempty" validate:"omitempty,datetime=15:04"`
	ClockOut *string `json:"clockOut,omitempty" validate:"omitempty,datetime=15:04"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=PRESENT LATE ABSENT"`
	Notes    *string `json:"notes,omitempty" validate:"omitempty,min=5"`
}
