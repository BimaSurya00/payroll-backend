package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/attendance/entity"
)

var (
	ErrAttendanceNotFound = errors.New("attendance not found")
)

type attendanceRepository struct {
	pool *pgxpool.Pool
}

func NewAttendanceRepository(pool *pgxpool.Pool) AttendanceRepository {
	return &attendanceRepository{pool: pool}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *entity.Attendance) error {
	query := `
		INSERT INTO attendances (id, employee_id, schedule_id, date, clock_in_time,
		                        clock_out_time, clock_in_lat, clock_in_long, clock_out_lat,
		                        clock_out_long, status, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		attendance.ID,
		attendance.EmployeeID,
		attendance.ScheduleID,
		attendance.Date,
		attendance.ClockInTime,
		attendance.ClockOutTime,
		attendance.ClockInLat,
		attendance.ClockInLong,
		attendance.ClockOutLat,
		attendance.ClockOutLong,
		attendance.Status,
		attendance.Notes,
		attendance.CreatedAt,
	)

	return err
}

func (r *attendanceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Attendance, error) {
	query := `
		SELECT id, employee_id, schedule_id, date, clock_in_time, clock_out_time,
		       clock_in_lat, clock_in_long, clock_out_lat, clock_out_long, status, notes, created_at
		FROM attendances
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var attendance entity.Attendance
	err := row.Scan(
		&attendance.ID,
		&attendance.EmployeeID,
		&attendance.ScheduleID,
		&attendance.Date,
		&attendance.ClockInTime,
		&attendance.ClockOutTime,
		&attendance.ClockInLat,
		&attendance.ClockInLong,
		&attendance.ClockOutLat,
		&attendance.ClockOutLong,
		&attendance.Status,
		&attendance.Notes,
		&attendance.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("failed to find attendance: %w", err)
	}

	return &attendance, nil
}

func (r *attendanceRepository) FindByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*entity.Attendance, error) {
	query := `
		SELECT id, employee_id, schedule_id, date, clock_in_time, clock_out_time,
		       clock_in_lat, clock_in_long, clock_out_lat, clock_out_long, status, notes, created_at
		FROM attendances
		WHERE employee_id = $1 AND DATE(date) = DATE($2)
		LIMIT 1
	`

	row := r.pool.QueryRow(ctx, query, employeeID, date)

	var attendance entity.Attendance
	err := row.Scan(
		&attendance.ID,
		&attendance.EmployeeID,
		&attendance.ScheduleID,
		&attendance.Date,
		&attendance.ClockInTime,
		&attendance.ClockOutTime,
		&attendance.ClockInLat,
		&attendance.ClockInLong,
		&attendance.ClockOutLat,
		&attendance.ClockOutLong,
		&attendance.Status,
		&attendance.Notes,
		&attendance.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("failed to find attendance: %w", err)
	}

	return &attendance, nil
}

func (r *attendanceRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, skip, limit int64) ([]*entity.Attendance, error) {
	query := `
		SELECT id, employee_id, schedule_id, date, clock_in_time, clock_out_time,
		       clock_in_lat, clock_in_long, clock_out_lat, clock_out_long, status, notes, created_at
		FROM attendances
		WHERE employee_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, employeeID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find attendances: %w", err)
	}
	defer rows.Close()

	var attendances []*entity.Attendance
	for rows.Next() {
		var attendance entity.Attendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.EmployeeID,
			&attendance.ScheduleID,
			&attendance.Date,
			&attendance.ClockInTime,
			&attendance.ClockOutTime,
			&attendance.ClockInLat,
			&attendance.ClockInLong,
			&attendance.ClockOutLat,
			&attendance.ClockOutLong,
			&attendance.Status,
			&attendance.Notes,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		attendances = append(attendances, &attendance)
	}

	return attendances, nil
}

func (r *attendanceRepository) CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM attendances WHERE employee_id = $1`

	var count int64
	err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count attendances: %w", err)
	}

	return count, nil
}

func (r *attendanceRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
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

	query := fmt.Sprintf("UPDATE attendances SET %s WHERE id = $%d", setClause, argPos)
	args = append(args, id)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrAttendanceNotFound
	}

	return nil
}

func (r *attendanceRepository) FindAll(ctx context.Context, filter AttendanceFilter, skip, limit int64) ([]*entity.Attendance, error) {
	query := `
		SELECT a.id, a.employee_id, a.schedule_id, a.date, a.clock_in_time, a.clock_out_time,
		       a.clock_in_lat, a.clock_in_long, a.clock_out_lat, a.clock_out_long, a.status, a.notes, a.created_at
		FROM attendances a
		JOIN employees e ON a.employee_id = e.id
		WHERE 1=1
	`
	args := make([]interface{}, 0)
	argPos := 1

	if filter.CompanyID != nil {
		query += fmt.Sprintf(" AND e.company_id = $%d", argPos)
		args = append(args, *filter.CompanyID)
		argPos++
	}

	if filter.EmployeeID != nil {
		query += fmt.Sprintf(" AND a.employee_id = $%d", argPos)
		args = append(args, *filter.EmployeeID)
		argPos++
	}

	if filter.ScheduleID != nil {
		query += fmt.Sprintf(" AND a.schedule_id = $%d", argPos)
		args = append(args, *filter.ScheduleID)
		argPos++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND a.status = $%d", argPos)
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND a.date >= $%d", argPos)
		args = append(args, *filter.DateFrom)
		argPos++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND a.date <= $%d", argPos)
		args = append(args, *filter.DateTo)
		argPos++
	}

	query += " ORDER BY a.date DESC LIMIT $" + fmt.Sprint(argPos) + " OFFSET $" + fmt.Sprint(argPos+1)
	args = append(args, limit, skip)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find attendances: %w", err)
	}
	defer rows.Close()

	var attendances []*entity.Attendance
	for rows.Next() {
		var attendance entity.Attendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.EmployeeID,
			&attendance.ScheduleID,
			&attendance.Date,
			&attendance.ClockInTime,
			&attendance.ClockOutTime,
			&attendance.ClockInLat,
			&attendance.ClockInLong,
			&attendance.ClockOutLat,
			&attendance.ClockOutLong,
			&attendance.Status,
			&attendance.Notes,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attendance: %w", err)
		}
		attendances = append(attendances, &attendance)
	}

	return attendances, nil
}

