-- Fix department unique constraint to be per-company
-- Drop the global unique constraint and create a composite unique constraint

-- Drop existing unique constraint on code
ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_code_key;

-- Add composite unique constraint: code + company_id
ALTER TABLE departments ADD CONSTRAINT departments_code_company_unique UNIQUE (code, company_id);

-- Also ensure company_id is not null (should already be the case from migration 000012)
ALTER TABLE departments ALTER COLUMN company_id SET NOT NULL;
