package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/dashboard/dto"
	sharedHelper "hris/shared/helper"
)

type DashboardService interface {
	GetSummary(ctx context.Context, companyID string) (*dto.DashboardSummary, error)
	GetAttendanceStats(ctx context.Context, companyID string, startDate, endDate string) (*dto.AttendanceStats, error)
	GetPayrollStats(ctx context.Context, companyID string, month, year int) (*dto.PayrollStats, error)
	GetEmployeeStats(ctx context.Context, companyID string) (*dto.EmployeeStats, error)
	GetRecentActivities(ctx context.Context, companyID string, limit int) (*dto.RecentActivitiesResponse, error)
}

type dashboardService struct {
	pool *pgxpool.Pool
}

func NewDashboardService(pool *pgxpool.Pool) DashboardService {
	return &dashboardService{
		pool: pool,
	}
}

func (s *dashboardService) GetSummary(ctx context.Context, companyID string) (*dto.DashboardSummary, error) {
	today := sharedHelper.Today()
	now := sharedHelper.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, sharedHelper.GetLocation())

	// 1. Attendance summary — filtered by company
	attendanceSummary := dto.AttendanceSummary{}
	err := s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN a.status = 'PRESENT' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN a.status = 'LATE'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN a.status = 'ABSENT'  THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN a.status = 'LEAVE'   THEN 1 ELSE 0 END), 0),
			(SELECT COUNT(*) FROM employees WHERE company_id = $2)
		FROM attendances a
		WHERE a.date = $1 AND a.company_id = $2
	`, today, companyID).Scan(
		&attendanceSummary.TodayPresent,
		&attendanceSummary.TodayLate,
		&attendanceSummary.TodayAbsent,
		&attendanceSummary.TodayLeave,
		&attendanceSummary.TotalEmployees,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance summary: %w", err)
	}

	// 2. Leave summary — filtered by company
	leaveSummary := dto.LeaveSummary{}
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'APPROVED' AND created_at >= $1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'REJECTED' AND created_at >= $1 THEN 1 ELSE 0 END), 0)
		FROM leave_requests
		WHERE company_id = $2
			AND status IN ('PENDING', 'APPROVED', 'REJECTED')
			AND (status = 'PENDING' OR created_at >= $1)
	`, monthStart, companyID).Scan(
		&leaveSummary.PendingRequests,
		&leaveSummary.ApprovedThisMonth,
		&leaveSummary.RejectedThisMonth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query leave summary: %w", err)
	}

	// 3. Payroll summary — filtered by company
	payrollSummary := dto.PayrollSummary{
		CurrentPeriod: today.Format("January 2006"),
	}
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'DRAFT'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'PAID'     THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(net_salary), 0)
		FROM payrolls
		WHERE company_id = $1 AND period_start >= $2 AND period_start <= $3
	`, companyID, monthStart, today).Scan(
		&payrollSummary.DraftCount,
		&payrollSummary.ApprovedCount,
		&payrollSummary.PaidCount,
		&payrollSummary.TotalNetSalary,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query payroll summary: %w", err)
	}

	// 4. Employee summary — filtered by company
	employeeSummary := dto.EmployeeSummary{}
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN employment_status NOT IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN created_at >= $1 THEN 1 ELSE 0 END), 0)
		FROM employees
		WHERE company_id = $2
	`, monthStart, companyID).Scan(
		&employeeSummary.TotalActive,
		&employeeSummary.TotalInactive,
		&employeeSummary.NewThisMonth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query employee summary: %w", err)
	}

	return &dto.DashboardSummary{
		Attendance: attendanceSummary,
		Leave:      leaveSummary,
		Payroll:    payrollSummary,
		Employee:   employeeSummary,
	}, nil
}

