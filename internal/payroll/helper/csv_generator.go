package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"hris/internal/payroll/dto"
)

// GeneratePayrollCSV creates CSV file for bank transfer
func GeneratePayrollCSV(payrolls []*dto.PayrollResponse) (*bytes.Buffer, string, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Write header
	header := []string{"Bank Name", "Account Number", "Account Holder", "Amount", "Description"}
	if err := writer.Write(header); err != nil {
		return nil, "", err
	}

	// Write data rows
	for _, p := range payrolls {
		description := fmt.Sprintf("Payroll %s - %s",
			p.PeriodStart, p.EmployeeName)

		row := []string{
			p.BankName,
			p.BankAccountNumber,
			p.BankAccountHolder,
			fmt.Sprintf("%.2f", p.NetSalary),
			description,
		}

		if err := writer.Write(row); err != nil {
			return nil, "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	// Generate filename
	filename := fmt.Sprintf("payroll_export_%s.csv", time.Now().Format("20060102_150405"))

	return buf, filename, nil
}
