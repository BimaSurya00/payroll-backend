# MVP-37: Fix Employee Repository Consistency

**Estimasi**: 3 jam  
**Impact**: HIGH — Data Consistency  
**Prerequisite**: MVP-34 (Users migrated to PostgreSQL) ✅

---

## 1. Problem Statement

The employee repository ([`internal/employee/repository/employee_repository.go`](file:///home/bima/Documents/example.com/hris/internal/employee/repository/employee_repository.go)) has **6 consistency issues** across its 7 query methods:

| # | Issue | Affected Methods | Impact |
|---|-------|-----------------|--------|
| 1 | `Create()` skips `full_name` and `department_id` | `Create` | New employees have no name or department |
| 2 | `Update()` skips `full_name` and `department_id` | `Update` | Cannot change name or department |
| 3 | 3 of 5 READ queries missing `LEFT JOIN departments` | `FindByUserID`, `FindAll`, `FindAllWithoutPagination` | No department info returned for most views |
| 4 | 3 of 5 READ queries missing `user_name`/`user_email` from `users` table JOIN | `FindByUserID`, `FindAll`, `FindAllWithoutPagination` | Name/email always empty string |
| 5 | `FindByIDs()` missing `full_name`, `gender`, `department_id` columns | `FindByIDs` | Payroll batch missing employee details |
| 6 | Stale "users are in MongoDB" comments | `FindByID`, `FindByUserID`, `FindAll`, `FindAllWithoutPagination` | Misleading comments (users are now in PostgreSQL after MVP-34) |

---

## 2. Current State Analysis

### Entity Structs (Lines 28–76)

**`Employee` struct** — has `FullName`, `DepartmentID` fields ✅  
**`EmployeeWithUser` struct** — has `UserName`, `UserEmail`, `DepartmentID`, `DepartmentName` fields ✅

### Method Consistency Matrix

| Method | `full_name` in cols | `department_id` in cols | Dept JOIN | User JOIN | Schedule JOIN | User data populated |
|--------|:---:|:---:|:---:|:---:|:---:|:---:|
| `Create` (line 86) | ❌ | ❌ | N/A | N/A | N/A | N/A |
| `Update` (line 364) | ❌ | ❌ | N/A | N/A | N/A | N/A |
| `FindByID` (line 115) | ❌ | ✅ | ✅ | ❌ | ✅ | ❌ hardcoded "" |
| `FindByUserID` (line 176) | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ hardcoded "" |
| `FindAll` (line 233) | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ hardcoded "" |
| `FindAllWithoutPagination` (line 305) | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ hardcoded "" |
| `FindByIDs` (line 409) | ❌ | ❌ | ❌ | ❌ | ❌ | N/A (returns `Employee` not `EmployeeWithUser`) |

**Target**: All READ methods should have consistent column lists, dept JOIN, user JOIN (now possible since MVP-34 migrated users to PostgreSQL), and schedule JOIN.

---

## 3. Implementation Steps

### Step 1: Fix `Create()` — Add `full_name` and `department_id`

**File**: `employee_repository.go` (lines 86–113)

**Current INSERT** (line 88):
```sql
INSERT INTO employees (id, user_id, position, phone_number, address, salary_base, join_date, bank_name, bank_account_number, bank_account_holder, schedule_id, employment_status, job_level, gender, division, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
```

**Change to:**
```sql
INSERT INTO employees (id, user_id, full_name, position, phone_number, address, salary_base, join_date, bank_name, bank_account_number, bank_account_holder, schedule_id, employment_status, job_level, gender, division, department_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
```

**Update the Exec call** to include `employee.FullName` (after `employee.UserID`) and `employee.DepartmentID` (after `employee.Division`).

---

### Step 2: Fix `Update()` — Add `full_name` and `department_id`

**Current UPDATE** (lines 365–371):
```sql
UPDATE employees
SET position = $2, phone_number = $3, address = $4, salary_base = $5,
    bank_name = $6, bank_account_number = $7, bank_account_holder = $8,
    schedule_id = $9, employment_status = $10, job_level = $11, gender = $12,
    division = $13, updated_at = $14
WHERE id = $1
```

**Change to:**
```sql
UPDATE employees
SET full_name = $2, position = $3, phone_number = $4, address = $5, salary_base = $6,
    bank_name = $7, bank_account_number = $8, bank_account_holder = $9,
    schedule_id = $10, employment_status = $11, job_level = $12, gender = $13,
    division = $14, department_id = $15, updated_at = $16
WHERE id = $1
```

**Update the Exec call** arguments accordingly, adding `employee.FullName` (first after ID) and `employee.DepartmentID` (after Division).

---

### Step 3: Create Common Query Builder Constants

To avoid copy-pasting the same SELECT + JOIN block into every method, define a constant at the top of the file for the standard column list. This is optional but recommended:

```go
const employeeWithUserSelectCols = `
    e.id, e.user_id, e.full_name,
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
    LEFT JOIN users u ON e.user_id = u.id::uuid
    LEFT JOIN schedules s ON e.schedule_id = s.id
    LEFT JOIN departments d ON e.department_id = d.id
`
```

> **Important**: Check the type of `users.id` — if it's `VARCHAR`/`TEXT` (from the MongoDB migration), you need `u.id::uuid` or `e.user_id::text = u.id`. If it's already `UUID`, just use `e.user_id = u.id`. Check the migration file `000013_add_users_table.up.sql` for the exact type.

Then create a standard scan helper:

```go
func scanEmployeeWithUser(row interface{ Scan(dest ...interface{}) error }) (*EmployeeWithUser, error) {
    var e EmployeeWithUser
    err := row.Scan(
        &e.ID, &e.UserID, &e.FullName,
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
```

> **Note**: Add a `FullName` field to `EmployeeWithUser` struct if it doesn't exist already. Currently the struct has `UserName` (from users table) but the employees table also has `full_name`.

---

### Step 4: Fix `FindByID()` — Add User JOIN, remove hardcoded empty strings

**Current** (lines 115–174): Has dept JOIN ✅ and schedule JOIN ✅, but **no user JOIN** ❌.

Replace the query with the standard query using the constants from Step 3:

```go
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
```

**Remove** the hardcoded:
```go
// UserName and UserEmail will be empty since users are in MongoDB
employee.UserName = ""
employee.UserEmail = ""
```

---

### Step 5: Fix `FindByUserID()` — Add Dept JOIN, User JOIN, remove hardcoded empty strings

**Current** (lines 176–231): Has schedule JOIN ✅, but **no dept JOIN** ❌ and **no user JOIN** ❌.

```go
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
```

---

### Step 6: Fix `FindAll()` — Add Dept JOIN, User JOIN, add search on `full_name`/`user_name`

**Current** (lines 233–303): Has schedule JOIN ✅, but **no dept JOIN** ❌ and **no user JOIN** ❌. Search parameter is unused.

**Update count query** to support search:
```go
countQuery := `SELECT COUNT(*) FROM employees e`
args := []interface{}{}
argNum := 1

if search != "" {
    countQuery += ` LEFT JOIN users u ON e.user_id = u.id::uuid WHERE (e.full_name ILIKE $1 OR u.name ILIKE $1 OR u.email ILIKE $1 OR e.position ILIKE $1)`
    args = append(args, "%"+search+"%")
    argNum++
}
```

**Update data query:**
```go
query := `SELECT ` + employeeWithUserSelectCols + employeeWithUserJoins

if search != "" {
    query += fmt.Sprintf(` WHERE (e.full_name ILIKE $%d OR u.name ILIKE $%d OR u.email ILIKE $%d OR e.position ILIKE $%d)`, argNum, argNum, argNum, argNum)
    // Note: reuse same arg number since it's the same parameter
}

query += fmt.Sprintf(` ORDER BY e.created_at DESC LIMIT $%d OFFSET $%d`, argNum, argNum+1)
```

> **Important**: Need to add `fmt` import if not already present.

**Remove** hardcoded empty `UserName`/`UserEmail`.

---

### Step 7: Fix `FindAllWithoutPagination()` — Add Dept JOIN, User JOIN

**Current** (lines 305–362): Same issues as `FindAll`.

```go
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
```

---

### Step 8: Fix `FindByIDs()` — Add `full_name`, `gender`, `department_id`

**Current** (lines 409–438): Missing `full_name`, `gender`, `department_id`.

```go
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
```

---

### Step 9: Add `FullName` to `EmployeeWithUser` Struct (if missing)

Check if `EmployeeWithUser` struct already has `FullName`. Based on the current code (lines 50–76), it does NOT. Add it:

```diff
 type EmployeeWithUser struct {
     ID                 uuid.UUID  `db:"id"`
     UserID             uuid.UUID  `db:"user_id"`
+    FullName           string     `db:"full_name"`
     UserName           string     `db:"user_name"`
     UserEmail          string     `db:"user_email"`
```

---

### Step 10: Check and Fix `users.id` Column Type for JOIN

**File**: [`database/migrations/000013_add_users_table.up.sql`](file:///home/bima/Documents/example.com/hris/database/migrations/000013_add_users_table.up.sql)

If `users.id` is `VARCHAR(36)` (string from MongoDB era), the JOIN needs type casting:
```sql
LEFT JOIN users u ON e.user_id::text = u.id
```

If `users.id` is `UUID`, the JOIN is simpler:
```sql
LEFT JOIN users u ON e.user_id = u.id
```

**Check the migration file and adjust the JOIN accordingly.**

---

## 4. Files Changed Summary

| # | File | What Changes |
|---|------|-------------|
| 1 | [`employee_repository.go`](file:///home/bima/Documents/example.com/hris/internal/employee/repository/employee_repository.go) | Fix all 7 methods, add `FullName` to `EmployeeWithUser`, add constants + scan helper, remove stale MongoDB comments |

**Only 1 file needs changes**, but the changes are extensive (touching all 7 methods + adding constants).

---

## 5. Before vs After

| Metric | Before | After |
|--------|--------|-------|
| Methods with dept JOIN | 1 of 5 (20%) | 5 of 5 (100%) |
| Methods with user JOIN | 0 of 5 (0%) | 5 of 5 (100%) |
| `Create()` columns | 17 (missing `full_name`, `department_id`) | 19 (complete) |
| `Update()` columns | 13 (missing `full_name`, `department_id`) | 15 (complete) |
| `FindByIDs()` columns | 16 (missing `full_name`, `gender`, `department_id`) | 19 (complete) |
| Stale MongoDB comments | 4 | 0 |
| Code duplication (Scan blocks) | 5 copies × ~20 lines | 1 helper function |

---

## 6. Verification Plan

### Build Verification
```bash
cd /home/bima/Documents/example.com/hris
go build ./...
go vet ./internal/employee/...
```

### Grep Verification
```bash
# No more "users are in MongoDB" comments
grep -rn "MongoDB" internal/employee/
# Should return NO matches

# All READ methods should JOIN users table
grep -n "LEFT JOIN users" internal/employee/repository/employee_repository.go
# Should return 1 match (the constant) or 5 matches (inline)

# full_name should be in Create and Update
grep -n "full_name" internal/employee/repository/employee_repository.go
# Should return multiple matches

# department_id should be in Create and Update
grep -n "department_id" internal/employee/repository/employee_repository.go
# Should return multiple matches
```

---

## 7. Risk: `full_name` Column Existence

Before implementing, verify that the `employees` table has a `full_name` column in the database migrations. If it doesn't exist in any migration, a new migration `000014_add_full_name_to_employees.sql` must be created:

```sql
-- 000014_add_full_name_to_employees.up.sql
ALTER TABLE employees ADD COLUMN IF NOT EXISTS full_name VARCHAR(255) NOT NULL DEFAULT '';

-- 000014_add_full_name_to_employees.down.sql
ALTER TABLE employees DROP COLUMN IF EXISTS full_name;
```

Check with:
```bash
grep -rn "full_name" database/migrations/
```
