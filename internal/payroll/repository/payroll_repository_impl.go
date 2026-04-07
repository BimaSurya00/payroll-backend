package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/payroll/entity"
)

var (
	ErrPayrollNotFound = errors.New("payroll not found")
)

type payrollRepository struct {
	pool *pgxpool.Pool
}

func NewPayrollRepository(pool *pgxpool.Pool) PayrollRepository {
	return &payrollRepository{pool: pool}
}

func (r *payrollRepository) CreateWithItems(ctx context.Context, payroll *entity.Payroll, items []*entity.PayrollItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert payroll
	query := `
		INSERT INTO payrolls (id, employee_id, period_start, period_end,
				     base_salary, total_allowance, total_deduction,
				     net_salary, status, generated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(ctx, query,
		payroll.ID, payroll.EmployeeID, payroll.PeriodStart, payroll.PeriodEnd,
		payroll.BaseSalary, payroll.TotalAllowance, payroll.TotalDeduction,
		payroll.NetSalary, payroll.Status, payroll.GeneratedAt,
		payroll.CreatedAt, payroll.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert payroll: %w", err)
	}

	// Insert payroll items
	for _, item := range items {
		itemQuery := `
			INSERT INTO payroll_items (id, payroll_id, name, amount, type, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.Exec(ctx, itemQuery, item.ID, item.PayrollID, item.Name, item.Amount, item.Type, item.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert payroll item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *payrollRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Payroll, error) {
	query := `
		SELECT id, employee_id, period_start, period_end,
		       base_salary, total_allowance, total_deduction,
		       net_salary, status, generated_at, created_at, updated_at
		FROM payrolls
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var payroll entity.Payroll
	err := row.Scan(
		&payroll.ID,
		&payroll.EmployeeID,
		&payroll.PeriodStart,
		&payroll.PeriodEnd,
		&payroll.BaseSalary,
		&payroll.TotalAllowance,
		&payroll.TotalDeduction,
		&payroll.NetSalary,
		&payroll.Status,
		&payroll.GeneratedAt,
		&payroll.CreatedAt,
		&payroll.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPayrollNotFound
		}
		return nil, fmt.Errorf("failed to find payroll: %w", err)
	}

	return &payroll, nil
}

func (r *payrollRepository) FindByIDWithItems(ctx context.Context, id uuid.UUID) (*entity.PayrollWithItems, error) {
	// Fetch payroll
	payroll, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Fetch items
	query := `
		SELECT id, payroll_id, name, amount, type, created_at
		FROM payroll_items
		WHERE payroll_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payroll items: %w", err)
	}
	defer rows.Close()

	var items []*entity.PayrollItem
	for rows.Next() {
		var item entity.PayrollItem
		err := rows.Scan(&item.ID, &item.PayrollID, &item.Name, &item.Amount, &item.Type, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll item: %w", err)
		}
		items = append(items, &item)
	}

	return &entity.PayrollWithItems{
		Payroll: payroll,
		Items:   items,
	}, nil
}

func (r *payrollRepository) FindAll(ctx context.Context, skip, limit int64) ([]*entity.Payroll, error) {
	query := `
		SELECT p.id, p.employee_id, p.period_start, p.period_end,
		       p.base_salary, p.total_allowance, p.total_deduction,
		       p.net_salary, p.status, p.generated_at, p.created_at, p.updated_at
		FROM payrolls p
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find payrolls: %w", err)
	}
	defer rows.Close()

	var payrolls []*entity.Payroll
	for rows.Next() {
		var payroll entity.Payroll
		err := rows.Scan(
			&payroll.ID,
			&payroll.EmployeeID,
			&payroll.PeriodStart,
			&payroll.PeriodEnd,
			&payroll.BaseSalary,
			&payroll.TotalAllowance,
			&payroll.TotalDeduction,
			&payroll.NetSalary,
			&payroll.Status,
			&payroll.GeneratedAt,
			&payroll.CreatedAt,
			&payroll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll: %w", err)
		}
		payrolls = append(payrolls, &payroll)
	}

	return payrolls, nil
}

func (r *payrollRepository) FindByPeriod(ctx context.Context, periodStart, periodEnd string) ([]*entity.Payroll, error) {
	query := `
		SELECT id, employee_id, period_start, period_end,
		       base_salary, total_allowance, total_deduction,
		       net_salary, status, generated_at, created_at, updated_at
		FROM payrolls
		WHERE period_start = $1 AND period_end = $2
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to find payrolls by period: %w", err)
	}
	defer rows.Close()

	var payrolls []*entity.Payroll
	for rows.Next() {
		var payroll entity.Payroll
		err := rows.Scan(
			&payroll.ID,
			&payroll.EmployeeID,
			&payroll.PeriodStart,
			&payroll.PeriodEnd,
			&payroll.BaseSalary,
			&payroll.TotalAllowance,
			&payroll.TotalDeduction,
			&payroll.NetSalary,
			&payroll.Status,
			&payroll.GeneratedAt,
			&payroll.CreatedAt,
			&payroll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll: %w", err)
		}
		payrolls = append(payrolls, &payroll)
	}

	return payrolls, nil
}

