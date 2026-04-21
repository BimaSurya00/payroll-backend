ALTER TABLE overtime_requests DROP CONSTRAINT IF EXISTS overtime_requests_approved_by_fkey;
ALTER TABLE overtime_requests ALTER COLUMN approved_by DROP NOT NULL;
