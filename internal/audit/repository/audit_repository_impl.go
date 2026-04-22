package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/audit/entity"
)

type auditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) AuditRepository {
	return &auditRepository{pool: pool}
}

func (r *auditRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.pool.Exec(ctx, query,
		log.ID,
		log.UserID,
		log.UserName,
		log.Action,
		log.ResourceType,
		log.ResourceID,
		log.OldData,
		log.NewData,
		log.Metadata,
		log.IPAddress,
	)

	return err
}

func (r *auditRepository) FindByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*entity.AuditLog, error) {
	query := `
		SELECT id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address, created_at FROM audit_logs
		WHERE resource_type = $1 AND resource_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.pool.Query(ctx, query, resourceType, resourceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		log := &entity.AuditLog{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserName,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldData,
			&log.NewData,
			&log.Metadata,
			&log.IPAddress,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *auditRepository) FindByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.AuditLog, error) {
	query := `
		SELECT id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address, created_at FROM audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		log := &entity.AuditLog{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserName,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldData,
			&log.NewData,
			&log.Metadata,
			&log.IPAddress,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *auditRepository) FindAll(ctx context.Context, filter AuditFilter, limit, offset int) ([]*entity.AuditLog, error) {
	query := `SELECT id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address, created_at FROM audit_logs WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(` AND user_id = $%d`, argNum)
		args = append(args, *filter.UserID)
		argNum++
	}

	if filter.Action != nil {
		query += fmt.Sprintf(` AND action = $%d`, argNum)
		args = append(args, *filter.Action)
		argNum++
	}

	if filter.ResourceType != nil {
		query += fmt.Sprintf(` AND resource_type = $%d`, argNum)
		args = append(args, *filter.ResourceType)
		argNum++
	}

	if filter.CompanyID != nil {
		query += fmt.Sprintf(` AND company_id = $%d`, argNum)
		args = append(args, *filter.CompanyID)
		argNum++
	}

	if filter.DateFrom != nil {
		query += fmt.Sprintf(` AND created_at >= $%d`, argNum)
		args = append(args, *filter.DateFrom)
		argNum++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(` AND created_at <= $%d`, argNum)
		args = append(args, *filter.DateTo)
		argNum++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.AuditLog
	for rows.Next() {
		log := &entity.AuditLog{}
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserName,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldData,
			&log.NewData,
			&log.Metadata,
			&log.IPAddress,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *auditRepository) Count(ctx context.Context, filter AuditFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`
	args := []interface{}{}
	argNum := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(` AND user_id = $%d`, argNum)
		args = append(args, *filter.UserID)
		argNum++
	}

	if filter.Action != nil {
		query += fmt.Sprintf(` AND action = $%d`, argNum)
		args = append(args, *filter.Action)
		argNum++
	}

	if filter.ResourceType != nil {
		query += fmt.Sprintf(` AND resource_type = $%d`, argNum)
		args = append(args, *filter.ResourceType)
		argNum++
	}

	if filter.CompanyID != nil {
		query += fmt.Sprintf(` AND company_id = $%d`, argNum)
		args = append(args, *filter.CompanyID)
		argNum++
	}

	if filter.DateFrom != nil {
		query += fmt.Sprintf(` AND created_at >= $%d`, argNum)
		args = append(args, *filter.DateFrom)
		argNum++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(` AND created_at <= $%d`, argNum)
		args = append(args, *filter.DateTo)
		argNum++
	}

	var count int64
	err := r.pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}