func (r *payrollRepository) FindByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, periodStart, periodEnd string) (*entity.Payroll, error) {
	query := `
		SELECT id, employee_id, period_start, period_end,
		       base_salary, total_allowance, total_deduction,
		       net_salary, status, generated_at, created_at, updated_at
		FROM payrolls
		WHERE employee_id = $1 AND period_start = $2 AND period_end = $3
	`

	row := r.pool.QueryRow(ctx, query, employeeID, periodStart, periodEnd)

	var payroll entity.Payroll
	err := row.Scan(
		&payroll.ID,
		&payroll.EmployeeID,
		&payroll.PeriodStart,
		&payroll.PeriodEnd,
		&payroll.BaseSalary,
		&payroll.TotalAllowance,
		&payroll.TotalDeduction,
		&payroll.NetSalary,
		&payroll.Status,
		&payroll.GeneratedAt,
		&payroll.CreatedAt,
		&payroll.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPayrollNotFound
		}
		return nil, fmt.Errorf("failed to find payroll: %w", err)
	}

	return &payroll, nil
}

func (r *payrollRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE payrolls
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	result, err := r.pool.Exec(ctx, query, status, now, id)
	if err != nil {
		return fmt.Errorf("failed to update payroll status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPayrollNotFound
	}

	return nil
}

func (r *payrollRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM payrolls`

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count payrolls: %w", err)
	}

	return count, nil
}

func (r *payrollRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payrolls WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payroll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPayrollNotFound
	}

	return nil
}

func (r *payrollRepository) FindByEmployeeID(ctx context.Context, employeeID string, skip, limit int64) ([]*entity.Payroll, error) {
	query := `SELECT id, employee_id, period_start, period_end, base_salary,
              total_allowance, total_deduction, net_salary, status, generated_at, created_at, updated_at
              FROM payrolls WHERE employee_id = $1
              ORDER BY period_start DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, employeeID, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payrolls []*entity.Payroll
	for rows.Next() {
		p := &entity.Payroll{}
		err := rows.Scan(&p.ID, &p.EmployeeID, &p.PeriodStart, &p.PeriodEnd,
			&p.BaseSalary, &p.TotalAllowance, &p.TotalDeduction, &p.NetSalary,
			&p.Status, &p.GeneratedAt, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payrolls = append(payrolls, p)
	}
	return payrolls, nil
}

func (r *payrollRepository) CountByEmployeeID(ctx context.Context, employeeID string) (int64, error) {
	query := `SELECT COUNT(*) FROM payrolls WHERE employee_id = $1`
	var count int64
	err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
	return count, err
}

func (r *payrollRepository) FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.Payroll, error) {
	query := `
		SELECT p.id, p.employee_id, p.period_start, p.period_end,
		       p.base_salary, p.total_allowance, p.total_deduction,
		       p.net_salary, p.status, p.generated_at, p.created_at, p.updated_at
		FROM payrolls p
		JOIN employees e ON p.employee_id = e.id
		WHERE p.id = $1 AND e.company_id = $2
	`

	row := r.pool.QueryRow(ctx, query, id, companyID)

	var payroll entity.Payroll
	err := row.Scan(
		&payroll.ID,
		&payroll.EmployeeID,
		&payroll.PeriodStart,
		&payroll.PeriodEnd,
		&payroll.BaseSalary,
		&payroll.TotalAllowance,
		&payroll.TotalDeduction,
		&payroll.NetSalary,
		&payroll.Status,
		&payroll.GeneratedAt,
		&payroll.CreatedAt,
		&payroll.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPayrollNotFound
		}
		return nil, fmt.Errorf("failed to find payroll: %w", err)
	}

	return &payroll, nil
}

