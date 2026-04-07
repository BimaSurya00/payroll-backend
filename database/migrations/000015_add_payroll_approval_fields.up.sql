-- Add payroll approval tracking fields
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS approved_by UUID REFERENCES users(id);
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ;
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS paid_at TIMESTAMPTZ;
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMPTZ;
ALTER TABLE payrolls ADD COLUMN IF NOT EXISTS notes TEXT DEFAULT '';

-- Add comment to document the workflow
COMMENT ON COLUMN payrolls.approved_by IS 'User who approved this payroll (ADMIN or SUPER_USER)';
COMMENT ON COLUMN payrolls.approved_at IS 'Timestamp when payroll was approved';
COMMENT ON COLUMN payrolls.paid_at IS 'Timestamp when payment was confirmed';
COMMENT ON COLUMN payrolls.cancelled_at IS 'Timestamp when payroll was cancelled (cannot be reversed)';
COMMENT ON COLUMN payrolls.notes IS 'Optional notes for approvals, cancellations, or rejections';
