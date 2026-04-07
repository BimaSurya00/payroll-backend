-- Migration: Add Employment Status, Job Level, Gender, and Division to Employees
-- Created: 2026-02-10

-- Add employment_status column
ALTER TABLE employees 
ADD COLUMN IF NOT EXISTS employment_status VARCHAR(20) DEFAULT 'PROBATION' 
CHECK (employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION'));

-- Add job_level column
ALTER TABLE employees 
ADD COLUMN IF NOT EXISTS job_level VARCHAR(20) DEFAULT 'STAFF' 
CHECK (job_level IN ('CEO', 'MANAGER', 'SUPERVISOR', 'STAFF'));

-- Add gender column
ALTER TABLE employees 
ADD COLUMN IF NOT EXISTS gender VARCHAR(10) 
CHECK (gender IN ('MALE', 'FEMALE'));

-- Add division column
ALTER TABLE employees 
ADD COLUMN IF NOT EXISTS division VARCHAR(100) DEFAULT 'GENERAL';

-- Add comments for documentation
COMMENT ON COLUMN employees.employment_status IS 'Employment status: PERMANENT, CONTRACT, PROBATION';
COMMENT ON COLUMN employees.job_level IS 'Job level: CEO, MANAGER, SUPERVISOR, STAFF';
COMMENT ON COLUMN employees.gender IS 'Gender: MALE, FEMALE';
COMMENT ON COLUMN employees.division IS 'Department/Division: IT, HR, FINANCE, MARKETING, OPERATIONS, GENERAL';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_employees_employment_status ON employees(employment_status);
CREATE INDEX IF NOT EXISTS idx_employees_job_level ON employees(job_level);
CREATE INDEX IF NOT EXISTS idx_employees_division ON employees(division);
CREATE INDEX IF NOT EXISTS idx_employees_gender ON employees(gender);

-- Create a separate divisions table for reference
CREATE TABLE IF NOT EXISTS divisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default divisions
INSERT INTO divisions (name, code, description) VALUES
('Information Technology', 'IT', 'IT and Software Development'),
('Human Resources', 'HR', 'HR and Recruitment'),
('Finance', 'FIN', 'Finance and Accounting'),
('Marketing', 'MKT', 'Marketing and Sales'),
('Operations', 'OPS', 'Operations and Logistics'),
('General', 'GEN', 'General Administration')
ON CONFLICT (name) DO NOTHING;

-- Add foreign key constraint for division (optional)
-- ALTER TABLE employees 
-- ADD CONSTRAINT fk_employee_division 
-- FOREIGN KEY (division) REFERENCES divisions(name);

-- Update existing records to have default values
UPDATE employees 
SET 
    employment_status = 'PERMANENT',
    job_level = 'STAFF',
    division = 'GENERAL'
WHERE employment_status IS NULL 
   OR job_level IS NULL 
   OR division IS NULL;
