-- Remove full_name column from employees table
DROP INDEX IF EXISTS idx_employees_full_name;
ALTER TABLE employees DROP COLUMN IF EXISTS full_name;
