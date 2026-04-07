package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/overtime/entity"
)

var (
	ErrOvertimeRequestNotFound = errors.New("overtime request not found")
)

type OvertimeRequestRepository interface {
	Create(ctx context.Context, request *entity.OvertimeRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.OvertimeRequest, error)
	FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.OvertimeRequest, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]entity.OvertimeRequest, error)
	FindAll(ctx context.Context, page, perPage int, status string, employeeID uuid.UUID) ([]entity.OvertimeRequest, error)
	FindAllByCompany(ctx context.Context, companyID string, page, perPage int, status string, employeeID uuid.UUID) ([]entity.OvertimeRequest, error)
	FindByEmployeeIDAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*entity.OvertimeRequest, error)
	FindPending(ctx context.Context) ([]entity.OvertimeRequest, error)
	FindPendingByCompany(ctx context.Context, companyID string) ([]entity.OvertimeRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, rejectionReason *string) error
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.OvertimeRequest, error)
	FindByDateRangeAndCompany(ctx context.Context, companyID string, startDate, endDate time.Time) ([]entity.OvertimeRequest, error)
}

type overtimeRequestRepository struct {
	pool *pgxpool.Pool
}

func NewOvertimeRequestRepository(pool *pgxpool.Pool) OvertimeRequestRepository {
	return &overtimeRequestRepository{pool: pool}
}

func (r *overtimeRequestRepository) Create(ctx context.Context, request *entity.OvertimeRequest) error {
	query := `
		INSERT INTO overtime_requests (id, employee_id, overtime_date, start_time, end_time, total_hours,
			reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.pool.Exec(ctx, query,
		request.ID,
		request.EmployeeID,
		request.OvertimeDate,
		request.StartTime,
		request.EndTime,
		request.TotalHours,
		request.Reason,
		request.OvertimePolicyID,
		request.Status,
		request.ApprovedBy,
		request.ApprovedAt,
		request.RejectionReason,
		request.CreatedAt,
		request.UpdatedAt,
	)

	return err
}

func (r *overtimeRequestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.OvertimeRequest, error) {
	query := `SELECT id, employee_id, overtime_date, start_time, end_time, total_hours, reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id FROM overtime_requests WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var request entity.OvertimeRequest
	err := row.Scan(
		&request.ID,
		&request.EmployeeID,
		&request.OvertimeDate,
		&request.StartTime,
		&request.EndTime,
		&request.TotalHours,
		&request.Reason,
		&request.OvertimePolicyID,
		&request.Status,
		&request.ApprovedBy,
		&request.ApprovedAt,
		&request.RejectionReason,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOvertimeRequestNotFound
		}
		return nil, err
	}

	return &request, nil
}

