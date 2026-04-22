package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/company/dto"
	"hris/internal/company/entity"
)

var ErrCompanyNotFound = errors.New("company not found")

type companyRepository struct {
	pool *pgxpool.Pool
}

func NewCompanyRepository(pool *pgxpool.Pool) CompanyRepository {
	return &companyRepository{pool: pool}
}

func (r *companyRepository) Create(ctx context.Context, company *entity.Company) error {
	query := `INSERT INTO companies (id, name, slug, is_active, plan, max_employees,
		office_lat, office_long, allowed_radius_meters, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query,
		company.ID, company.Name, company.Slug, company.IsActive,
		company.Plan, company.MaxEmployees, company.OfficeLat, company.OfficeLong,
		company.AllowedRadiusMeters, now, now,
	)
	return err
}

func (r *companyRepository) FindByID(ctx context.Context, id string) (*entity.Company, error) {
	query := `SELECT id, name, slug, is_active, plan, max_employees,
		office_lat, office_long, allowed_radius_meters, created_at, updated_at
		FROM companies WHERE id = $1`

	var c entity.Company
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.IsActive,
		&c.Plan, &c.MaxEmployees, &c.OfficeLat, &c.OfficeLong,
		&c.AllowedRadiusMeters, &c.CreatedAt, &c.UpdatedAt,
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
	query := `SELECT id, name, slug, is_active, plan, max_employees,
		office_lat, office_long, allowed_radius_meters, created_at, updated_at
		FROM companies WHERE slug = $1`

	var c entity.Company
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&c.ID, &c.Name, &c.Slug, &c.IsActive,
		&c.Plan, &c.MaxEmployees, &c.OfficeLat, &c.OfficeLong,
		&c.AllowedRadiusMeters, &c.CreatedAt, &c.UpdatedAt,
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
	var total int64
	countQuery := `SELECT COUNT(*) FROM companies`
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `SELECT id, name, slug, is_active, plan, max_employees,
		office_lat, office_long, allowed_radius_meters, created_at, updated_at
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
			&c.Plan, &c.MaxEmployees, &c.OfficeLat, &c.OfficeLong,
			&c.AllowedRadiusMeters, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		companies = append(companies, &c)
	}

	return companies, total, nil
}

func (r *companyRepository) Update(ctx context.Context, company *entity.Company) error {
	query := `UPDATE companies SET name = $1, slug = $2, is_active = $3, plan = $4,
		max_employees = $5, office_lat = $6, office_long = $7, allowed_radius_meters = $8,
		updated_at = $9 WHERE id = $10`

	_, err := r.pool.Exec(ctx, query,
		company.Name, company.Slug, company.IsActive,
		company.Plan, company.MaxEmployees, company.OfficeLat, company.OfficeLong,
		company.AllowedRadiusMeters, time.Now(), company.ID,
	)
	return err
}

func (r *companyRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *companyRepository) FindAllWithStats(ctx context.Context, page, perPage int) ([]*dto.CompanyListItem, int64, error) {
	var total int64
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM companies`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT c.id, c.name, c.slug, c.is_active, c.plan, c.max_employees, c.created_at,
			COALESCE((SELECT COUNT(*) FROM users WHERE company_id = c.id), 0),
			COALESCE((SELECT COUNT(*) FROM employees WHERE company_id = c.id AND deleted_at IS NULL), 0)
		FROM companies c
		ORDER BY c.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []*dto.CompanyListItem
	for rows.Next() {
		var item dto.CompanyListItem
		var createdAt time.Time
		if err := rows.Scan(
			&item.ID, &item.Name, &item.Slug, &item.IsActive,
			&item.Plan, &item.MaxEmployees, &createdAt,
			&item.UserCount, &item.EmployeeCount,
		); err != nil {
			return nil, 0, err
		}
		item.CreatedAt = createdAt.Format("2006-01-02")
		items = append(items, &item)
	}

	return items, total, nil
}

func (r *companyRepository) GetStats(ctx context.Context, companyID string) (*dto.CompanyStatsResponse, error) {
	var c entity.Company
	if err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, is_active, plan, max_employees,
			office_lat, office_long, allowed_radius_meters, created_at, updated_at
		FROM companies WHERE id = $1`, companyID,
	).Scan(
		&c.ID, &c.Name, &c.Slug, &c.IsActive,
		&c.Plan, &c.MaxEmployees, &c.OfficeLat, &c.OfficeLong,
		&c.AllowedRadiusMeters, &c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, err
	}

	stats := &dto.CompanyStatsResponse{
		CompanyID:   c.ID,
		CompanyName: c.Name,
	}

	queries := []struct {
		sql  string
		dest *int
	}{
		{`SELECT COUNT(*) FROM users WHERE company_id = $1`, &stats.UserCount},
		{`SELECT COUNT(*) FROM employees WHERE company_id = $1 AND deleted_at IS NULL`, &stats.EmployeeCount},
		{`SELECT COUNT(*) FROM departments WHERE company_id = $1`, &stats.DepartmentCount},
		{`SELECT COUNT(*) FROM schedules WHERE company_id = $1`, &stats.ScheduleCount},
	}

	for _, q := range queries {
		if err := r.pool.QueryRow(ctx, q.sql, companyID).Scan(q.dest); err != nil {
			*q.dest = 0
		}
	}

	return stats, nil
}