// GetAttendanceStats returns detailed attendance statistics for a period
func (s *dashboardService) GetAttendanceStats(ctx context.Context, companyID string, startDate, endDate string) (*dto.AttendanceStats, error) {
	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	// Get daily breakdown - filtered by company
	query := `
		SELECT 
			date,
			COALESCE(SUM(CASE WHEN status = 'PRESENT' THEN 1 ELSE 0 END), 0) as present,
			COALESCE(SUM(CASE WHEN status = 'LATE' THEN 1 ELSE 0 END), 0) as late,
			COALESCE(SUM(CASE WHEN status = 'ABSENT' THEN 1 ELSE 0 END), 0) as absent,
			COALESCE(SUM(CASE WHEN status = 'LEAVE' THEN 1 ELSE 0 END), 0) as leave_count,
			COUNT(*) as total
		FROM attendances
		WHERE date BETWEEN $1 AND $2 AND company_id = $3
		GROUP BY date
		ORDER BY date
	`

	rows, err := s.pool.Query(ctx, query, startDate, endDate, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance stats: %w", err)
	}
	defer rows.Close()

	var dailyBreakdown []dto.DailyAttendanceStat
	var totalPresent, totalLate, totalAbsent, totalLeave, dayCount int

	for rows.Next() {
		var stat dto.DailyAttendanceStat
		var leaveCount int
		err := rows.Scan(&stat.Date, &stat.Present, &stat.Late, &stat.Absent, &leaveCount, &stat.Total)
		if err != nil {
			continue
		}
		stat.Leave = leaveCount
		dailyBreakdown = append(dailyBreakdown, stat)
		totalPresent += stat.Present
		totalLate += stat.Late
		totalAbsent += stat.Absent
		totalLeave += stat.Leave
		dayCount++
	}

	// Calculate averages
	var avgPresent, avgLate float64
	if dayCount > 0 {
		avgPresent = float64(totalPresent) / float64(dayCount)
		avgLate = float64(totalLate) / float64(dayCount)
	}

	return &dto.AttendanceStats{
		Period:         start.Format("Jan 2, 2006") + " - " + end.Format("Jan 2, 2006"),
		DailyBreakdown: dailyBreakdown,
		Summary: dto.AttendancePeriodSummary{
			TotalPresent: totalPresent,
			TotalLate:    totalLate,
			TotalAbsent:  totalAbsent,
			TotalLeave:   totalLeave,
			AvgPresent:   avgPresent,
			AvgLate:      avgLate,
		},
	}, nil
}

// GetPayrollStats returns detailed payroll statistics
func (s *dashboardService) GetPayrollStats(ctx context.Context, companyID string, month, year int) (*dto.PayrollStats, error) {
	periodStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, -1)
	periodStr := periodStart.Format("January 2006")

	// Get status breakdown - filtered by company
	var statusBreakdown dto.PayrollStatusBreakdown
	err := s.pool.QueryRow(ctx, `
		SELECT 
			COALESCE(SUM(CASE WHEN status = 'DRAFT' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'PAID' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'DRAFT' THEN net_salary ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'APPROVED' THEN net_salary ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'PAID' THEN net_salary ELSE 0 END), 0)
		FROM payrolls
		WHERE company_id = $1 AND period_start >= $2 AND period_start <= $3
	`, companyID, periodStart, periodEnd).Scan(
		&statusBreakdown.DraftCount,
		&statusBreakdown.ApprovedCount,
		&statusBreakdown.PaidCount,
		&statusBreakdown.DraftAmount,
		&statusBreakdown.ApprovedAmount,
		&statusBreakdown.PaidAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query payroll status breakdown: %w", err)
	}

	// Get total payrolls and average - filtered by company
	var totalPayrolls int
	var totalAmount, averageSalary float64
	err = s.pool.QueryRow(ctx, `
		SELECT 
			COUNT(*),
			COALESCE(SUM(net_salary), 0),
			COALESCE(AVG(net_salary), 0)
		FROM payrolls
		WHERE company_id = $1 AND period_start >= $2 AND period_start <= $3
	`, companyID, periodStart, periodEnd).Scan(&totalPayrolls, &totalAmount, &averageSalary)
	if err != nil {
		return nil, fmt.Errorf("failed to query payroll totals: %w", err)
	}

	// Get department-wise payroll stats - filtered by company
	deptRows, err := s.pool.Query(ctx, `
		SELECT 
			d.id,
			d.name,
			COUNT(e.id),
			COALESCE(SUM(p.net_salary), 0)
		FROM departments d
		LEFT JOIN employees e ON e.department_id = d.id AND e.company_id = d.company_id
		LEFT JOIN payrolls p ON p.employee_id = e.id 
			AND p.period_start >= $1 AND p.period_start <= $2 AND p.company_id = $3
		WHERE d.company_id = $3
		GROUP BY d.id, d.name
		ORDER BY d.name
	`, periodStart, periodEnd, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query department payroll stats: %w", err)
	}
	defer deptRows.Close()

	var deptStats []dto.DepartmentPayrollStat
	for deptRows.Next() {
		var stat dto.DepartmentPayrollStat
		err := deptRows.Scan(&stat.DepartmentID, &stat.DepartmentName, &stat.EmployeeCount, &stat.TotalPayroll)
		if err != nil {
			continue
		}
		deptStats = append(deptStats, stat)
	}

	return &dto.PayrollStats{
		Period:          periodStr,
		TotalPayrolls:   totalPayrolls,
		TotalAmount:     totalAmount,
		AverageSalary:   averageSalary,
		StatusBreakdown: statusBreakdown,
		DepartmentStats: deptStats,
	}, nil
}

