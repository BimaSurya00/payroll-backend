package helper

import (
	"time"

	"hris/internal/payroll/entity"
)

// Deprecated: Use CalculateSalaryFromConfig instead
const (
	TransportAllowance = 500000 // Rp 500.000
	MealAllowance      = 300000 // Rp 300.000
	LateDeductionPerDay = 50000 // Rp 50.000 per late day
)

// CalculateSalary computes net salary with allowances and deductions
// Deprecated: Use CalculateSalaryFromConfig for configurable values
func CalculateSalary(baseSalary float64, lateDays int) (allowance, deduction, netSalary float64) {
	// Calculate fixed allowances
	allowance = TransportAllowance + MealAllowance

	// Calculate late deduction
	deduction = float64(lateDays) * LateDeductionPerDay

	// Calculate net salary
	netSalary = baseSalary + allowance - deduction

	if netSalary < 0 {
		netSalary = 0
	}

	return
}

// CalculateSalaryFromConfig computes net salary based on database configs
func CalculateSalaryFromConfig(baseSalary float64, lateDays, absentDays int,
	configs []*entity.PayrollConfig) (totalAllowance, totalDeduction, netSalary float64, items []*entity.PayrollItem) {

	for _, config := range configs {
		if !config.IsActive {
			continue
		}

		switch config.Type {
		case "EARNING":
			totalAllowance += config.Amount
			items = append(items, &entity.PayrollItem{
				Name:   config.Name,
				Amount: config.Amount,
				Type:   "EARNING",
			})
		case "DEDUCTION":
			var deductionAmount float64
			switch config.CalculationType {
			case "PER_DAY":
				if config.Code == "LATE_DEDUCTION" {
					deductionAmount = config.Amount * float64(lateDays)
				} else if config.Code == "ABSENT_DEDUCTION" {
					dailySalary := baseSalary / 22
					deductionAmount = dailySalary * float64(absentDays)
				}
			case "FIXED":
				deductionAmount = config.Amount
			}
			if deductionAmount > 0 {
				totalDeduction += deductionAmount
				items = append(items, &entity.PayrollItem{
					Name:   config.Name,
					Amount: deductionAmount,
					Type:   "DEDUCTION",
				})
			}
		}
	}

	netSalary = baseSalary + totalAllowance - totalDeduction
	if netSalary < 0 {
		netSalary = 0
	}

	return
}

// GetPeriodRange returns start and end date for given month/year
func GetPeriodRange(month, year int) (time.Time, time.Time) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	return start, end
}

