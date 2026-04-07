package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	attendancerepository "hris/internal/attendance/repository"
	auditservice "hris/internal/audit/service"
	employeerepository "hris/internal/employee/repository"
	"hris/internal/payroll/dto"
	"hris/internal/payroll/entity"
	"hris/internal/payroll/helper"
	payrollrepository "hris/internal/payroll/repository"
)

var (
	ErrPayrollNotFound      = errors.New("payroll not found")
	ErrPayrollAlreadyExists = errors.New("payroll for this period already exists")
	ErrInvalidPeriod        = errors.New("invalid period")
	ErrInvalidStatus        = errors.New("invalid status transition")
	ErrNoEmployeesFound     = errors.New("no employees found")
	ErrGenerateFailed       = errors.New("failed to generate payroll")
	ErrEmployeeNotFound     = errors.New("employee not found")
)

type payrollServiceImpl struct {
	payrollRepo       payrollrepository.PayrollRepository
	employeeRepo      employeerepository.EmployeeRepository
	attendanceRepo    attendancerepository.AttendanceRepository
	payrollConfigRepo payrollrepository.PayrollConfigRepository
	auditService      auditservice.AuditService
	pool              *pgxpool.Pool
}

func NewPayrollService(
	payrollRepo payrollrepository.PayrollRepository,
	employeeRepo employeerepository.EmployeeRepository,
	attendanceRepo attendancerepository.AttendanceRepository,
	payrollConfigRepo payrollrepository.PayrollConfigRepository,
	auditService auditservice.AuditService,
	pool *pgxpool.Pool,
) PayrollService {
	return &payrollServiceImpl{
		payrollRepo:       payrollRepo,
		employeeRepo:      employeeRepo,
		attendanceRepo:    attendanceRepo,
		payrollConfigRepo: payrollConfigRepo,
		auditService:      auditService,
		pool:              pool,
	}
}