// GetEmployeeStats returns detailed employee statistics
func (s *dashboardService) GetEmployeeStats(ctx context.Context, companyID string) (*dto.EmployeeStats, error) {
	// Get total count and status breakdown - filtered by company
	var totalCount int
	var statusBreakdown dto.EmployeeStatusBreakdown
	err := s.pool.QueryRow(ctx, `
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN employment_status = 'PERMANENT' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN employment_status = 'CONTRACT' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN employment_status = 'PROBATION' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN employment_status = 'INTERN' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN employment_status = 'RESIGNED' THEN 1 ELSE 0 END), 0)
		FROM employees
		WHERE company_id = $1
	`, companyID).Scan(
		&totalCount,
		&statusBreakdown.Permanent,
		&statusBreakdown.Contract,
		&statusBreakdown.Probation,
		&statusBreakdown.Intern,
		&statusBreakdown.Resigned,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query employee status breakdown: %w", err)
	}

	// Get department stats - filtered by company
	deptRows, err := s.pool.Query(ctx, `
		SELECT 
			d.id,
			d.name,
			COUNT(e.id) as total,
			COALESCE(SUM(CASE WHEN e.employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0) as active
		FROM departments d
		LEFT JOIN employees e ON e.department_id = d.id AND e.company_id = d.company_id
		WHERE d.company_id = $1
		GROUP BY d.id, d.name
		ORDER BY total DESC
	`, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query department stats: %w", err)
	}
	defer deptRows.Close()

	var deptStats []dto.DepartmentEmployeeStat
	for deptRows.Next() {
		var stat dto.DepartmentEmployeeStat
		err := deptRows.Scan(&stat.DepartmentID, &stat.DepartmentName, &stat.EmployeeCount, &stat.ActiveCount)
		if err != nil {
			continue
		}
		deptStats = append(deptStats, stat)
	}

	// Get job level stats - filtered by company
	levelRows, err := s.pool.Query(ctx, `
		SELECT 
			COALESCE(job_level, 'Unspecified') as level,
			COUNT(*) as count
		FROM employees
		WHERE company_id = $1 AND employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION')
		GROUP BY job_level
		ORDER BY count DESC
	`, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query job level stats: %w", err)
	}
	defer levelRows.Close()

	var levelStats []dto.JobLevelStat
	for levelRows.Next() {
		var stat dto.JobLevelStat
		err := levelRows.Scan(&stat.Level, &stat.EmployeeCount)
		if err != nil {
			continue
		}
		levelStats = append(levelStats, stat)
	}

	// Get recent hires (last 5) - filtered by company
	recentRows, err := s.pool.Query(ctx, `
		SELECT 
			e.id,
			e.full_name,
			e.position,
			e.join_date
		FROM employees e
		WHERE e.company_id = $1
		ORDER BY e.join_date DESC
		LIMIT 5
	`, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent hires: %w", err)
	}
	defer recentRows.Close()

	var recentHires []dto.RecentHire
	for recentRows.Next() {
		var hire dto.RecentHire
		var joinDate time.Time
		err := recentRows.Scan(&hire.EmployeeID, &hire.EmployeeName, &hire.Position, &joinDate)
		if err != nil {
			continue
		}
		hire.JoinDate = joinDate.Format("2006-01-02")
		recentHires = append(recentHires, hire)
	}

	return &dto.EmployeeStats{
		TotalCount:      totalCount,
		StatusBreakdown: statusBreakdown,
		DepartmentStats: deptStats,
		JobLevelStats:   levelStats,
		RecentHires:     recentHires,
	}, nil
}

// GetRecentActivities returns recent system activities
func (s *dashboardService) GetRecentActivities(ctx context.Context, companyID string, limit int) (*dto.RecentActivitiesResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// Query filtered by company
	query := `
		SELECT 
			id,
			created_at,
			action,
			resource_type,
			resource_id,
			user_id,
			COALESCE(user_name, 'Unknown'),
			action || ' ' || resource_type || ' (ID: ' || COALESCE(resource_id, 'N/A') || ')'
		FROM audit_logs
		WHERE company_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, query, companyID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent activities: %w", err)
	}
	defer rows.Close()

	var activities []dto.RecentActivity
	for rows.Next() {
		var activity dto.RecentActivity
		var timestamp time.Time
		err := rows.Scan(
			&activity.ID,
			&timestamp,
			&activity.Action,
			&activity.ResourceType,
			&activity.ResourceID,
			&activity.UserID,
			&activity.UserName,
			&activity.Description,
		)
		if err != nil {
			continue
		}
		activity.Timestamp = timestamp.Format("2006-01-02 15:04:05")
		activities = append(activities, activity)
	}

	// Get total count - filtered by company
	var total int
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_logs WHERE company_id = $1`, companyID).Scan(&total)
	if err != nil {
		total = len(activities)
	}

	return &dto.RecentActivitiesResponse{
		Activities: activities,
		Total:      total,
	}, nil
}