func (r *attendanceRepository) CountAll(ctx context.Context, filter AttendanceFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM attendances a JOIN employees e ON a.employee_id = e.id WHERE 1=1`
	args := make([]interface{}, 0)
	argPos := 1

	if filter.CompanyID != nil {
		query += fmt.Sprintf(" AND e.company_id = $%d", argPos)
		args = append(args, *filter.CompanyID)
		argPos++
	}

	if filter.EmployeeID != nil {
		query += fmt.Sprintf(" AND a.employee_id = $%d", argPos)
		args = append(args, *filter.EmployeeID)
		argPos++
	}

	if filter.ScheduleID != nil {
		query += fmt.Sprintf(" AND a.schedule_id = $%d", argPos)
		args = append(args, *filter.ScheduleID)
		argPos++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND a.status = $%d", argPos)
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND a.date >= $%d", argPos)
		args = append(args, *filter.DateFrom)
		argPos++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND a.date <= $%d", argPos)
		args = append(args, *filter.DateTo)
		argPos++
	}

	var count int64
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count attendances: %w", err)
	}

	return count, nil
}

func (r *attendanceRepository) CountLateByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM attendances
              WHERE employee_id = $1 AND date >= $2 AND date <= $3 AND status = 'LATE'`
	var count int
	err := r.pool.QueryRow(ctx, query, employeeID, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *attendanceRepository) CountAbsentByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM attendances
              WHERE employee_id = $1 AND date >= $2 AND date <= $3 AND status = 'ABSENT'`
	var count int
	err := r.pool.QueryRow(ctx, query, employeeID, startDate, endDate).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *attendanceRepository) GetAttendanceSummaryByPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*AttendanceSummary, error) {
	query := `SELECT
                COALESCE(SUM(CASE WHEN status = 'PRESENT' THEN 1 ELSE 0 END), 0) as total_present,
                COALESCE(SUM(CASE WHEN status = 'LATE' THEN 1 ELSE 0 END), 0) as total_late,
                COALESCE(SUM(CASE WHEN status = 'ABSENT' THEN 1 ELSE 0 END), 0) as total_absent,
                COALESCE(SUM(CASE WHEN status = 'LEAVE' THEN 1 ELSE 0 END), 0) as total_leave,
                COUNT(*) as total_days
              FROM attendances
              WHERE employee_id = $1 AND date >= $2 AND date <= $3`

	summary := &AttendanceSummary{}
	err := r.pool.QueryRow(ctx, query, employeeID, startDate, endDate).Scan(
		&summary.TotalPresent, &summary.TotalLate, &summary.TotalAbsent,
		&summary.TotalLeave, &summary.TotalDays,
	)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (r *attendanceRepository) CountByStatusAndDate(ctx context.Context, status string, date time.Time) (int64, error) {
	query := `SELECT COUNT(*) FROM attendances WHERE status = $1 AND date = $2`
	var count int64
	err := r.pool.QueryRow(ctx, query, status, date).Scan(&count)
	return count, err
}

func (r *attendanceRepository) GetMonthlySummaryAll(ctx context.Context, startDate, endDate time.Time) ([]AttendanceMonthlySummary, error) {
	query := `SELECT employee_id,
		COALESCE(SUM(CASE WHEN status = 'PRESENT' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'LATE' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'ABSENT' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'LEAVE' THEN 1 ELSE 0 END), 0),
		COUNT(*)
		FROM attendances
		WHERE date >= $1 AND date <= $2
		GROUP BY employee_id`

	rows, err := r.pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []AttendanceMonthlySummary
	for rows.Next() {
		s := AttendanceMonthlySummary{}
		err := rows.Scan(&s.EmployeeID, &s.TotalPresent, &s.TotalLate,
			&s.TotalAbsent, &s.TotalLeave, &s.TotalDays)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (r *attendanceRepository) GetMonthlySummaryByCompany(ctx context.Context, companyID string, startDate, endDate time.Time) ([]AttendanceMonthlySummary, error) {
	query := `SELECT a.employee_id,
		COALESCE(SUM(CASE WHEN a.status = 'PRESENT' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status = 'LATE' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status = 'ABSENT' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status = 'LEAVE' THEN 1 ELSE 0 END), 0),
		COUNT(*)
		FROM attendances a
		JOIN employees e ON a.employee_id = e.id
		WHERE e.company_id = $1 AND a.date >= $2 AND a.date <= $3
		GROUP BY a.employee_id`

	rows, err := r.pool.Query(ctx, query, companyID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []AttendanceMonthlySummary
	for rows.Next() {
		s := AttendanceMonthlySummary{}
		err := rows.Scan(&s.EmployeeID, &s.TotalPresent, &s.TotalLate,
			&s.TotalAbsent, &s.TotalLeave, &s.TotalDays)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}
