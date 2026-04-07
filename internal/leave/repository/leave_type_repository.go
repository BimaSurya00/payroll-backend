package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/leave/entity"
)

var (
	ErrLeaveTypeNotFound = errors.New("leave type not found")
)

type LeaveTypeRepository interface {
	Create(ctx context.Context, leaveType *entity.LeaveType) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.LeaveType, error)
	FindByCode(ctx context.Context, code string) (*entity.LeaveType, error)
	FindAll(ctx context.Context) ([]entity.LeaveType, error)
	FindActive(ctx context.Context) ([]entity.LeaveType, error)
	Update(ctx context.Context, leaveType *entity.LeaveType) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.LeaveType, error)
}

type leaveTypeRepository struct {
	pool *pgxpool.Pool
}

func NewLeaveTypeRepository(pool *pgxpool.Pool) LeaveTypeRepository {
	return &leaveTypeRepository{pool: pool}
}

func (r *leaveTypeRepository) Create(ctx context.Context, leaveType *entity.LeaveType) error {
	query := `
		INSERT INTO leave_types (id, name, code, description, max_days_per_year, default_days, is_paid, requires_approval, is_active, color)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(ctx, query,
		leaveType.ID,
		leaveType.Name,
		leaveType.Code,
		leaveType.Description,
		leaveType.MaxDaysPerYear,
		leaveType.DefaultDays,
		leaveType.IsPaid,
		leaveType.RequiresApproval,
		leaveType.IsActive,
		leaveType.Color,
	)

	return err
}

func (r *leaveTypeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.LeaveType, error) {
	query := `SELECT id, name, code, description, max_days_per_year, default_days, is_paid, requires_approval, is_active, color, created_at, updated_at, company_id FROM leave_types WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var leaveType entity.LeaveType
	err := row.Scan(
		&leaveType.ID,
		&leaveType.Name,
		&leaveType.Code,
		&leaveType.Description,
		&leaveType.MaxDaysPerYear,
		&leaveType.DefaultDays,
		&leaveType.IsPaid,
		&leaveType.RequiresApproval,
		&leaveType.IsActive,
		&leaveType.Color,
		&leaveType.CreatedAt,
		&leaveType.UpdatedAt,
		&leaveType.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLeaveTypeNotFound
		}
		return nil, err
	}

	return &leaveType, nil
}

func (r *leaveTypeRepository) FindByCode(ctx context.Context, code string) (*entity.LeaveType, error) {
	query := `SELECT id, name, code, description, max_days_per_year, default_days, is_paid, requires_approval, is_active, color, created_at, updated_at, company_id FROM leave_types WHERE code = $1 AND is_active = true`

	row := r.pool.QueryRow(ctx, query, code)

	var leaveType entity.LeaveType
	err := row.Scan(
		&leaveType.ID,
		&leaveType.Name,
		&leaveType.Code,
		&leaveType.Description,
		&leaveType.MaxDaysPerYear,
		&leaveType.DefaultDays,
		&leaveType.IsPaid,
		&leaveType.RequiresApproval,
		&leaveType.IsActive,
		&leaveType.Color,
		&leaveType.CreatedAt,
		&leaveType.UpdatedAt,
		&leaveType.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLeaveTypeNotFound
		}
		return nil, err
	}

	return &leaveType, nil
}

func (r *leaveTypeRepository) FindAll(ctx context.Context) ([]entity.LeaveType, error) {
	query := `SELECT id, name, code, description, max_days_per_year, default_days, is_paid, requires_approval, is_active, color, created_at, updated_at, company_id FROM leave_types ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaveTypes []entity.LeaveType
	for rows.Next() {
		var leaveType entity.LeaveType
		err := rows.Scan(
			&leaveType.ID,
			&leaveType.Name,
			&leaveType.Code,
			&leaveType.Description,
			&leaveType.MaxDaysPerYear,
			&leaveType.DefaultDays,
			&leaveType.IsPaid,
			&leaveType.RequiresApproval,
			&leaveType.IsActive,
			&leaveType.Color,
			&leaveType.CreatedAt,
			&leaveType.UpdatedAt,
			&leaveType.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		leaveTypes = append(leaveTypes, leaveType)
	}

	return leaveTypes, nil
}

func (r *leaveTypeRepository) FindActive(ctx context.Context) ([]entity.LeaveType, error) {
	query := `SELECT id, name, code, description, max_days_per_year, default_days, is_paid, requires_approval, is_active, color, created_at, updated_at, company_id FROM leave_types WHERE is_active = true ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaveTypes []entity.LeaveType
	for rows.Next() {
		var leaveType entity.LeaveType
		err := rows.Scan(
			&leaveType.ID,
			&leaveType.Name,
			&leaveType.Code,
			&leaveType.Description,
			&leaveType.MaxDaysPerYear,
			&leaveType.DefaultDays,
			&leaveType.IsPaid,
			&leaveType.RequiresApproval,
			&leaveType.IsActive,
			&leaveType.Color,
			&leaveType.CreatedAt,
			&leaveType.UpdatedAt,
			&leaveType.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		leaveTypes = append(leaveTypes, leaveType)
	}

	return leaveTypes, nil
}

func (r *leaveTypeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.LeaveType, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `SELECT id, name, code, description, max_days_per_year, default_days, is_paid, requires_approval, is_active, color, created_at, updated_at, company_id
              FROM leave_types WHERE id = ANY($1)`
	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*entity.LeaveType
	for rows.Next() {
		lt := &entity.LeaveType{}
		err := rows.Scan(&lt.ID, &lt.Name, &lt.Code, &lt.Description,
			&lt.MaxDaysPerYear, &lt.DefaultDays, &lt.IsPaid, &lt.RequiresApproval, &lt.IsActive,
			&lt.Color, &lt.CreatedAt, &lt.UpdatedAt, &lt.CompanyID)
		if err != nil {
			return nil, err
		}
		types = append(types, lt)
	}
	return types, nil
}

func (r *leaveTypeRepository) Update(ctx context.Context, leaveType *entity.LeaveType) error {
	query := `
		UPDATE leave_types
		SET name = $2, description = $3, max_days_per_year = $4,
			is_paid = $5, requires_approval = $6, is_active = $7, updated_at = $8,
			default_days = $9, color = $10
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		leaveType.ID,
		leaveType.Name,
		leaveType.Description,
		leaveType.MaxDaysPerYear,
		leaveType.IsPaid,
		leaveType.RequiresApproval,
		leaveType.IsActive,
		time.Now(),
		leaveType.DefaultDays,
		leaveType.Color,
	)

	return err
}

func (r *leaveTypeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM leave_types WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrLeaveTypeNotFound
	}

	return nil
}
