package helper

import (
	"time"

	"hris/internal/payroll/dto"
	"hris/internal/payroll/entity"
)

func PayrollToResponse(payroll *entity.PayrollWithItems, employeeName, bankName, bankAccountNumber, bankAccountHolder string) *dto.PayrollResponse {
	items := make([]dto.PayrollItemResponse, len(payroll.Items))
	for i, item := range payroll.Items {
		items[i] = dto.PayrollItemResponse{
			ID:     item.ID,
			Name:   item.Name,
			Amount: item.Amount,
			Type:   item.Type,
		}
	}

	return &dto.PayrollResponse{
		ID:                payroll.Payroll.ID,
		EmployeeID:        payroll.Payroll.EmployeeID,
		EmployeeName:      employeeName,
		BankName:          bankName,
		BankAccountNumber: bankAccountNumber,
		BankAccountHolder: bankAccountHolder,
		PeriodStart:       payroll.Payroll.PeriodStart.Format("2006-01-02"),
		PeriodEnd:         payroll.Payroll.PeriodEnd.Format("2006-01-02"),
		BaseSalary:        payroll.Payroll.BaseSalary,
		TotalAllowance:    payroll.Payroll.TotalAllowance,
		TotalDeduction:    payroll.Payroll.TotalDeduction,
		NetSalary:         payroll.Payroll.NetSalary,
		Status:            payroll.Payroll.Status,
		Items:             items,
		GeneratedAt:       payroll.Payroll.GeneratedAt.Format(time.RFC3339),
		CreatedAt:         payroll.Payroll.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         payroll.Payroll.UpdatedAt.Format(time.RFC3339),
	}
}

func PayrollToListResponse(payroll *entity.Payroll, employeeName string) *dto.PayrollListResponse {
	return &dto.PayrollListResponse{
		ID:           payroll.ID,
		EmployeeName: employeeName,
		Period:       payroll.PeriodStart.Format("2006-01") + " - " + payroll.PeriodEnd.Format("2006-01"),
		NetSalary:    payroll.NetSalary,
		Status:       payroll.Status,
		GeneratedAt:  payroll.GeneratedAt.Format(time.RFC3339),
	}
}
