package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"hris/internal/user/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type repository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &repository{pool: pool}
}

const userColumns = `id, company_id, name, email, password, role, is_active, profile_image_url, created_at, updated_at`

func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User
	var profileImageUrl *string
	var companyID *string
	err := row.Scan(
		&u.ID, &companyID, &u.Name, &u.Email, &u.Password,
		&u.Role, &u.IsActive, &profileImageUrl, &u.CreatedAt, &u.UpdatedAt,
	)
	if companyID != nil {
		u.CompanyID = *companyID
	}
	if profileImageUrl != nil {
		u.ProfileImageUrl = *profileImageUrl
	}
	return &u, err
}

func (r *repository) Create(ctx context.Context, user *entity.User) error {
	query := fmt.Sprintf(`INSERT INTO users (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, userColumns)

	now := time.Now()
	var profileUrl *string
	if user.ProfileImageUrl != "" {
		profileUrl = &user.ProfileImageUrl
	}
	var companyID *string
	if user.CompanyID != "" {
		companyID = &user.CompanyID
	}

	_, err := r.pool.Exec(ctx, query,
		user.ID, companyID, user.Name, user.Email, user.Password,
		user.Role, user.IsActive, profileUrl, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1`, userColumns)

	user, err := scanUser(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE email = $1`, userColumns)

	user, err := scanUser(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}

func (r *repository) FindAll(ctx context.Context, skip, limit int64) ([]*entity.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, userColumns)

	rows, err := r.pool.Query(ctx, query, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		var profileImageUrl *string
		var companyID *string
		if err := rows.Scan(
			&u.ID, &companyID, &u.Name, &u.Email, &u.Password,
			&u.Role, &u.IsActive, &profileImageUrl, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if companyID != nil {
			u.CompanyID = *companyID
		}
		if profileImageUrl != nil {
			u.ProfileImageUrl = *profileImageUrl
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET name = $1, email = $2, password = $3, role = $4, 
		is_active = $5, profile_image_url = $6, company_id = $7, updated_at = $8
		WHERE id = $9`

	var profileUrl *string
	if user.ProfileImageUrl != "" {
		profileUrl = &user.ProfileImageUrl
	}
	var companyID *string
	if user.CompanyID != "" {
		companyID = &user.CompanyID
	}

	result, err := r.pool.Exec(ctx, query,
		user.Name, user.Email, user.Password, user.Role,
		user.IsActive, profileUrl, companyID, time.Now(), user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) FindByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = ANY($1)`, userColumns)

	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by ids: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		var profileImageUrl *string
		var companyID *string
		if err := rows.Scan(
			&u.ID, &companyID, &u.Name, &u.Email, &u.Password,
			&u.Role, &u.IsActive, &profileImageUrl, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if companyID != nil {
			u.CompanyID = *companyID
		}
		if profileImageUrl != nil {
			u.ProfileImageUrl = *profileImageUrl
		}
		users = append(users, &u)
	}

	return users, nil
}

// FindAllByCompany returns users filtered by company_id for multi-tenancy
func (r *repository) FindAllByCompany(ctx context.Context, companyID string, skip, limit int64) ([]*entity.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users WHERE company_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userColumns)

	rows, err := r.pool.Query(ctx, query, companyID, limit, skip)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by company: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		var profileImageUrl *string
		var companyID *string
		if err := rows.Scan(
			&u.ID, &companyID, &u.Name, &u.Email, &u.Password,
			&u.Role, &u.IsActive, &profileImageUrl, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if companyID != nil {
			u.CompanyID = *companyID
		}
		if profileImageUrl != nil {
			u.ProfileImageUrl = *profileImageUrl
		}
		users = append(users, &u)
	}

	return users, nil
}

// CountByCompany returns the count of users in a specific company
func (r *repository) CountByCompany(ctx context.Context, companyID string) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE company_id = $1`, companyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by company: %w", err)
	}
	return count, nil
}
