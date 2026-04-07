-- Migration: Add payroll_configs table
-- Description: Dynamic allowance and deduction configuration

CREATE TABLE IF NOT EXISTS payroll_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'FIXED',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed default configs (migrate from hardcoded values)
INSERT INTO payroll_configs (name, code, type, amount, calculation_type, description) VALUES
    ('Transport Allowance', 'TRANSPORT_ALLOWANCE', 'EARNING', 500000, 'FIXED', 'Tunjangan transport bulanan'),
    ('Meal Allowance', 'MEAL_ALLOWANCE', 'EARNING', 300000, 'FIXED', 'Tunjangan makan bulanan'),
    ('Late Deduction', 'LATE_DEDUCTION', 'DEDUCTION', 50000, 'PER_DAY', 'Potongan keterlambatan per hari'),
    ('Absent Deduction', 'ABSENT_DEDUCTION', 'DEDUCTION', 0, 'PER_DAY', 'Potongan absen per hari (hitung dari gaji harian)');
