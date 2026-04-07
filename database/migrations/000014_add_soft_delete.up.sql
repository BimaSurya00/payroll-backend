-- Add deleted_at column to all relevant tables for soft delete functionality
-- Employees: payroll history, attendance, leave, overtime references
ALTER TABLE employees ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_employees_not_deleted ON employees (id) WHERE deleted_at IS NULL;

-- Users: audit trail references
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_users_not_deleted ON users (id) WHERE deleted_at IS NULL;

-- Leave requests: compliance history
ALTER TABLE leave_requests ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_leave_requests_not_deleted ON leave_requests (id) WHERE deleted_at IS NULL;

-- Payrolls: accounting history (5 years retention)
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_payrolls_not_deleted ON payrolls (id) WHERE deleted_at IS NULL;

-- Overtime requests: compliance history
ALTER TABLE overtime_requests ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_overtime_not_deleted ON overtime_requests (id) WHERE deleted_at IS NULL;

-- Departments: referential integrity with employees
ALTER TABLE departments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_departments_not_deleted ON departments (id) WHERE deleted_at IS NULL;

-- Schedules: referential integrity with employees
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_schedules_not_deleted ON schedules (id) WHERE deleted_at IS NULL;
