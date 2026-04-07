package repository

import (
	"context"
	"errors"
	"time"

	"hris/internal/company/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCompanyNotFound = errors.New("company not found")

type companyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) CompanyRepository {
	return &companyRepository{pool: pool}
}

func (r *companyRepository) Create(ctx context.Context, company *entity.Company) error {
	query := `INSERT INTO companies (id, name, slug, is_active, plan, max_employees, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query,
		company.ID, company.Name, company.Slug, company.IsActive,
		company.Plan, company.MaxEmployees, now, now,
	)
	return err
}

func (r *companyRepository) FindByID(ctx context.Context, id string) (*entity.Company, error) {
	query := `SELECT id, name, slug, is_active, plan, max_employees, created_at, updated_at
		FROM companies WHERE id = $1`

	var c entity.Company
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.IsActive,
		&c.Plan, &c.MaxEmployees, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) FindBySlug(ctx context.Context, slug string) (*entity.Company, error) {
	query := `SELECT id, name, slug, is_active, plan, max_employees, created_at, updated_at
		FROM companies WHERE slug = $1`

	var c entity.Company
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&c.ID, &c.Name, &c.Slug, &c.IsActive,
		&c.Plan, &c.MaxEmployees, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *companyRepository) FindAll(ctx context.Context, page, perPage int) ([]*entity.Company, int64, error) {
	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM companies`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch page
	offset := (page - 1) * perPage
	query := `SELECT id, name, slug, is_active, plan, max_employees, created_at, updated_at
		FROM companies ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var companies []*entity.Company
	for rows.Next() {
		var c entity.Company
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.IsActive,
			&c.Plan, &c.MaxEmployees, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		companies = append(companies, &c)
	}

	return companies, total, nil
}

func (r *companyRepository) Update(ctx context.Context, company *entity.Company) error {
	query := `UPDATE companies SET name = $1, slug = $2, is_active = $3, plan = $4, 
		max_employees = $5, updated_at = $6 WHERE id = $7`

	_, err := r.pool.Exec(ctx, query,
		company.Name, company.Slug, company.IsActive,
		company.Plan, company.MaxEmployees, time.Now(), company.ID,
	)
	return err
}

func (r *companyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
