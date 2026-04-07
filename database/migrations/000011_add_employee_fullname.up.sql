-- Add full_name column to employees table for denormalized employee name
-- This avoids JOIN with MongoDB users table for every query

ALTER TABLE employees ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);

-- Create index for faster searching
CREATE INDEX IF NOT EXISTS idx_employees_full_name ON employees(full_name);

-- Populate from existing data (temporary fallback)
-- NOTE: This should be updated from actual User data in MongoDB
UPDATE employees
SET full_name = COALESCE(
    -- Try to get from position (temporary)
    position,
    'Unknown'
)
WHERE full_name IS NULL;

COMMENT ON COLUMN employees.full_name IS 'Denormalized full name from users collection (MongoDB)';
