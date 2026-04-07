-- Reverse multi-tenancy migration: drop company_id from all tables

ALTER TABLE employees DROP COLUMN IF EXISTS company_id;
ALTER TABLE schedules DROP COLUMN IF EXISTS company_id;
ALTER TABLE attendances DROP COLUMN IF EXISTS company_id;
ALTER TABLE payrolls DROP COLUMN IF EXISTS company_id;
ALTER TABLE payroll_items DROP COLUMN IF EXISTS company_id;
ALTER TABLE payroll_configs DROP COLUMN IF EXISTS company_id;
ALTER TABLE departments DROP COLUMN IF EXISTS company_id;
ALTER TABLE holidays DROP COLUMN IF EXISTS company_id;
ALTER TABLE leave_types DROP COLUMN IF EXISTS company_id;
ALTER TABLE leave_balances DROP COLUMN IF EXISTS company_id;
ALTER TABLE leave_requests DROP COLUMN IF EXISTS company_id;
ALTER TABLE overtime_policies DROP COLUMN IF EXISTS company_id;
ALTER TABLE overtime_requests DROP COLUMN IF EXISTS company_id;
ALTER TABLE overtime_attendance DROP COLUMN IF EXISTS company_id;
ALTER TABLE audit_logs DROP COLUMN IF EXISTS company_id;

DROP TABLE IF EXISTS companies;
