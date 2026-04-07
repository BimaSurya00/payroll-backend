-- ============================================
-- Leave Management Module Migration
-- ============================================

-- 1. Create Leave Types Table
CREATE TABLE IF NOT EXISTS leave_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    max_days_per_year INT NOT NULL DEFAULT 12,
    is_paid BOOLEAN NOT NULL DEFAULT true,
    requires_approval BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Default Leave Types
INSERT INTO leave_types (name, code, max_days_per_year, is_paid) VALUES
    ('Cuti Tahunan', 'ANNUAL', 12, true),
    ('Cuti Sakit', 'SICK', 14, true),
    ('Cuti Melahirkan', 'MATERNITY', 90, true),
    ('Cuti Khusus', 'SPECIAL', 3, true),
    ('Cuti Tidak Berbayar', 'UNPAID', 30, false)
ON CONFLICT DO NOTHING;

-- 2. Create Leave Balances Table
CREATE TABLE IF NOT EXISTS leave_balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID NOT NULL REFERENCES leave_types(id) ON DELETE CASCADE,
    year INT NOT NULL,
    balance INT NOT NULL DEFAULT 0,
    used INT NOT NULL DEFAULT 0,
    pending INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employee_id, leave_type_id, year)
);

-- 3. Create Leave Requests Table
CREATE TABLE IF NOT EXISTS leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    leave_type_id UUID NOT NULL REFERENCES leave_types(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INT NOT NULL,
    reason TEXT,
    attachment_url VARCHAR(500),
    emergency_contact VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    approved_by UUID REFERENCES employees(id),
    approved_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Indexes for Leave Tables
CREATE INDEX IF NOT EXISTS idx_leave_types_code ON leave_types(code);
CREATE INDEX IF NOT EXISTS idx_leave_types_active ON leave_types(is_active);

CREATE INDEX IF NOT EXISTS idx_leave_balances_employee_id ON leave_balances(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_balances_year ON leave_balances(year);
CREATE INDEX IF NOT EXISTS idx_leave_balances_employee_year ON leave_balances(employee_id, year);

CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests(status);
CREATE INDEX IF NOT EXISTS idx_leave_requests_dates ON leave_requests(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_leave_type_id ON leave_requests(leave_type_id);

-- Add Comment for Documentation
COMMENT ON TABLE leave_types IS 'Master data untuk jenis-jenis cuti';
COMMENT ON TABLE leave_balances IS 'Tracking saldo cuti per employee per tahun';
COMMENT ON TABLE leave_requests IS 'Pengajuan cuti oleh karyawan';

COMMENT ON COLUMN leave_balances.balance IS 'Total alokasi cuti yang tersedia';
COMMENT ON COLUMN leave_balances.used IS 'Total cuti yang sudah digunakan (approved)';
COMMENT ON COLUMN leave_balances.pending IS 'Total cuti yang sedang diajukan (menunggu approval)';

COMMENT ON COLUMN leave_requests.attachment_url IS 'URL dokumen pendukung (surat dokter, dll) di MinIO';
COMMENT ON COLUMN leave_requests.status IS 'PENDING, APPROVED, REJECTED, CANCELLED';
