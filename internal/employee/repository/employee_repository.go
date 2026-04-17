package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmployeeNotFound = errors.New("employee not found")
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *Employee) error
	FindByID(ctx context.Context, id uuid.UUID) (*EmployeeWithUser, error)
	FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*EmployeeWithUser, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*EmployeeWithUser, error)
	FindByUserIDAndCompany(ctx context.Context, userID uuid.UUID, companyID string) (*EmployeeWithUser, error)
	FindAll(ctx context.Context, page, perPage int, search string) ([]EmployeeWithUser, int64, error)
	FindAllByCompany(ctx context.Context, companyID string, page, perPage int, search string) ([]EmployeeWithUser, int64, error)
	FindAllWithoutPagination(ctx context.Context) ([]EmployeeWithUser, error)
	FindAllWithoutPaginationByCompany(ctx context.Context, companyID string) ([]EmployeeWithUser, error)
	Update(ctx context.Context, employee *Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByCompany(ctx context.Context, id uuid.UUID, companyID string) error
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*Employee, error)
	CountByDepartmentIDs(ctx context.Context, ids []uuid.UUID, companyID string) (map[uuid.UUID]int, error)
}

type Employee struct {
	ID                uuid.UUID  `db:"id"`
	CompanyID         uuid.UUID  `db:"company_id"`
	UserID            uuid.UUID  `db:"user_id"`
	FullName          string     `db:"full_name"`
	Position          string     `db:"position"`
	PhoneNumber       string     `db:"phone_number"`
	Address           string     `db:"address"`
	SalaryBase        float64    `db:"salary_base"`
	JoinDate          time.Time  `db:"join_date"`
	BankName          string     `db:"bank_name"`
	BankAccountNumber string     `db:"bank_account_number"`
	BankAccountHolder string     `db:"bank_account_holder"`
	ScheduleID        *uuid.UUID `db:"schedule_id"`
	EmploymentStatus  string     `db:"employment_status"`
	JobLevel          string     `db:"job_level"`
	Gender            string     `db:"gender"`
	Division          string     `db:"division"`
	DepartmentID      *uuid.UUID `db:"department_id"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

type EmployeeWithUser struct {
	ID                         uuid.UUID  `db:"id"`
	CompanyID                  uuid.UUID  `db:"company_id"`
	UserID                     uuid.UUID  `db:"user_id"`
	FullName                   string     `db:"full_name"`
	UserName                   *string    `db:"user_name"`
	UserEmail                  *string    `db:"user_email"`
	Position                   string     `db:"position"`
	PhoneNumber                string     `db:"phone_number"`
	Address                    string     `db:"address"`
	SalaryBase                 float64    `db:"salary_base"`
	JoinDate                   time.Time  `db:"join_date"`
	BankName                   string     `db:"bank_name"`
	BankAccountNumber          string     `db:"bank_account_number"`
	BankAccountHolder          string     `db:"bank_account_holder"`
	ScheduleID                 *uuid.UUID `db:"schedule_id"`
	ScheduleName               *string    `db:"schedule_name"`
	ScheduleTimeIn             *string    `db:"schedule_time_in"`
	ScheduleTimeOut            *string    `db:"schedule_time_out"`
	ScheduleAllowedLateMinutes *int       `db:"schedule_allowed_late_minutes"`
	EmploymentStatus           string     `db:"employment_status"`
	JobLevel                   string     `db:"job_level"`
	Gender                     string     `db:"gender"`
	Division                   string     `db:"division"`
	DepartmentID               *uuid.UUID `db:"department_id"`
	DepartmentName             *string    `db:"department_name"`
	CreatedAt                  time.Time  `db:"created_at"`
	UpdatedAt                  time.Time  `db:"updated_at"`
}

type employeeRepository struct {
	pool *pgxpool.Pool
}

const employeeWithUserSelectCols = `
	e.id, e.company_id, e.user_id, e.full_name,
	u.name as user_name, u.email as user_email,
	e.position, e.phone_number, e.address, e.salary_base, e.join_date,
	e.bank_name, e.bank_account_number, e.bank_account_holder, e.schedule_id,
	s.name as schedule_name, s.time_in as schedule_time_in,
	s.time_out as schedule_time_out, s.allowed_late_minutes as schedule_allowed_late_minutes,
	e.employment_status, e.job_level, e.gender, e.division, e.department_id,
	d.name as department_name,
	e.created_at, e.updated_at
`

const employeeWithUserJoins = `
	FROM employees e
	LEFT JOIN users u ON e.user_id = u.id
	LEFT JOIN schedules s ON e.schedule_id = s.id
	LEFT JOIN departments d ON e.department_id = d.id
`

func NewEmployeeRepository(pool *pgxpool.Pool) EmployeeRepository {
	return &employeeRepository{pool: pool}
}

func scanEmployeeWithUser(row interface {
	Scan(dest ...interface{}) error
}) (*EmployeeWithUser, error) {
	var e EmployeeWithUser
	err := row.Scan(
		&e.ID, &e.CompanyID, &e.UserID, &e.FullName,
		&e.UserName, &e.UserEmail,
		&e.Position, &e.PhoneNumber, &e.Address, &e.SalaryBase, &e.JoinDate,
		&e.BankName, &e.BankAccountNumber, &e.BankAccountHolder, &e.ScheduleID,
		&e.ScheduleName, &e.ScheduleTimeIn, &e.ScheduleTimeOut, &e.ScheduleAllowedLateMinutes,
		&e.EmploymentStatus, &e.JobLevel, &e.Gender, &e.Division, &e.DepartmentID,
		&e.DepartmentName,
		&e.CreatedAt, &e.UpdatedAt,
	)
	return &e, err
}

func (r *employeeRepository) Create(ctx context.Context, employee *Employee) error {
	query := `
		INSERT INTO employees (id, company_id, user_id, full_name, position, phone_number, address, salary_base, join_date, bank_name, bank_account_number, bank_account_holder, schedule_id, employment_status, job_level, gender, division, department_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := r.pool.Exec(ctx, query,
		employee.ID,
		employee.CompanyID,
		employee.UserID,
		employee.FullName,
		employee.Position,
		employee.PhoneNumber,
		employee.Address,
		employee.SalaryBase,
		employee.JoinDate,
		employee.BankName,
		employee.BankAccountNumber,
		employee.BankAccountHolder,
		employee.ScheduleID,
		employee.EmploymentStatus,
		employee.JobLevel,
		employee.Gender,
		employee.Division,
		employee.DepartmentID,
		employee.CreatedAt,
		employee.UpdatedAt,
	)

	return err
}

func (r *employeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*EmployeeWithUser, error) {
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` WHERE e.id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	emp, err := scanEmployeeWithUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return emp, nil
}

func (r *employeeRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*EmployeeWithUser, error) {
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` WHERE e.user_id = $1`

	row := r.pool.QueryRow(ctx, query, userID)
	emp, err := scanEmployeeWithUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return emp, nil
}

func (r *employeeRepository) FindAll(ctx context.Context, page, perPage int, search string) ([]EmployeeWithUser, int64, error) {
	offset := (page - 1) * perPage

	// Build count query with optional search
	countQuery := `SELECT COUNT(*) FROM employees e`
	args := []interface{}{}
	argNum := 1

	if search != "" {
		countQuery += ` LEFT JOIN users u ON e.user_id = u.id WHERE (e.full_name ILIKE $1 OR u.name ILIKE $1 OR u.email ILIKE $1 OR e.position ILIKE $1)`
		args = append(args, "%"+search+"%")
		argNum++
	}

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build data query
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins

	if search != "" {
		query += ` WHERE (e.full_name ILIKE $1 OR u.name ILIKE $1 OR u.email ILIKE $1 OR e.position ILIKE $1)`
	}

	query += ` ORDER BY e.created_at DESC LIMIT $` + fmt.Sprintf("%d", argNum) + ` OFFSET $` + fmt.Sprintf("%d", argNum+1)
	args = append(args, perPage, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []EmployeeWithUser
	for rows.Next() {
		emp, err := scanEmployeeWithUser(rows)
		if err != nil {
			return nil, 0, err
		}
		employees = append(employees, *emp)
	}

	return employees, total, nil
}

func (r *employeeRepository) FindAllWithoutPagination(ctx context.Context) ([]EmployeeWithUser, error) {
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` ORDER BY e.created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []EmployeeWithUser
	for rows.Next() {
		emp, err := scanEmployeeWithUser(rows)
		if err != nil {
			return nil, err
		}
		employees = append(employees, *emp)
	}

	return employees, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *Employee) error {
	query := `
		UPDATE employees
		SET full_name = $2, position = $3, phone_number = $4, address = $5, salary_base = $6,
		    bank_name = $7, bank_account_number = $8, bank_account_holder = $9,
		    schedule_id = $10, employment_status = $11, job_level = $12, gender = $13,
		    division = $14, department_id = $15, updated_at = $16
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		employee.ID,
		employee.FullName,
		employee.Position,
		employee.PhoneNumber,
		employee.Address,
		employee.SalaryBase,
		employee.BankName,
		employee.BankAccountNumber,
		employee.BankAccountHolder,
		employee.ScheduleID,
		employee.EmploymentStatus,
		employee.JobLevel,
		employee.Gender,
		employee.Division,
		employee.DepartmentID,
		employee.UpdatedAt,
	)

	return err
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) FindByIDAndCompany(ctx context.Context, id uuid.UUID, companyID string) (*EmployeeWithUser, error) {
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` WHERE e.id = $1 AND e.company_id = $2`

	row := r.pool.QueryRow(ctx, query, id, companyID)
	emp, err := scanEmployeeWithUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return emp, nil
}

func (r *employeeRepository) FindByUserIDAndCompany(ctx context.Context, userID uuid.UUID, companyID string) (*EmployeeWithUser, error) {
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` WHERE e.user_id = $1 AND e.company_id = $2`

	row := r.pool.QueryRow(ctx, query, userID, companyID)
	emp, err := scanEmployeeWithUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	return emp, nil
}

func (r *employeeRepository) FindAllByCompany(ctx context.Context, companyID string, page, perPage int, search string) ([]EmployeeWithUser, int64, error) {
	offset := (page - 1) * perPage

	// Build count query with optional search and company filter
	countQuery := `SELECT COUNT(*) FROM employees e WHERE e.company_id = $1`
	countArgs := []interface{}{companyID}
	argNum := 2

	var dataQueryArgs []interface{}
	dataQueryArgs = append(dataQueryArgs, companyID)

	if search != "" {
		countQuery += ` AND (e.full_name ILIKE $` + fmt.Sprintf("%d", argNum) + ` OR e.position ILIKE $` + fmt.Sprintf("%d", argNum) + `)`
		countArgs = append(countArgs, "%"+search+"%")
		dataQueryArgs = append(dataQueryArgs, "%"+search+"%")
		argNum++
	}

	var total int64
	err := r.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build data query
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` WHERE e.company_id = $1`

	if search != "" {
		query += ` AND (e.full_name ILIKE $` + fmt.Sprintf("%d", 2) + ` OR e.position ILIKE $` + fmt.Sprintf("%d", 2) + `)`
	}

	query += ` ORDER BY e.created_at DESC LIMIT $` + fmt.Sprintf("%d", argNum) + ` OFFSET $` + fmt.Sprintf("%d", argNum+1)
	dataQueryArgs = append(dataQueryArgs, perPage, offset)

	rows, err := r.pool.Query(ctx, query, dataQueryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []EmployeeWithUser
	for rows.Next() {
		emp, err := scanEmployeeWithUser(rows)
		if err != nil {
			return nil, 0, err
		}
		employees = append(employees, *emp)
	}

	return employees, total, nil
}

func (r *employeeRepository) FindAllWithoutPaginationByCompany(ctx context.Context, companyID string) ([]EmployeeWithUser, error) {
	query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins + ` WHERE e.company_id = $1 ORDER BY e.created_at DESC`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []EmployeeWithUser
	for rows.Next() {
		emp, err := scanEmployeeWithUser(rows)
		if err != nil {
			return nil, err
		}
		employees = append(employees, *emp)
	}

	return employees, nil
}

func (r *employeeRepository) DeleteByCompany(ctx context.Context, id uuid.UUID, companyID string) error {
	query := `DELETE FROM employees WHERE id = $1 AND company_id = $2`

	result, err := r.pool.Exec(ctx, query, id, companyID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*Employee, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `SELECT id, user_id, full_name, position, salary_base, phone_number, address,
              employment_status, join_date, schedule_id, bank_name, bank_account_number,
              bank_account_holder, division, job_level, gender, department_id, created_at, updated_at
              FROM employees WHERE id = ANY($1)`

	rows, err := r.pool.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*Employee
	for rows.Next() {
		emp := &Employee{}
		err := rows.Scan(&emp.ID, &emp.UserID, &emp.FullName, &emp.Position, &emp.SalaryBase,
			&emp.PhoneNumber, &emp.Address, &emp.EmploymentStatus, &emp.JoinDate,
			&emp.ScheduleID, &emp.BankName, &emp.BankAccountNumber, &emp.BankAccountHolder,
			&emp.Division, &emp.JobLevel, &emp.Gender, &emp.DepartmentID, &emp.CreatedAt, &emp.UpdatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}
	return employees, nil
}

func (r *employeeRepository) CountByDepartmentIDs(ctx context.Context, ids []uuid.UUID, companyID string) (map[uuid.UUID]int, error) {
	if len(ids) == 0 {
		return make(map[uuid.UUID]int), nil
	}

	query := `
		SELECT e.department_id, COUNT(e.id)
		FROM employees e
		WHERE e.department_id = ANY($1) AND e.company_id = $2
		GROUP BY e.department_id
	`

	rows, err := r.pool.Query(ctx, query, ids, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]int)
	for rows.Next() {
		var deptID uuid.UUID
		var count int
		if err := rows.Scan(&deptID, &count); err != nil {
			return nil, err
		}
		result[deptID] = count
	}

	return result, nil
}