func (s *payrollServiceImpl) GenerateBulk(ctx context.Context, req *dto.GeneratePayrollRequest, companyID string) (*dto.GeneratePayrollResponse, error) {
	// Calculate period range
	periodStart, periodEnd := helper.GetPeriodRange(req.PeriodMonth, req.PeriodYear)
	periodStartStr := periodStart.Format("2006-01-02")
	periodEndStr := periodEnd.Format("2006-01-02")

	// Check if payroll already exists for this period and company
	existingPayrolls, err := s.payrollRepo.FindByPeriodAndCompany(ctx, companyID, periodStartStr, periodEndStr)
	if err == nil && len(existingPayrolls) > 0 {
		return nil, ErrPayrollAlreadyExists
	}

	// Fetch all employees for the company
	employees, err := s.employeeRepo.FindAllWithoutPaginationByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoEmployeesFound, err)
	}

	if len(employees) == 0 {
		return nil, ErrNoEmployeesFound
	}

	// Fetch payroll configs for calculation
	configs, err := s.payrollConfigRepo.FindAll(ctx)
	if err != nil {
		zap.L().Warn("failed to fetch payroll configs, using defaults", zap.Error(err))
		configs = nil
	}

	// === BEGIN TRANSACTION ===
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback if not committed

	generatedCount := 0
	now := time.Now()

	// Generate payroll for each employee
	for _, emp := range employees {
		// Hitung attendance summary dari database
		summary, err := s.attendanceRepo.GetAttendanceSummaryByPeriod(ctx, emp.ID, periodStart, periodEnd)
		lateDays := 0
		absentDays := 0
		if err == nil && summary != nil {
			lateDays = summary.TotalLate
			absentDays = summary.TotalAbsent
		}

		// Calculate salary dengan config-based calculator
		var allowance, deduction, netSalary float64
		var items []*entity.PayrollItem

		if len(configs) > 0 {
			// Use config-based calculation
			allowance, deduction, netSalary, items = helper.CalculateSalaryFromConfig(
				emp.SalaryBase,
				lateDays,
				absentDays,
				configs,
			)
		} else {
			// Fallback to deprecated calculator
			allowance, deduction, netSalary = helper.CalculateSalary(
				emp.SalaryBase,
				lateDays,
			)
			// Create default items
			items = []*entity.PayrollItem{
				{
					ID:        uuid.New().String(),
					PayrollID: "", // Will set after payrollID is created
					Name:      "Transport Allowance",
					Amount:    helper.TransportAllowance,
					Type:      "EARNING",
					CreatedAt: now,
				},
				{
					ID:        uuid.New().String(),
					PayrollID: "", // Will set after payrollID is created
					Name:      "Meal Allowance",
					Amount:    helper.MealAllowance,
					Type:      "EARNING",
					CreatedAt: now,
				},
			}
			if deduction > 0 {
				items = append(items, &entity.PayrollItem{
					ID:        uuid.New().String(),
					PayrollID: "",
					Name:      "Late Deduction",
					Amount:    deduction,
					Type:      "DEDUCTION",
					CreatedAt: now,
				})
			}
		}

		// Create payroll
		payrollID := uuid.New().String()

		payroll := &entity.Payroll{
			ID:             payrollID,
			EmployeeID:     emp.ID.String(),
			PeriodStart:    periodStart,
			PeriodEnd:      periodEnd,
			BaseSalary:     emp.SalaryBase,
			TotalAllowance: allowance,
			TotalDeduction: deduction,
			NetSalary:      netSalary,
			Status:         "DRAFT",
			GeneratedAt:    now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// Update payrollID in items
		for _, item := range items {
			item.PayrollID = payrollID
		}

		// Add absent deduction if any
		if absentDays > 0 {
			dailySalary := emp.SalaryBase / 22 // asumsi 22 hari kerja
			absentDeduction := float64(absentDays) * dailySalary
			deduction += absentDeduction
			netSalary -= absentDeduction

			// Update total deduction in payroll
			payroll.TotalDeduction = deduction
			payroll.NetSalary = netSalary

			items = append(items, &entity.PayrollItem{
				ID:        uuid.New().String(),
				PayrollID: payrollID,
				Name:      fmt.Sprintf("Absent Deduction (%d days)", absentDays),
				Amount:    absentDeduction,
				Type:      "DEDUCTION",
				CreatedAt: now,
			})
		}

		// INSERT menggunakan transaction (tx, bukan pool)
		if err := s.insertPayrollWithTx(ctx, tx, payroll, items); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrGenerateFailed, err)
		}

		generatedCount++
	}

	// === COMMIT TRANSACTION ===
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Audit log — payroll batch generated
	_ = s.auditService.Log(ctx, auditservice.AuditEntry{
		Action:       "GENERATE",
		ResourceType: "payroll",
		NewData: map[string]interface{}{
			"period_start": periodStartStr,
			"period_end":   periodEndStr,
			"count":        generatedCount,
		},
	})

	return &dto.GeneratePayrollResponse{
		TotalGenerated: generatedCount,
		PeriodStart:    periodStartStr,
		PeriodEnd:      periodEndStr,
		Message:        fmt.Sprintf("Successfully generated %d payrolls", generatedCount),
	}, nil
}

