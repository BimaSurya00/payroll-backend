-- Rollback Migration: Remove Employment Status, Job Level, Gender, and Division from Employees

-- Drop indexes
DROP INDEX IF EXISTS idx_employees_employment_status;
DROP INDEX IF EXISTS idx_employees_job_level;
DROP INDEX IF EXISTS idx_employees_division;
DROP INDEX IF EXISTS idx_employees_gender;

-- Drop foreign key constraint if exists
-- ALTER TABLE employees DROP CONSTRAINT IF EXISTS fk_employee_division;

-- Drop columns
ALTER TABLE employees DROP COLUMN IF EXISTS employment_status;
ALTER TABLE employees DROP COLUMN IF EXISTS job_level;
ALTER TABLE employees DROP COLUMN IF EXISTS gender;
ALTER TABLE employees DROP COLUMN IF EXISTS division;

-- Drop divisions table
DROP TABLE IF EXISTS divisions;
