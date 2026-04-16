package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/schedule/entity"
)

var (
	ErrScheduleNotFound = errors.New("schedule not found")
)

type scheduleRepository struct {
	pool *pgxpool.Pool
}

func NewScheduleRepository(pool *pgxpool.Pool) ScheduleRepository {
	return &scheduleRepository{pool: pool}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *entity.Schedule) error {
	query := `
		INSERT INTO schedules (id, company_id, name, time_in, time_out, allowed_late_minutes,
		                      description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		schedule.ID,
		schedule.CompanyID,
		schedule.Name,
		schedule.TimeIn,
		schedule.TimeOut,
		schedule.AllowedLateMinutes,
		schedule.Description,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)

	return err
}

func (r *scheduleRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	query := `
		SELECT id, name, time_in, time_out, allowed_late_minutes,
		       COALESCE(description, ''), created_at, updated_at
		FROM schedules
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var schedule entity.Schedule
	err := row.Scan(
		&schedule.ID,
		&schedule.Name,
		&schedule.TimeIn,
		&schedule.TimeOut,
		&schedule.AllowedLateMinutes,
		&schedule.Description,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to find schedule: %w", err)
	}

	return &schedule, nil
}

func (r *scheduleRepository) FindAll(ctx context.Context, skip, limit int64) ([]*entity.Schedule, error) {
	query := `
		SELECT id, name, time_in, time_out, allowed_late_minutes,
		       COALESCE(description, ''), created_at, updated_at
		FROM schedules
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*entity.Schedule
	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.Name,
			&schedule.TimeIn,
			&schedule.TimeOut,
			&schedule.AllowedLateMinutes,
			&schedule.Description,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, nil
}

func (r *scheduleRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM schedules`

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	return count, nil
}

func (r *scheduleRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()

	setClause := ""
	args := make([]interface{}, 0, len(updates)+1)
	argPos := 1

	for key := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", key, argPos)
		args = append(args, updates[key])
		argPos++
	}

	query := fmt.Sprintf("UPDATE schedules SET %s WHERE id = $%d", setClause, argPos)
	args = append(args, id)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedules WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrScheduleNotFound
	}

	return nil
}

func (r *scheduleRepository) FindByEmployeeCount(ctx context.Context, employeeID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM employees WHERE schedule_id = $1`

	var count int64
	err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check schedule usage: %w", err)
	}

	return count, nil
}

func (r *scheduleRepository) FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.Schedule, error) {
	query := `
		SELECT id, company_id, name, time_in, time_out, allowed_late_minutes,
		       COALESCE(description, ''), created_at, updated_at
		FROM schedules
		WHERE id = $1 AND company_id = $2
	`

	row := r.pool.QueryRow(ctx, query, id, companyID)

	var schedule entity.Schedule
	err := row.Scan(
		&schedule.ID,
		&schedule.CompanyID,
		&schedule.Name,
		&schedule.TimeIn,
		&schedule.TimeOut,
		&schedule.AllowedLateMinutes,
		&schedule.Description,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to find schedule: %w", err)
	}

	return &schedule, nil
}

func (r *scheduleRepository) FindAllByCompany(ctx context.Context, companyID string, skip, limit int64) ([]*entity.Schedule, error) {
	query := `
		SELECT id, company_id, name, time_in, time_out, allowed_late_minutes,
		       COALESCE(description, ''), created_at, updated_at
		FROM schedules
		WHERE company_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, companyID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*entity.Schedule
	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.CompanyID,
			&schedule.Name,
			&schedule.TimeIn,
			&schedule.TimeOut,
			&schedule.AllowedLateMinutes,
			&schedule.Description,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, nil
}

func (r *scheduleRepository) CountByCompany(ctx context.Context, companyID string) (int64, error) {
	query := `SELECT COUNT(*) FROM schedules WHERE company_id = $1`

	var count int64
	err := r.pool.QueryRow(ctx, query, companyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	return count, nil
}

func (r *scheduleRepository) DeleteByCompany(ctx context.Context, id uuid.UUID, companyID string) error {
	query := `DELETE FROM schedules WHERE id = $1 AND company_id = $2`

	result, err := r.pool.Exec(ctx, query, id, companyID)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrScheduleNotFound
	}

	return nil
}