func (r *payrollRepository) FindByIDWithItemsAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.PayrollWithItems, error) {
	payroll, err := r.FindByIDAndCompany(ctx, id, companyID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, payroll_id, name, amount, type, created_at
		FROM payroll_items
		WHERE payroll_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payroll items: %w", err)
	}
	defer rows.Close()

	var items []*entity.PayrollItem
	for rows.Next() {
		var item entity.PayrollItem
		err := rows.Scan(&item.ID, &item.PayrollID, &item.Name, &item.Amount, &item.Type, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll item: %w", err)
		}
		items = append(items, &item)
	}

	return &entity.PayrollWithItems{
		Payroll: payroll,
		Items:   items,
	}, nil
}

func (r *payrollRepository) FindAllByCompany(ctx context.Context, companyID string, skip, limit int64) ([]*entity.Payroll, error) {
	query := `
		SELECT p.id, p.employee_id, p.period_start, p.period_end,
		       p.base_salary, p.total_allowance, p.total_deduction,
		       p.net_salary, p.status, p.generated_at, p.created_at, p.updated_at
		FROM payrolls p
		JOIN employees e ON p.employee_id = e.id
		WHERE e.company_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, companyID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find payrolls: %w", err)
	}
	defer rows.Close()

	var payrolls []*entity.Payroll
	for rows.Next() {
		var payroll entity.Payroll
		err := rows.Scan(
			&payroll.ID,
			&payroll.EmployeeID,
			&payroll.PeriodStart,
			&payroll.PeriodEnd,
			&payroll.BaseSalary,
			&payroll.TotalAllowance,
			&payroll.TotalDeduction,
			&payroll.NetSalary,
			&payroll.Status,
			&payroll.GeneratedAt,
			&payroll.CreatedAt,
			&payroll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll: %w", err)
		}
		payrolls = append(payrolls, &payroll)
	}

	return payrolls, nil
}

func (r *payrollRepository) FindByPeriodAndCompany(ctx context.Context, companyID string, periodStart, periodEnd string) ([]*entity.Payroll, error) {
	query := `
		SELECT p.id, p.employee_id, p.period_start, p.period_end,
		       p.base_salary, p.total_allowance, p.total_deduction,
		       p.net_salary, p.status, p.generated_at, p.created_at, p.updated_at
		FROM payrolls p
		JOIN employees e ON p.employee_id = e.id
		WHERE e.company_id = $1 AND p.period_start = $2 AND p.period_end = $3
		ORDER BY p.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, companyID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to find payrolls by period: %w", err)
	}
	defer rows.Close()

	var payrolls []*entity.Payroll
	for rows.Next() {
		var payroll entity.Payroll
		err := rows.Scan(
			&payroll.ID,
			&payroll.EmployeeID,
			&payroll.PeriodStart,
			&payroll.PeriodEnd,
			&payroll.BaseSalary,
			&payroll.TotalAllowance,
			&payroll.TotalDeduction,
			&payroll.NetSalary,
			&payroll.Status,
			&payroll.GeneratedAt,
			&payroll.CreatedAt,
			&payroll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payroll: %w", err)
		}
		payrolls = append(payrolls, &payroll)
	}

	return payrolls, nil
}

func (r *payrollRepository) CountByCompany(ctx context.Context, companyID string) (int64, error) {
	query := `
		SELECT COUNT(*) FROM payrolls p
		JOIN employees e ON p.employee_id = e.id
		WHERE e.company_id = $1
	`
	var count int64
	err := r.pool.QueryRow(ctx, query, companyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count payrolls: %w", err)
	}
	return count, nil
}

func (r *payrollRepository) DeleteByCompany(ctx context.Context, id uuid.UUID, companyID string) error {
	query := `
		DELETE FROM payrolls p
		USING employees e
		WHERE p.id = $1 AND p.employee_id = e.id AND e.company_id = $2
	`

	result, err := r.pool.Exec(ctx, query, id, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete payroll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrPayrollNotFound
	}

	return nil
}
