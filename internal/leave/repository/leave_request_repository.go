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
	ErrLeaveRequestNotFound = errors.New("leave request not found")
)

const leaveRequestCols = `id, employee_id, leave_type_id, start_date, end_date, total_days, reason, attachment_url, emergency_contact, status, approved_by, approved_at, rejection_reason, created_at, updated_at, company_id`

type LeaveRequestRepository interface {
	Create(ctx context.Context, request *entity.LeaveRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.LeaveRequest, error)
	FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.LeaveRequest, error)
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]entity.LeaveRequest, error)
	CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error)
	FindAll(ctx context.Context, limit, offset int) ([]entity.LeaveRequest, error)
	FindAllByCompany(ctx context.Context, companyID string, limit, offset int) ([]entity.LeaveRequest, error)
	CountAll(ctx context.Context) (int64, error)
	CountAllByCompany(ctx context.Context, companyID string) (int64, error)
	FindPending(ctx context.Context) ([]entity.LeaveRequest, error)
	FindPendingByCompany(ctx context.Context, companyID string) ([]entity.LeaveRequest, error)
	FindByEmployeeAndDateRange(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) ([]entity.LeaveRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, rejectionReason *string) error
}

type leaveRequestRepository struct {
	pool *pgxpool.Pool
}

func NewLeaveRequestRepository(pool *pgxpool.Pool) LeaveRequestRepository {
	return &leaveRequestRepository{pool: pool}
}

func (r *leaveRequestRepository) Create(ctx context.Context, request *entity.LeaveRequest) error {
	query := `
		INSERT INTO leave_requests
		(id, employee_id, leave_type_id, start_date, end_date, total_days, reason, attachment_url, emergency_contact, status, company_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.pool.Exec(ctx, query,
		request.ID,
		request.EmployeeID,
		request.LeaveTypeID,
		request.StartDate,
		request.EndDate,
		request.TotalDays,
		request.Reason,
		request.AttachmentURL,
		request.EmergencyContact,
		request.Status,
		request.CompanyID,
	)

	return err
}

func (r *leaveRequestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.LeaveRequest, error) {
	query := `SELECT ${leaveRequestCols} FROM leave_requests WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var request entity.LeaveRequest
	err := row.Scan(
		&request.ID,
		&request.EmployeeID,
		&request.LeaveTypeID,
		&request.StartDate,
		&request.EndDate,
		&request.TotalDays,
		&request.Reason,
		&request.AttachmentURL,
		&request.EmergencyContact,
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
			return nil, ErrLeaveRequestNotFound
		}
		return nil, err
	}

	return &request, nil
}

func (r *leaveRequestRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]entity.LeaveRequest, error) {
	query := `
		SELECT  FROM leave_requests
		WHERE employee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, employeeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.LeaveRequest
	for rows.Next() {
		var request entity.LeaveRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.LeaveTypeID,
			&request.StartDate,
			&request.EndDate,
			&request.TotalDays,
			&request.Reason,
			&request.AttachmentURL,
			&request.EmergencyContact,
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

func (r *leaveRequestRepository) FindPending(ctx context.Context) ([]entity.LeaveRequest, error) {
	query := `
		SELECT  FROM leave_requests
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.LeaveRequest
	for rows.Next() {
		var request entity.LeaveRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.LeaveTypeID,
			&request.StartDate,
			&request.EndDate,
			&request.TotalDays,
			&request.Reason,
			&request.AttachmentURL,
			&request.EmergencyContact,
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

func (r *leaveRequestRepository) FindByEmployeeAndDateRange(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) ([]entity.LeaveRequest, error) {
	query := `
		SELECT  FROM leave_requests
		WHERE employee_id = $1
		  AND status IN ('PENDING', 'APPROVED')
		  AND (
			  (start_date <= $2 AND end_date >= $2) OR
			  (start_date <= $3 AND end_date >= $3) OR
			  (start_date >= $2 AND end_date <= $3)
		  )
	`

	rows, err := r.pool.Query(ctx, query, employeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.LeaveRequest
	for rows.Next() {
		var request entity.LeaveRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.LeaveTypeID,
			&request.StartDate,
			&request.EndDate,
			&request.TotalDays,
			&request.Reason,
			&request.AttachmentURL,
			&request.EmergencyContact,
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

func (r *leaveRequestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, rejectionReason *string) error {
	query := `
		UPDATE leave_requests
		SET status = $2, approved_by = $3, approved_at = $4, rejection_reason = $5, updated_at = NOW()
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.pool.Exec(ctx, query, id, status, approvedBy, &now, rejectionReason)
	return err
}

func (r *leaveRequestRepository) CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM leave_requests WHERE employee_id = $1`
	var count int64
	err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
	return count, err
}

func (r *leaveRequestRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.LeaveRequest, error) {
	query := `
		SELECT  FROM leave_requests
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.LeaveRequest
	for rows.Next() {
		var request entity.LeaveRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.LeaveTypeID,
			&request.StartDate,
			&request.EndDate,
			&request.TotalDays,
			&request.Reason,
			&request.AttachmentURL,
			&request.EmergencyContact,
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

func (r *leaveRequestRepository) CountAll(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM leave_requests`
	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func (r *leaveRequestRepository) FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*entity.LeaveRequest, error) {
	query := `SELECT  FROM leave_requests lr lr JOIN employees e ON lr.employee_id = e.id WHERE lr.id = $1 AND e.company_id = $2`

	row := r.pool.QueryRow(ctx, query, id, companyID)

	var request entity.LeaveRequest
	err := row.Scan(
		&request.ID,
		&request.EmployeeID,
		&request.LeaveTypeID,
		&request.StartDate,
		&request.EndDate,
		&request.TotalDays,
		&request.Reason,
		&request.AttachmentURL,
		&request.EmergencyContact,
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
			return nil, ErrLeaveRequestNotFound
		}
		return nil, err
	}

	return &request, nil
}

func (r *leaveRequestRepository) FindAllByCompany(ctx context.Context, companyID string, limit, offset int) ([]entity.LeaveRequest, error) {
	query := `
		SELECT  FROM leave_requests lr lr
		JOIN employees e ON lr.employee_id = e.id
		WHERE e.company_id = $1
		ORDER BY lr.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.LeaveRequest
	for rows.Next() {
		var request entity.LeaveRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.LeaveTypeID,
			&request.StartDate,
			&request.EndDate,
			&request.TotalDays,
			&request.Reason,
			&request.AttachmentURL,
			&request.EmergencyContact,
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

func (r *leaveRequestRepository) CountAllByCompany(ctx context.Context, companyID string) (int64, error) {
	query := `SELECT COUNT(*) FROM leave_requests lr JOIN employees e ON lr.employee_id = e.id WHERE e.company_id = $1`
	var count int64
	err := r.pool.QueryRow(ctx, query, companyID).Scan(&count)
	return count, err
}

func (r *leaveRequestRepository) FindPendingByCompany(ctx context.Context, companyID string) ([]entity.LeaveRequest, error) {
	query := `
		SELECT  FROM leave_requests lr lr
		JOIN employees e ON lr.employee_id = e.id
		WHERE lr.status = 'PENDING' AND e.company_id = $1
		ORDER BY lr.created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.LeaveRequest
	for rows.Next() {
		var request entity.LeaveRequest
		err := rows.Scan(
			&request.ID,
			&request.EmployeeID,
			&request.LeaveTypeID,
			&request.StartDate,
			&request.EndDate,
			&request.TotalDays,
			&request.Reason,
			&request.AttachmentURL,
			&request.EmergencyContact,
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
