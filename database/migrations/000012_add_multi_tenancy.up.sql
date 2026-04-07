-- ============================================
-- Multi-Tenancy Foundation Migration
-- Adds company_id to all tenant-scoped tables
-- ============================================

-- 1. Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    plan VARCHAR(50) NOT NULL DEFAULT 'free',
    max_employees INT NOT NULL DEFAULT 25,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_companies_slug ON companies(slug);

-- 2. Insert default company for existing data migration
INSERT INTO companies (id, name, slug, plan, max_employees)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Company', 'default', 'pro', 9999)
ON CONFLICT DO NOTHING;

-- ============================================
-- 3. Add company_id to all tenant-scoped tables
-- Pattern: ADD COLUMN → DEFAULT backfill → SET NOT NULL → INDEX
-- ============================================

-- employees
ALTER TABLE employees ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE employees SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE employees ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_employees_company ON employees(company_id);

-- schedules
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE schedules SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE schedules ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_schedules_company ON schedules(company_id);

-- attendances
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE attendances SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE attendances ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_attendances_company ON attendances(company_id);

-- payrolls
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE payrolls SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE payrolls ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_payrolls_company ON payrolls(company_id);

-- payroll_items
ALTER TABLE payroll_items ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE payroll_items SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE payroll_items ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_payroll_items_company ON payroll_items(company_id);

-- payroll_configs
ALTER TABLE payroll_configs ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE payroll_configs SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE payroll_configs ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_payroll_configs_company ON payroll_configs(company_id);

-- departments
ALTER TABLE departments ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE departments SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE departments ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_departments_company ON departments(company_id);

-- holidays
ALTER TABLE holidays ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE holidays SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE holidays ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_holidays_company ON holidays(company_id);

-- leave_types
ALTER TABLE leave_types ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE leave_types SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE leave_types ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_leave_types_company ON leave_types(company_id);

-- leave_balances
ALTER TABLE leave_balances ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE leave_balances SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE leave_balances ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_leave_balances_company ON leave_balances(company_id);

-- leave_requests
ALTER TABLE leave_requests ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE leave_requests SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE leave_requests ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_leave_requests_company ON leave_requests(company_id);

-- overtime_policies
ALTER TABLE overtime_policies ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE overtime_policies SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE overtime_policies ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_overtime_policies_company ON overtime_policies(company_id);

-- overtime_requests
ALTER TABLE overtime_requests ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE overtime_requests SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE overtime_requests ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_overtime_requests_company ON overtime_requests(company_id);

-- overtime_attendance
ALTER TABLE overtime_attendance ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE overtime_attendance SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE overtime_attendance ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_overtime_attendance_company ON overtime_attendance(company_id);

-- audit_logs
ALTER TABLE audit_logs ADD COLUMN IF NOT EXISTS company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';
UPDATE audit_logs SET company_id = '00000000-0000-0000-0000-000000000001' WHERE company_id IS NULL;
ALTER TABLE audit_logs ALTER COLUMN company_id SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_logs_company ON audit_logs(company_id);

-- ============================================
-- Comments
-- ============================================
COMMENT ON TABLE companies IS 'Tenant companies for multi-tenancy SaaS';
COMMENT ON COLUMN companies.slug IS 'URL-safe unique identifier for the company';
COMMENT ON COLUMN companies.plan IS 'Subscription plan: free, starter, pro, enterprise';
COMMENT ON COLUMN companies.max_employees IS 'Max employees allowed by plan';
