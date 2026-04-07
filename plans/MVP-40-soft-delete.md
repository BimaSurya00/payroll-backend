# MVP-40: Soft Delete

**Estimasi**: 2 hari  
**Impact**: 🔴 TINGGI — Data Integrity & Compliance

---

## 1. Problem

Semua `DELETE` query menggunakan hard delete (`DELETE FROM`). Data yang dihapus **hilang permanen** — tidak bisa di-undo, tidak bisa di-audit, dan melanggar compliance.

## 2. Scope

### Tabel yang perlu soft delete:

| # | Tabel | Alasan | Foreign Key Impact |
|---|-------|--------|--------------------|
| 1 | `employees` | Data karyawan harus disimpan untuk payroll history | payrolls, attendances, leave_requests, overtime_requests |
| 2 | `users` | Akun harus disimpan untuk audit trail | employees, audit_logs |
| 3 | `leave_requests` | History cuti untuk compliance | leave_balances |
| 4 | `payrolls` | History gaji **wajib** disimpan min 5 tahun | payroll_items |
| 5 | `overtime_requests` | History lembur untuk compliance | - |
| 6 | `departments` | Referential integrity | employees |
| 7 | `schedules` | Referential integrity | employees |

### Tabel yang TIDAK perlu soft delete:
- `attendances` — time-series data, tidak pernah dihapus
- `leave_types` — menggunakan `is_active` flag (sudah ada)
- `leave_balances` — re-generated setiap tahun
- `audit_logs` — audit trail, tidak boleh dihapus
- `holidays` — calendar data, bisa hard delete
- `companies` — tenant data, NEVER delete

## 3. Implementation Steps

### Step 1: Migration `000014_add_soft_delete.up.sql`

```sql
-- Add deleted_at column to all relevant tables
ALTER TABLE employees ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE leave_requests ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE overtime_requests ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE departments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;

-- Index for performance (partial index — only non-deleted rows)
CREATE INDEX IF NOT EXISTS idx_employees_not_deleted ON employees (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_not_deleted ON users (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_leave_requests_not_deleted ON leave_requests (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_payrolls_not_deleted ON payrolls (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_overtime_not_deleted ON overtime_requests (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_departments_not_deleted ON departments (id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_schedules_not_deleted ON schedules (id) WHERE deleted_at IS NULL;
```

### Step 2: Migration `000014_add_soft_delete.down.sql`

```sql
ALTER TABLE employees DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE leave_requests DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE payrolls DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE overtime_requests DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE departments DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE schedules DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_employees_not_deleted;
DROP INDEX IF EXISTS idx_users_not_deleted;
DROP INDEX IF EXISTS idx_leave_requests_not_deleted;
DROP INDEX IF EXISTS idx_payrolls_not_deleted;
DROP INDEX IF EXISTS idx_overtime_not_deleted;
DROP INDEX IF EXISTS idx_departments_not_deleted;
DROP INDEX IF EXISTS idx_schedules_not_deleted;
```

### Step 3: Add `DeletedAt` to Entity Structs

For each entity, add:
```go
DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
```

**Files to modify:**
- `internal/employee/repository/employee_repository.go` — `Employee` and `EmployeeWithUser` structs
- `internal/user/entity/user.go` — `User` struct
- `internal/leave/entity/leave.go` — `LeaveRequest` struct
- `internal/payroll/entity/payroll.go` — `Payroll` struct
- `internal/overtime/entity/overtime.go` — overtime request struct
- `internal/department/entity/` or `repository/` — Department struct
- `internal/schedule/entity/` or `repository/` — Schedule struct

### Step 4: Update All `DELETE` Operations → Soft Delete

Replace `DELETE FROM table WHERE id = $1` with:
```sql
UPDATE table SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
```

**Files to modify:**
- `internal/employee/repository/employee_repository.go` — `Delete()`
- `internal/user/repository/user_repository_impl.go` — `Delete()`
- `internal/leave/repository/leave_request_repository.go` — if exists
- `internal/department/repository/` — `Delete()`
- `internal/schedule/repository/` — `Delete()`

### Step 5: Update All `SELECT`/`COUNT` Queries — Add `WHERE deleted_at IS NULL`

**CRITICAL**: Every query that reads data must exclude soft-deleted records.

Pattern to apply:
```diff
- WHERE e.id = $1
+ WHERE e.id = $1 AND e.deleted_at IS NULL

- SELECT COUNT(*) FROM employees
+ SELECT COUNT(*) FROM employees WHERE deleted_at IS NULL
```

This affects **every repository** in the codebase that touches the 7 tables.

### Step 6: Update Unique Constraints (If Any)

If there are unique constraints (e.g., `users(company_id, email)`), they need to be modified to exclude soft-deleted rows:

```sql
-- Drop old unique constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_company_id_email_key;
-- Create partial unique index
CREATE UNIQUE INDEX users_company_email_unique ON users (company_id, email) WHERE deleted_at IS NULL;
```

## 4. Files Changed Summary

| # | Category | Files |
|---|----------|-------|
| 1 | Migrations | `000014_add_soft_delete.up.sql`, `000014_add_soft_delete.down.sql` |
| 2 | Entity structs | 7 files (add `DeletedAt` field) |
| 3 | Repository DELETE methods | 5-7 files (change to UPDATE SET deleted_at) |
| 4 | Repository SELECT/COUNT queries | ALL repositories (~15+ queries) |
| 5 | Unique constraints | Migration for any unique indexes |

## 5. Verification

```bash
go build ./...
grep -rn "DELETE FROM" internal/  # Should find 0 hard deletes on soft-delete tables
grep -rn "deleted_at IS NULL" internal/  # Should find matches in all repositories
```
