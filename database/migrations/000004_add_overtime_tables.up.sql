-- ============================================
-- Overtime Management Module Migration
-- ============================================

-- 1. Create Overtime Policies Table
CREATE TABLE IF NOT EXISTS overtime_policies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    rate_type VARCHAR(20) NOT NULL DEFAULT 'MULTIPLIER', -- FIXED, MULTIPLIER
    rate_multiplier DECIMAL(4,2) DEFAULT 1.5,
    fixed_amount DECIMAL(10,2),
    min_overtime_minutes INT DEFAULT 60,
    max_overtime_hours_per_day DECIMAL(4,2) DEFAULT 4.0,
    max_overtime_hours_per_month DECIMAL(5,2) DEFAULT 40.0,
    requires_approval BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Default Overtime Policy
INSERT INTO overtime_policies (name, rate_type, rate_multiplier)
VALUES ('Standard Overtime (1.5x)', 'MULTIPLIER', 1.5)
ON CONFLICT DO NOTHING;

-- 2. Create Overtime Requests Table
CREATE TABLE IF NOT EXISTS overtime_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    overtime_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    total_hours DECIMAL(4,2) NOT NULL,
    reason TEXT NOT NULL,
    overtime_policy_id UUID REFERENCES overtime_policies(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    approved_by UUID REFERENCES employees(id),
    approved_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, overtime_date)
);

-- 3. Create Overtime Attendance Table
CREATE TABLE IF NOT EXISTS overtime_attendance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    overtime_request_id UUID NOT NULL REFERENCES overtime_requests(id),
    employee_id UUID NOT NULL REFERENCES employees(id),
    clock_in_time TIMESTAMP,
    clock_out_time TIMESTAMP,
    actual_hours DECIMAL(4,2),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Indexes
CREATE INDEX IF NOT EXISTS idx_overtime_policies_active ON overtime_policies(is_active);
CREATE INDEX IF NOT EXISTS idx_overtime_requests_employee_id ON overtime_requests(employee_id);
CREATE INDEX IF NOT EXISTS idx_overtime_requests_status ON overtime_requests(status);
CREATE INDEX IF NOT EXISTS idx_overtime_requests_date ON overtime_requests(overtime_date);
CREATE INDEX IF NOT EXISTS idx_overtime_attendance_request_id ON overtime_attendance(overtime_request_id);

-- Add Comments
COMMENT ON TABLE overtime_policies IS 'Kebijakan perhitungan lembur';
COMMENT ON TABLE overtime_requests IS 'Pengajuan lembur oleh karyawan';
COMMENT ON TABLE overtime_attendance IS 'Absensi lembur (clock in/out)';

COMMENT ON COLUMN overtime_policies.rate_multiplier IS 'Multiplier untuk hourly rate (1.5x, 2x, 3x)';
COMMENT ON COLUMN overtime_requests.total_hours IS 'Durasi lembur yang diajukan';
COMMENT ON COLUMN overtime_attendance.actual_hours IS 'Durasi lembur aktual berdasarkan clock in/out';
