package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/overtime/entity"
)

var (
	ErrOvertimePolicyNotFound = errors.New("overtime policy not found")
)

type OvertimePolicyRepository interface {
	Create(ctx context.Context, policy *entity.OvertimePolicy) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.OvertimePolicy, error)
	FindAll(ctx context.Context) ([]entity.OvertimePolicy, error)
	FindActive(ctx context.Context) ([]entity.OvertimePolicy, error)
	Update(ctx context.Context, policy *entity.OvertimePolicy) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type overtimePolicyRepository struct {
	pool *pgxpool.Pool
}

func NewOvertimePolicyRepository(pool *pgxpool.Pool) OvertimePolicyRepository {
	return &overtimePolicyRepository{pool: pool}
}

func (r *overtimePolicyRepository) Create(ctx context.Context, policy *entity.OvertimePolicy) error {
	query := `
		INSERT INTO overtime_policies (id, name, description, rate_type, rate_multiplier, fixed_amount,
			min_overtime_minutes, max_overtime_hours_per_day, max_overtime_hours_per_month, requires_approval, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		policy.ID,
		policy.Name,
		policy.Description,
		policy.RateType,
		policy.RateMultiplier,
		policy.FixedAmount,
		policy.MinOvertimeMinutes,
		policy.MaxOvertimeHoursPerDay,
		policy.MaxOvertimeHoursPerMonth,
		policy.RequiresApproval,
		policy.IsActive,
		policy.CreatedAt,
		policy.UpdatedAt,
	)

	return err
}

func (r *overtimePolicyRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.OvertimePolicy, error) {
	query := `SELECT id, name, description, rate_type, rate_multiplier, fixed_amount, min_overtime_minutes, max_overtime_hours_per_day, max_overtime_hours_per_month, requires_approval, is_active, created_at, updated_at, company_id FROM overtime_policies WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var policy entity.OvertimePolicy
	err := row.Scan(
		&policy.ID,
		&policy.Name,
		&policy.Description,
		&policy.RateType,
		&policy.RateMultiplier,
		&policy.FixedAmount,
		&policy.MinOvertimeMinutes,
		&policy.MaxOvertimeHoursPerDay,
		&policy.MaxOvertimeHoursPerMonth,
		&policy.RequiresApproval,
		&policy.IsActive,
		&policy.CreatedAt,
		&policy.UpdatedAt,
		&policy.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOvertimePolicyNotFound
		}
		return nil, err
	}

	return &policy, nil
}

func (r *overtimePolicyRepository) FindAll(ctx context.Context) ([]entity.OvertimePolicy, error) {
	query := `SELECT id, name, description, rate_type, rate_multiplier, fixed_amount, min_overtime_minutes, max_overtime_hours_per_day, max_overtime_hours_per_month, requires_approval, is_active, created_at, updated_at, company_id FROM overtime_policies ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []entity.OvertimePolicy
	for rows.Next() {
		var policy entity.OvertimePolicy
		err := rows.Scan(
			&policy.ID,
			&policy.Name,
			&policy.Description,
			&policy.RateType,
			&policy.RateMultiplier,
			&policy.FixedAmount,
			&policy.MinOvertimeMinutes,
			&policy.MaxOvertimeHoursPerDay,
			&policy.MaxOvertimeHoursPerMonth,
			&policy.RequiresApproval,
			&policy.IsActive,
			&policy.CreatedAt,
			&policy.UpdatedAt,
			&policy.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

func (r *overtimePolicyRepository) FindActive(ctx context.Context) ([]entity.OvertimePolicy, error) {
	query := `SELECT id, name, description, rate_type, rate_multiplier, fixed_amount, min_overtime_minutes, max_overtime_hours_per_day, max_overtime_hours_per_month, requires_approval, is_active, created_at, updated_at, company_id FROM overtime_policies WHERE is_active = true ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var policies []entity.OvertimePolicy
	for rows.Next() {
		var policy entity.OvertimePolicy
		err := rows.Scan(
			&policy.ID,
			&policy.Name,
			&policy.Description,
			&policy.RateType,
			&policy.RateMultiplier,
			&policy.FixedAmount,
			&policy.MinOvertimeMinutes,
			&policy.MaxOvertimeHoursPerDay,
			&policy.MaxOvertimeHoursPerMonth,
			&policy.RequiresApproval,
			&policy.IsActive,
			&policy.CreatedAt,
			&policy.UpdatedAt,
			&policy.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}

	return policies, nil
}

func (r *overtimePolicyRepository) Update(ctx context.Context, policy *entity.OvertimePolicy) error {
	query := `
		UPDATE overtime_policies
		SET name = $2, description = $3, rate_type = $4, rate_multiplier = $5, fixed_amount = $6,
			min_overtime_minutes = $7, max_overtime_hours_per_day = $8, max_overtime_hours_per_month = $9,
			requires_approval = $10, is_active = $11, updated_at = $12
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		policy.ID,
		policy.Name,
		policy.Description,
		policy.RateType,
		policy.RateMultiplier,
		policy.FixedAmount,
		policy.MinOvertimeMinutes,
		policy.MaxOvertimeHoursPerDay,
		policy.MaxOvertimeHoursPerMonth,
		policy.RequiresApproval,
		policy.IsActive,
		policy.UpdatedAt,
	)

	return err
}

func (r *overtimePolicyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM overtime_policies WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
