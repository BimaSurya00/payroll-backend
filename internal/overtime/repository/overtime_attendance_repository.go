package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/overtime/entity"
)

var (
	ErrOvertimeAttendanceNotFound = errors.New("overtime attendance not found")
)

type OvertimeAttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.OvertimeAttendance) error
	FindByRequestID(ctx context.Context, requestID uuid.UUID) (*entity.OvertimeAttendance, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]entity.OvertimeAttendance, error)
	UpdateClockIn(ctx context.Context, id uuid.UUID, clockInTime time.Time, notes string) error
	UpdateClockOut(ctx context.Context, id uuid.UUID, clockOutTime time.Time, actualHours float64, notes string) error
}

type overtimeAttendanceRepository struct {
	pool *pgxpool.Pool
}

func NewOvertimeAttendanceRepository(pool *pgxpool.Pool) OvertimeAttendanceRepository {
	return &overtimeAttendanceRepository{pool: pool}
}

func (r *overtimeAttendanceRepository) Create(ctx context.Context, attendance *entity.OvertimeAttendance) error {
	query := `
		INSERT INTO overtime_attendance (id, overtime_request_id, employee_id, clock_in_time, clock_out_time,
			actual_hours, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.pool.Exec(ctx, query,
		attendance.ID,
		attendance.OvertimeRequestID,
		attendance.EmployeeID,
		attendance.ClockInTime,
		attendance.ClockOutTime,
		attendance.ActualHours,
		attendance.Notes,
		attendance.CreatedAt,
	)

	return err
}

func (r *overtimeAttendanceRepository) FindByRequestID(ctx context.Context, requestID uuid.UUID) (*entity.OvertimeAttendance, error) {
	query := `SELECT * FROM overtime_attendance WHERE overtime_request_id = $1`

	row := r.pool.QueryRow(ctx, query, requestID)

	var attendance entity.OvertimeAttendance
	err := row.Scan(
		&attendance.ID,
		&attendance.OvertimeRequestID,
		&attendance.EmployeeID,
		&attendance.ClockInTime,
		&attendance.ClockOutTime,
		&attendance.ActualHours,
		&attendance.Notes,
		&attendance.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOvertimeAttendanceNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

func (r *overtimeAttendanceRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]entity.OvertimeAttendance, error) {
	offset := (page - 1) * perPage
	query := `
		SELECT * FROM overtime_attendance
		WHERE employee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, employeeID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []entity.OvertimeAttendance
	for rows.Next() {
		var attendance entity.OvertimeAttendance
		err := rows.Scan(
			&attendance.ID,
			&attendance.OvertimeRequestID,
			&attendance.EmployeeID,
			&attendance.ClockInTime,
			&attendance.ClockOutTime,
			&attendance.ActualHours,
			&attendance.Notes,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}

	return attendances, nil
}

func (r *overtimeAttendanceRepository) UpdateClockIn(ctx context.Context, id uuid.UUID, clockInTime time.Time, notes string) error {
	query := `
		UPDATE overtime_attendance
		SET clock_in_time = $2, notes = $3
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id, clockInTime, notes)
	return err
}

func (r *overtimeAttendanceRepository) UpdateClockOut(ctx context.Context, id uuid.UUID, clockOutTime time.Time, actualHours float64, notes string) error {
	query := `
		UPDATE overtime_attendance
		SET clock_out_time = $2, actual_hours = $3, notes = $4
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, id, clockOutTime, actualHours, notes)
	return err
}
