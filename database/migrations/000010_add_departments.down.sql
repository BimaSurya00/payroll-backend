-- Rollback: Remove departments table

ALTER TABLE employees DROP COLUMN IF EXISTS department_id;
DROP TABLE IF EXISTS departments;
