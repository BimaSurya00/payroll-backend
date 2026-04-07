package dto

type PayrollResponse struct {
	ID                string                `json:"id"`
	EmployeeID        string                `json:"employeeId"`
	EmployeeName      string                `json:"employeeName"`
	BankName          string                `json:"bankName"`
	BankAccountNumber string                `json:"bankAccountNumber"`
	BankAccountHolder string                `json:"bankAccountHolder"`
	PeriodStart       string                `json:"periodStart"`
	PeriodEnd         string                `json:"periodEnd"`
	BaseSalary        float64               `json:"baseSalary"`
	TotalAllowance    float64               `json:"totalAllowance"`
	TotalDeduction    float64               `json:"totalDeduction"`
	NetSalary         float64               `json:"netSalary"`
	Status            string                `json:"status"`
	Items             []PayrollItemResponse `json:"items"`
	GeneratedAt       string                `json:"generatedAt"`
	CreatedAt         string                `json:"createdAt"`
	UpdatedAt         string                `json:"updatedAt"`
}

type PayrollItemResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Amount float64 `json:"amount"`
	Type   string `json:"type"`
}

type PayrollListResponse struct {
	ID          string  `json:"id"`
	EmployeeName string `json:"employeeName"`
	Period      string  `json:"period"`
	NetSalary   float64 `json:"netSalary"`
	Status      string  `json:"status"`
	GeneratedAt string  `json:"generatedAt"`
}
