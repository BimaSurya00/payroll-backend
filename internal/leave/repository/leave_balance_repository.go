package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/leave/entity"
)

var (
	ErrLeaveBalanceNotFound = errors.New("leave balance not found")
)

type LeaveBalanceRepository interface {
	Create(ctx context.Context, balance *entity.LeaveBalance) error
	FindByEmployeeAndYear(ctx context.Context, employeeID uuid.UUID, year int) ([]entity.LeaveBalance, error)
	FindByEmployeeTypeAndYear(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*entity.LeaveBalance, error)
	AddToPending(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, days int) error
	MoveFromPendingToUsed(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, days int) error
	ReturnFromPending(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, days int) error
}

type leaveBalanceRepository struct {
	pool *pgxpool.Pool
}

func NewLeaveBalanceRepository(pool *pgxpool.Pool) LeaveBalanceRepository {
	return &leaveBalanceRepository{pool: pool}
}

func (r *leaveBalanceRepository) Create(ctx context.Context, balance *entity.LeaveBalance) error {
	query := `
		INSERT INTO leave_balances (id, employee_id, leave_type_id, year, balance, used, pending)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (employee_id, leave_type_id, year)
		DO UPDATE SET balance = EXCLUDED.balance, updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query,
		balance.ID,
		balance.EmployeeID,
		balance.LeaveTypeID,
		balance.Year,
		balance.Balance,
		balance.Used,
		balance.Pending,
	)

	return err
}

func (r *leaveBalanceRepository) FindByEmployeeAndYear(ctx context.Context, employeeID uuid.UUID, year int) ([]entity.LeaveBalance, error) {
	query := `
		SELECT * FROM leave_balances
		WHERE employee_id = $1 AND year = $2
		ORDER BY leave_type_id ASC
	`

	rows, err := r.pool.Query(ctx, query, employeeID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []entity.LeaveBalance
	for rows.Next() {
		var balance entity.LeaveBalance
		err := rows.Scan(
			&balance.ID,
			&balance.EmployeeID,
			&balance.LeaveTypeID,
			&balance.Year,
			&balance.Balance,
			&balance.Used,
			&balance.Pending,
			&balance.CreatedAt,
			&balance.UpdatedAt,
			&balance.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

func (r *leaveBalanceRepository) FindByEmployeeTypeAndYear(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*entity.LeaveBalance, error) {
	query := `
		SELECT * FROM leave_balances
		WHERE employee_id = $1 AND leave_type_id = $2 AND year = $3
	`

	row := r.pool.QueryRow(ctx, query, employeeID, leaveTypeID, year)

	var balance entity.LeaveBalance
	err := row.Scan(
		&balance.ID,
		&balance.EmployeeID,
		&balance.LeaveTypeID,
		&balance.Year,
		&balance.Balance,
		&balance.Used,
		&balance.Pending,
		&balance.CreatedAt,
		&balance.UpdatedAt,
		&balance.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLeaveBalanceNotFound
		}
		return nil, err
	}

	return &balance, nil
}

func (r *leaveBalanceRepository) AddToPending(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, days int) error {
	query := `
		INSERT INTO leave_balances (id, employee_id, leave_type_id, year, balance, used, pending)
		VALUES (gen_random_uuid(), $1, $2, $3,
			(SELECT max_days_per_year FROM leave_types WHERE id = $2), 0, $4)
		ON CONFLICT (employee_id, leave_type_id, year)
		DO UPDATE SET pending = leave_balances.pending + $4, updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query, employeeID, leaveTypeID, year, days)
	return err
}

func (r *leaveBalanceRepository) MoveFromPendingToUsed(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, days int) error {
	query := `
		UPDATE leave_balances
		SET pending = pending - $1,
			used = used + $1,
			updated_at = NOW()
		WHERE employee_id = $2 AND leave_type_id = $3 AND year = $4
	`

	result, err := r.pool.Exec(ctx, query, days, employeeID, leaveTypeID, year)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrLeaveBalanceNotFound
	}

	return nil
}

func (r *leaveBalanceRepository) ReturnFromPending(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int, days int) error {
	query := `
		UPDATE leave_balances
		SET pending = pending - $1,
			updated_at = NOW()
		WHERE employee_id = $2 AND leave_type_id = $3 AND year = $4
	`

	result, err := r.pool.Exec(ctx, query, days, employeeID, leaveTypeID, year)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrLeaveBalanceNotFound
	}

	return nil
}
