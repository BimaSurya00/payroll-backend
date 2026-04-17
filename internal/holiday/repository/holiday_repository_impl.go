package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/holiday/entity"
)

var (
	ErrHolidayNotFound = errors.New("holiday not found")
)

type holidayRepository struct {
	pool *pgxpool.Pool
}

func NewHolidayRepository(pool *pgxpool.Pool) HolidayRepository {
	return &holidayRepository{pool: pool}
}

func (r *holidayRepository) Create(ctx context.Context, holiday *entity.Holiday) error {
	query := `
		INSERT INTO holidays (id, name, date, type, is_recurring, year, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		holiday.ID,
		holiday.Name,
		holiday.Date,
		holiday.Type,
		holiday.IsRecurring,
		holiday.Year,
		holiday.Description,
	)

	return err
}

func (r *holidayRepository) FindByID(ctx context.Context, id string) (*entity.Holiday, error) {
	query := `SELECT * FROM holidays WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var holiday entity.Holiday
	err := row.Scan(
		&holiday.ID,
		&holiday.CompanyID,
		&holiday.Name,
		&holiday.Date,
		&holiday.Type,
		&holiday.IsRecurring,
		&holiday.Year,
		&holiday.Description,
		&holiday.CreatedAt,
		&holiday.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrHolidayNotFound
		}
		return nil, err
	}

	return &holiday, nil
}

func (r *holidayRepository) FindAll(ctx context.Context, year int) ([]*entity.Holiday, error) {
	query := `SELECT * FROM holidays WHERE year = $1 OR is_recurring = true ORDER BY date ASC`

	rows, err := r.pool.Query(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []*entity.Holiday
	for rows.Next() {
		var holiday entity.Holiday
		err := rows.Scan(
			&holiday.ID,
			&holiday.CompanyID,
			&holiday.Name,
			&holiday.Date,
			&holiday.Type,
			&holiday.IsRecurring,
			&holiday.Year,
			&holiday.Description,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, &holiday)
	}

	return holidays, nil
}

func (r *holidayRepository) FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error) {
	query := `SELECT * FROM holidays WHERE date >= $1 AND date <= $2 ORDER BY date ASC`

	rows, err := r.pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []*entity.Holiday
	for rows.Next() {
		var holiday entity.Holiday
		err := rows.Scan(
			&holiday.ID,
			&holiday.CompanyID,
			&holiday.Name,
			&holiday.Date,
			&holiday.Type,
			&holiday.IsRecurring,
			&holiday.Year,
			&holiday.Description,
			&holiday.CreatedAt,
			&holiday.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		holidays = append(holidays, &holiday)
	}

	return holidays, nil
}

func (r *holidayRepository) Update(ctx context.Context, holiday *entity.Holiday) error {
	query := `
		UPDATE holidays
		SET name = $2, date = $3, type = $4, is_recurring = $5, year = $6, description = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		holiday.ID,
		holiday.Name,
		holiday.Date,
		holiday.Type,
		holiday.IsRecurring,
		holiday.Year,
		holiday.Description,
		time.Now(),
	)

	return err
}

func (r *holidayRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM holidays WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrHolidayNotFound
	}

	return nil
}

func (r *holidayRepository) IsHoliday(ctx context.Context, date time.Time) (bool, error) {
	query := `SELECT COUNT(*) FROM holidays WHERE date = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, date).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
