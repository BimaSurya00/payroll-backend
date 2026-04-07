-- Migration: Add departments table
-- Description: Create departments master data table and migrate existing division data

CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    head_employee_id UUID REFERENCES employees(id),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_departments_code ON departments(code);
CREATE INDEX idx_departments_active ON departments(is_active);

-- Migrate existing divisions to departments
INSERT INTO departments (name, code)
SELECT DISTINCT division, UPPER(REPLACE(division, ' ', '_'))
FROM employees
WHERE division IS NOT NULL AND division != ''
ON CONFLICT (code) DO NOTHING;

-- Add department_id column to employees
ALTER TABLE employees ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id);

-- Populate department_id from existing division data
UPDATE employees e
SET department_id = d.id
FROM departments d
WHERE d.name = e.division;
