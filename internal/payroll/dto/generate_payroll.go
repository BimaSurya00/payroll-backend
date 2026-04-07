package dto

type GeneratePayrollRequest struct {
	PeriodMonth int `json:"periodMonth" validate:"required,min=1,max=12"`
	PeriodYear  int `json:"periodYear" validate:"required,min=2020,max=2100"`
}

type GeneratePayrollResponse struct {
	TotalGenerated int    `json:"totalGenerated"`
	PeriodStart    string `json:"periodStart"`
	PeriodEnd      string `json:"periodEnd"`
	Message        string `json:"message"`
}