func (s *payrollServiceImpl) GetMyPayrolls(ctx context.Context, userID string, page, perPage int, path string) (*helper.PayrollPagination, error) {
	// Find employee by userID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, employeerepository.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}

	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	payrolls, err := s.payrollRepo.FindByEmployeeID(ctx, employee.ID.String(), skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.payrollRepo.CountByEmployeeID(ctx, employee.ID.String())
	if err != nil {
		return nil, err
	}

	// Convert to list response
	employeeName := employee.UserName
	if employeeName == nil || *employeeName == "" {
		name := employee.Position
		employeeName = &name
	}

	data := make([]dto.PayrollListResponse, len(payrolls))
	for i, p := range payrolls {
		data[i] = *helper.PayrollToListResponse(p, *employeeName)
	}

	return helper.BuildPayrollPagination(data, page, perPage, total, path), nil
}

func (s *payrollServiceImpl) GetMyPayrollByID(ctx context.Context, userID string, payrollID string) (*dto.PayrollResponse, error) {
	// Verify ownership
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
	if err != nil {
		return nil, ErrPayrollNotFound
	}

	payrollUUID, err := uuid.Parse(payrollID)
	if err != nil {
		return nil, ErrPayrollNotFound
	}

	payrollWithItems, err := s.payrollRepo.FindByIDWithItems(ctx, payrollUUID)
	if err != nil {
		return nil, ErrPayrollNotFound
	}

	// SECURITY: Pastikan payroll milik employee ini
	if payrollWithItems.Payroll.EmployeeID != employee.ID.String() {
		return nil, ErrPayrollNotFound // Jangan expose info bahwa payroll ada tapi bukan miliknya
	}

	employeeName := employee.UserName
	if employeeName == nil || *employeeName == "" {
		name := employee.Position
		employeeName = &name
	}

	return helper.PayrollToResponse(payrollWithItems, *employeeName,
		employee.BankName, employee.BankAccountNumber, employee.BankAccountHolder), nil
}

// insertPayrollWithTx inserts payroll dan items menggunakan transaction
func (s *payrollServiceImpl) insertPayrollWithTx(ctx context.Context, tx pgx.Tx, payroll *entity.Payroll, items []*entity.PayrollItem) error {
	// Insert payroll
	_, err := tx.Exec(ctx,
		`INSERT INTO payrolls (id, employee_id, period_start, period_end, base_salary,
         total_allowance, total_deduction, net_salary, status, generated_at, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		payroll.ID, payroll.EmployeeID, payroll.PeriodStart, payroll.PeriodEnd,
		payroll.BaseSalary, payroll.TotalAllowance, payroll.TotalDeduction,
		payroll.NetSalary, payroll.Status, payroll.GeneratedAt, payroll.CreatedAt, payroll.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payroll: %w", err)
	}

	// Insert items
	for _, item := range items {
		_, err := tx.Exec(ctx,
			`INSERT INTO payroll_items (id, payroll_id, name, amount, type) VALUES ($1,$2,$3,$4,$5)`,
			item.ID, item.PayrollID, item.Name, item.Amount, item.Type,
		)
		if err != nil {
			return fmt.Errorf("failed to insert payroll item: %w", err)
		}
	}

	return nil
}

func (s *payrollServiceImpl) GetByID(ctx context.Context, id string, companyID string) (*dto.PayrollResponse, error) {
	payrollUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrPayrollNotFound
	}

	payrollWithItems, err := s.payrollRepo.FindByIDWithItemsAndCompany(ctx, payrollUUID, companyID)
	if err != nil {
		if err == payrollrepository.ErrPayrollNotFound {
			return nil, ErrPayrollNotFound
		}
		return nil, err
	}

	// Get employee info
	employeeUUID, err := uuid.Parse(payrollWithItems.Payroll.EmployeeID)
	if err != nil {
		return nil, err
	}

	employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
	if err != nil {
		return nil, err
	}

	// Use UserName if available, otherwise use Position as fallback
	employeeName := employee.UserName
	if employeeName == nil || *employeeName == "" {
		name := employee.Position
		employeeName = &name
	}

	response := helper.PayrollToResponse(payrollWithItems, *employeeName, employee.BankName, employee.BankAccountNumber, employee.BankAccountHolder)
	return response, nil
}

func (s *payrollServiceImpl) GetAll(ctx context.Context, page, perPage int, path string, companyID string) (*helper.PayrollPagination, error) {
	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	payrolls, err := s.payrollRepo.FindAllByCompany(ctx, companyID, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := s.payrollRepo.CountByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}

	// Collect unique employee IDs
	employeeIDSet := make(map[uuid.UUID]bool)
	for _, p := range payrolls {
		empUUID, err := uuid.Parse(p.EmployeeID)
		if err == nil {
			employeeIDSet[empUUID] = true
		}
	}

	employeeIDs := make([]uuid.UUID, 0, len(employeeIDSet))
	for id := range employeeIDSet {
		employeeIDs = append(employeeIDs, id)
	}

	// Batch fetch employees (1 query instead of N)
	employees, err := s.employeeRepo.FindByIDs(ctx, employeeIDs)
	employeeMap := make(map[string]*employeerepository.Employee)
	if err == nil {
		for _, emp := range employees {
			employeeMap[emp.ID.String()] = emp
		}
	}

	// Build response using map
	data := make([]dto.PayrollListResponse, len(payrolls))
	for i, payroll := range payrolls {
		employeeName := payroll.EmployeeID // fallback
		if emp, ok := employeeMap[payroll.EmployeeID]; ok {
			// Employee struct doesn't have UserName, use Position as name
			employeeName = emp.Position
			if employeeName == "" {
				employeeName = payroll.EmployeeID
			}
		}
		data[i] = *helper.PayrollToListResponse(payroll, employeeName)
	}

	return helper.BuildPayrollPagination(data, page, perPage, total, path), nil
}

func (s *payrollServiceImpl) UpdateStatus(ctx context.Context, id string, companyID string, req *dto.UpdatePayrollStatusRequest) error {
	payrollUUID, err := uuid.Parse(id)
	if err != nil {
		return ErrPayrollNotFound
	}

	// Validate status transition with company check
	payroll, err := s.payrollRepo.FindByIDAndCompany(ctx, payrollUUID, companyID)
	if err != nil {
		if err == payrollrepository.ErrPayrollNotFound {
			return ErrPayrollNotFound
		}
		return err
	}

	// Validate status transition rules
	if !isValidStatusTransition(payroll.Status, req.Status) {
		return ErrInvalidStatus
	}

	oldStatus := payroll.Status
	err = s.payrollRepo.UpdateStatus(ctx, payrollUUID, req.Status)
	if err != nil {
		if err == payrollrepository.ErrPayrollNotFound {
			return ErrPayrollNotFound
		}
		return err
	}

	// Audit log — payroll status updated
	_ = s.auditService.Log(ctx, auditservice.AuditEntry{
		Action:       "UPDATE",
		ResourceType: "payroll",
		ResourceID:   id,
		OldData:      map[string]interface{}{"status": oldStatus},
		NewData:      map[string]interface{}{"status": req.Status},
	})

	return nil
}

func (s *payrollServiceImpl) ExportCSV(ctx context.Context, month, year int, companyID string) ([]byte, string, error) {
	periodStart, periodEnd := helper.GetPeriodRange(month, year)

	// Fetch approved payrolls for company
	payrolls, err := s.payrollRepo.FindByPeriodAndCompany(ctx, companyID, periodStart.Format("2006-01-02"), periodEnd.Format("2006-01-02"))
	if err != nil {
		return nil, "", err
	}

	// Filter by status APPROVED only and convert to response
	var approvedPayrolls []*dto.PayrollResponse
	for _, p := range payrolls {
		if p.Status == "APPROVED" {
			payrollUUID, err := uuid.Parse(p.ID)
			if err != nil {
				continue
			}

			payrollWithItems, err := s.payrollRepo.FindByIDWithItems(ctx, payrollUUID)
			if err != nil {
				continue
			}

			employeeUUID, err := uuid.Parse(p.EmployeeID)
			if err != nil {
				continue
			}

			employee, err := s.employeeRepo.FindByID(ctx, employeeUUID)
			if err != nil {
				continue
			}

			// Use UserName if available, otherwise use Position as fallback
			employeeName := employee.UserName
			if employeeName == nil || *employeeName == "" {
				name := employee.Position
				employeeName = &name
			}

			response := helper.PayrollToResponse(payrollWithItems, *employeeName, employee.BankName, employee.BankAccountNumber, employee.BankAccountHolder)
			approvedPayrolls = append(approvedPayrolls, response)
		}
	}

	// Generate CSV
	buf, filename, err := helper.GeneratePayrollCSV(approvedPayrolls)
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), filename, nil
}

// isValidStatusTransition validates status transition rules
func isValidStatusTransition(currentStatus, newStatus string) bool {
	transitions := map[string][]string{
		"DRAFT":    {"APPROVED", "DRAFT"},
		"APPROVED": {"PAID", "APPROVED"},
		"PAID":     {"PAID"},
	}

	allowedTransitions, exists := transitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}

	return false
}
