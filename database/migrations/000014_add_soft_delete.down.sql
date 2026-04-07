-- Remove deleted_at column from all tables (rollback)
ALTER TABLE employees DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE leave_requests DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE payrolls DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE overtime_requests DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE departments DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE schedules DROP COLUMN IF EXISTS deleted_at;

-- Drop partial indexes
DROP INDEX IF EXISTS idx_employees_not_deleted;
DROP INDEX IF EXISTS idx_users_not_deleted;
DROP INDEX IF EXISTS idx_leave_requests_not_deleted;
DROP INDEX IF EXISTS idx_payrolls_not_deleted;
DROP INDEX IF EXISTS idx_overtime_not_deleted;
DROP INDEX IF EXISTS idx_departments_not_deleted;
DROP INDEX IF EXISTS idx_schedules_not_deleted;