func (r *overtimeRequestRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, page, perPage int) ([]entity.OvertimeRequest, error) {
	offset := (page - 1) * perPage
	query := `
		SELECT id, employee_id, overtime_date, start_time, end_time, total_hours, reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id FROM overtime_requests
		WHERE employee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, employeeID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *overtimeRequestRepository) FindAll(ctx context.Context, page, perPage int, status string, employeeID uuid.UUID) ([]entity.OvertimeRequest, error) {
	offset := (page - 1) * perPage

	var query string
	var args []interface{}

	baseQuery := `SELECT id, employee_id, overtime_date, start_time, end_time, total_hours, reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id FROM overtime_requests`
	whereClause := "WHERE 1=1"
	orderClause := "ORDER BY created_at DESC"
	limitClause := "LIMIT $1 OFFSET $2"

	argIndex := 3
	if status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if employeeID != uuid.Nil {
		whereClause += fmt.Sprintf(" AND employee_id = $%d", argIndex)
		args = append(args, employeeID)
		argIndex++
	}

	query = fmt.Sprintf("%s %s %s %s", baseQuery, whereClause, orderClause, limitClause)
	args = append([]interface{}{perPage, offset}, args...)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *overtimeRequestRepository) FindByEmployeeIDAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*entity.OvertimeRequest, error) {
	query := `
		SELECT id, employee_id, overtime_date, start_time, end_time, total_hours, reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id FROM overtime_requests
		WHERE employee_id = $1 AND overtime_date = $2
	`

	row := r.pool.QueryRow(ctx, query, employeeID, date)

	var request entity.OvertimeRequest
	err := row.Scan(
		&request.ID,
		&request.EmployeeID,
		&request.OvertimeDate,
		&request.StartTime,
		&request.EndTime,
		&request.TotalHours,
		&request.Reason,
		&request.OvertimePolicyID,
		&request.Status,
		&request.ApprovedBy,
		&request.ApprovedAt,
		&request.RejectionReason,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOvertimeRequestNotFound
		}
		return nil, err
	}

	return &request, nil
}

func (r *overtimeRequestRepository) FindPending(ctx context.Context) ([]entity.OvertimeRequest, error) {
	query := `
		SELECT id, employee_id, overtime_date, start_time, end_time, total_hours, reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id FROM overtime_requests
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *overtimeRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, rejectionReason *string) error {
	query := `
		UPDATE overtime_requests
		SET status = $2, approved_by = $3, approved_at = $4, rejection_reason = $5, updated_at = $6
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query, id, status, approvedBy, now, rejectionReason, now)
	return err
}

func (r *overtimeRequestRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.OvertimeRequest, error) {
	query := `
		SELECT id, employee_id, overtime_date, start_time, end_time, total_hours, reason, overtime_policy_id, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id FROM overtime_requests
		WHERE overtime_date BETWEEN $1 AND $2
		ORDER BY overtime_date ASC
	`

	rows, err := r.pool.Query(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *overtimeRequestRepository) FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.OvertimeRequest, error) {
	query := `SELECT o.* FROM overtime_requests o JOIN employees e ON o.employee_id = e.id WHERE o.id = $1 AND e.company_id = $2`

	row := r.pool.QueryRow(ctx, query, id, companyID)

	var request entity.OvertimeRequest
	err := row.Scan(
		&request.ID,
		&request.EmployeeID,
		&request.OvertimeDate,
		&request.StartTime,
		&request.EndTime,
		&request.TotalHours,
		&request.Reason,
		&request.OvertimePolicyID,
		&request.Status,
		&request.ApprovedBy,
		&request.ApprovedAt,
		&request.RejectionReason,
		&request.CreatedAt,
		&request.UpdatedAt,
		&request.CompanyID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOvertimeRequestNotFound
		}
		return nil, err
	}

	return &request, nil
}

func (r *overtimeRequestRepository) FindAllByCompany(ctx context.Context, companyID string, page, perPage int, status string, employeeID uuid.UUID) ([]entity.OvertimeRequest, error) {
	offset := (page - 1) * perPage

	var query string
	var args []interface{}

	baseQuery := `SELECT o.* FROM overtime_requests o JOIN employees e ON o.employee_id = e.id WHERE e.company_id = $1`
	whereClause := ""
	orderClause := "ORDER BY o.created_at DESC"
	limitClause := fmt.Sprintf("LIMIT $2 OFFSET $3")

	args = append(args, companyID, perPage, offset)
	argIndex := 4

	if status != "" {
		whereClause += fmt.Sprintf(" AND o.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if employeeID != uuid.Nil {
		whereClause += fmt.Sprintf(" AND o.employee_id = $%d", argIndex)
		args = append(args, employeeID)
		argIndex++
	}

	query = fmt.Sprintf("%s%s %s %s", baseQuery, whereClause, orderClause, limitClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *overtimeRequestRepository) FindPendingByCompany(ctx context.Context, companyID string) ([]entity.OvertimeRequest, error) {
	query := `
		SELECT o.* FROM overtime_requests o
		JOIN employees e ON o.employee_id = e.id
		WHERE o.status = 'PENDING' AND e.company_id = $1
		ORDER BY o.created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *overtimeRequestRepository) FindByDateRangeAndCompany(ctx context.Context, companyID string, startDate, endDate time.Time) ([]entity.OvertimeRequest, error) {
	query := `
		SELECT o.* FROM overtime_requests o
		JOIN employees e ON o.employee_id = e.id
		WHERE e.company_id = $1 AND o.overtime_date BETWEEN $2 AND $3
		ORDER BY o.overtime_date ASC
	`

	rows, err := r.pool.Query(ctx, query, companyID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.OvertimeRequest
	for rows.Next() {
		var request entity.OvertimeRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.OvertimeDate,
			&request.StartTime,
			&request.EndTime,
			&request.TotalHours,
			&request.Reason,
			&request.OvertimePolicyID,
			&request.Status,
			&request.ApprovedBy,
			&request.ApprovedAt,
			&request.RejectionReason,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.CompanyID,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}
