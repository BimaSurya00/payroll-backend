package repository

import (
	"context"
	"errors"
	"time"

	"hris/internal/payroll/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPayrollConfigNotFound = errors.New("payroll config not found")
)

type payrollConfigRepository struct {
	pool *pgxpool.Pool
}

func NewPayrollConfigRepository(pool *pgxpool.Pool) PayrollConfigRepository {
	return &payrollConfigRepository{pool: pool}
}

func (r *payrollConfigRepository) Create(ctx context.Context, config *entity.PayrollConfig) error {
	query := `
		INSERT INTO payroll_configs (id, name, code, type, amount, calculation_type, is_active, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query,
		config.ID,
		config.Name,
		config.Code,
		config.Type,
		config.Amount,
		config.CalculationType,
		config.IsActive,
		config.Description,
	)

	return err
}

func (r *payrollConfigRepository) FindByID(ctx context.Context, id string) (*entity.PayrollConfig, error) {
	query := `SELECT * FROM payroll_configs WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var config entity.PayrollConfig
	err := row.Scan(
		&config.ID,
		&config.Name,
		&config.Code,
		&config.Type,
		&config.Amount,
		&config.CalculationType,
		&config.IsActive,
		&config.Description,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPayrollConfigNotFound
		}
		return nil, err
	}

	return &config, nil
}

func (r *payrollConfigRepository) FindByCode(ctx context.Context, code string) (*entity.PayrollConfig, error) {
	query := `SELECT * FROM payroll_configs WHERE code = $1`

	row := r.pool.QueryRow(ctx, query, code)

	var config entity.PayrollConfig
	err := row.Scan(
		&config.ID,
		&config.Name,
		&config.Code,
		&config.Type,
		&config.Amount,
		&config.CalculationType,
		&config.IsActive,
		&config.Description,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPayrollConfigNotFound
		}
		return nil, err
	}

	return &config, nil
}

func (r *payrollConfigRepository) FindAll(ctx context.Context) ([]*entity.PayrollConfig, error) {
	query := `SELECT * FROM payroll_configs ORDER BY type, name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*entity.PayrollConfig
	for rows.Next() {
		config := &entity.PayrollConfig{}
		err := rows.Scan(
			&config.ID,
			&config.Name,
			&config.Code,
			&config.Type,
			&config.Amount,
			&config.CalculationType,
			&config.IsActive,
			&config.Description,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (r *payrollConfigRepository) FindActiveByType(ctx context.Context, configType string) ([]*entity.PayrollConfig, error) {
	query := `SELECT * FROM payroll_configs WHERE is_active = true`
	if configType != "" {
		query += ` AND type = $1`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, configType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*entity.PayrollConfig
	for rows.Next() {
		config := &entity.PayrollConfig{}
		err := rows.Scan(
			&config.ID,
			&config.Name,
			&config.Code,
			&config.Type,
			&config.Amount,
			&config.CalculationType,
			&config.IsActive,
			&config.Description,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (r *payrollConfigRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	query := `
		UPDATE payroll_configs
		SET name = $2, code = $3, type = $4, amount = $5, calculation_type = $6, is_active = $7, description = $8, updated_at = $9
		WHERE id = $1
	`

	updates["updated_at"] = time.Now()

	_, err := r.pool.Exec(ctx, query,
		id,
		updates["name"],
		updates["code"],
		updates["type"],
		updates["amount"],
		updates["calculation_type"],
		updates["is_active"],
		updates["description"],
	)
	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrPayrollConfigNotFound
	}

	return nil
}

func (r *payrollConfigRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM payroll_configs WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrPayrollConfigNotFound
	}

	return nil
}
