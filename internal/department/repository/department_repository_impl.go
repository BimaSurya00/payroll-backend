package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/department/entity"
)

var (
	ErrDepartmentNotFound = errors.New("department not found")
)

type departmentRepository struct {
	pool *pgxpool.Pool
}

func NewDepartmentRepository(pool *pgxpool.Pool) DepartmentRepository {
	return &departmentRepository{pool: pool}
}

func (r *departmentRepository) Create(ctx context.Context, department *entity.Department) error {
	query := `
		INSERT INTO departments (id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		department.ID,
		department.CompanyID,
		department.Name,
		department.Code,
		department.Description,
		department.HeadEmployeeID,
		department.IsActive,
		department.CreatedAt,
		department.UpdatedAt,
	)

	return err
}

func (r *departmentRepository) FindByID(ctx context.Context, id string) (*entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var department entity.Department
	err := row.Scan(
		&department.ID,
		&department.CompanyID,
		&department.Name,
		&department.Code,
		&department.Description,
		&department.HeadEmployeeID,
		&department.IsActive,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDepartmentNotFound
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) FindByCode(ctx context.Context, code string) (*entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE code = $1`

	row := r.pool.QueryRow(ctx, query, code)

	var department entity.Department
	err := row.Scan(
		&department.ID,
		&department.CompanyID,
		&department.Name,
		&department.Code,
		&department.Description,
		&department.HeadEmployeeID,
		&department.IsActive,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDepartmentNotFound
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) FindAll(ctx context.Context) ([]entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []entity.Department
	for rows.Next() {
		var department entity.Department
		err := rows.Scan(
			&department.ID,
			&department.CompanyID,
			&department.Name,
			&department.Code,
			&department.Description,
			&department.HeadEmployeeID,
			&department.IsActive,
			&department.CreatedAt,
			&department.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *departmentRepository) FindActive(ctx context.Context) ([]entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE is_active = true ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []entity.Department
	for rows.Next() {
		var department entity.Department
		err := rows.Scan(
			&department.ID,
			&department.CompanyID,
			&department.Name,
			&department.Code,
			&department.Description,
			&department.HeadEmployeeID,
			&department.IsActive,
			&department.CreatedAt,
			&department.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *departmentRepository) Update(ctx context.Context, department *entity.Department) error {
	query := `
		UPDATE departments
		SET name = $2, code = $3, description = $4, head_employee_id = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		department.ID,
		department.Name,
		department.Code,
		department.Description,
		department.HeadEmployeeID,
		department.IsActive,
		time.Now(),
	)

	return err
}

func (r *departmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM departments WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDepartmentNotFound
	}

	return nil
}

func (r *departmentRepository) FindByIDAndCompany(ctx context.Context, id string, companyID string) (*entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE id = $1 AND company_id = $2`

	row := r.pool.QueryRow(ctx, query, id, companyID)

	var department entity.Department
	err := row.Scan(
		&department.ID,
		&department.CompanyID,
		&department.Name,
		&department.Code,
		&department.Description,
		&department.HeadEmployeeID,
		&department.IsActive,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDepartmentNotFound
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) FindByCodeAndCompany(ctx context.Context, code string, companyID string) (*entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE code = $1 AND company_id = $2`

	row := r.pool.QueryRow(ctx, query, code, companyID)

	var department entity.Department
	err := row.Scan(
		&department.ID,
		&department.CompanyID,
		&department.Name,
		&department.Code,
		&department.Description,
		&department.HeadEmployeeID,
		&department.IsActive,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDepartmentNotFound
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) FindAllByCompany(ctx context.Context, companyID string) ([]entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE company_id = $1 ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []entity.Department
	for rows.Next() {
		var department entity.Department
		err := rows.Scan(
			&department.ID,
			&department.CompanyID,
			&department.Name,
			&department.Code,
			&department.Description,
			&department.HeadEmployeeID,
			&department.IsActive,
			&department.CreatedAt,
			&department.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *departmentRepository) FindActiveByCompany(ctx context.Context, companyID string) ([]entity.Department, error) {
	query := `SELECT id, company_id, name, code, description, head_employee_id, is_active, created_at, updated_at FROM departments WHERE company_id = $1 AND is_active = true ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []entity.Department
	for rows.Next() {
		var department entity.Department
		err := rows.Scan(
			&department.ID,
			&department.CompanyID,
			&department.Name,
			&department.Code,
			&department.Description,
			&department.HeadEmployeeID,
			&department.IsActive,
			&department.CreatedAt,
			&department.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *departmentRepository) DeleteByCompany(ctx context.Context, id string, companyID string) error {
	query := `DELETE FROM departments WHERE id = $1 AND company_id = $2`

	result, err := r.pool.Exec(ctx, query, id, companyID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDepartmentNotFound
	}

	return nil
}
