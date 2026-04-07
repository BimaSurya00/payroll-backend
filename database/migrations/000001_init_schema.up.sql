-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Note: Users table is stored in MongoDB, not PostgreSQL
-- Migration below creates all PostgreSQL tables for the application

-- 1. Create Schedules Table
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    time_in TIME NOT NULL,
    time_out TIME NOT NULL,
    allowed_late_minutes INT DEFAULT 15,
    office_lat FLOAT NOT NULL,
    office_long FLOAT NOT NULL,
    allowed_radius_meters INT DEFAULT 50,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Default Schedule (09:00 - 17:00)
INSERT INTO schedules (name, time_in, time_out, office_lat, office_long)
VALUES ('Regular Office Hour', '09:00', '17:00', -6.200000, 106.816666) -- Default Jakarta coords, update later
ON CONFLICT DO NOTHING;

-- 3. Create Employees Table
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE, -- One-to-One with Users
    position VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    address TEXT,
    salary_base DECIMAL(10, 2) NOT NULL DEFAULT 0,
    join_date DATE NOT NULL DEFAULT CURRENT_DATE,
    bank_name VARCHAR(50),
    bank_account_number VARCHAR(50),
    bank_account_holder VARCHAR(100),
    schedule_id UUID REFERENCES schedules(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create Attendances Table
CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    schedule_id UUID REFERENCES schedules(id), -- Snapshot of schedule
    date DATE NOT NULL,
    clock_in_time TIMESTAMP,
    clock_out_time TIMESTAMP,
    clock_in_lat FLOAT,
    clock_in_long FLOAT,
    clock_out_lat FLOAT,
    clock_out_long FLOAT,
    status VARCHAR(20) NOT NULL, -- PRESENT, LATE, ABSENT, LEAVE
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, date)
);

-- 5. Create Payrolls Table
CREATE TABLE IF NOT EXISTS payrolls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    base_salary DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_allowance DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_deduction DECIMAL(10, 2) NOT NULL DEFAULT 0,
    net_salary DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, APPROVED, PAID
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Create Payroll Items Table (Details)
CREATE TABLE IF NOT EXISTS payroll_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payroll_id UUID NOT NULL REFERENCES payrolls(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    type VARCHAR(20) NOT NULL -- EARNING, DEDUCTION
);

-- Create Indexes
CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances(date);
CREATE INDEX IF NOT EXISTS idx_attendances_employee_id ON attendances(employee_id);
CREATE INDEX IF NOT EXISTS idx_payrolls_period ON payrolls(period_start, period_end);
