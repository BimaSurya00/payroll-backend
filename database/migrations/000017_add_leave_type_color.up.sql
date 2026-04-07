-- Migration: Add color and default_days columns to leave_types table
-- Description: Add color and default_days fields for UI customization

ALTER TABLE leave_types ADD COLUMN IF NOT EXISTS color VARCHAR(7) DEFAULT '#3b82f6';
ALTER TABLE leave_types ADD COLUMN IF NOT EXISTS default_days INT DEFAULT 0;

-- Update existing records with default values
UPDATE leave_types SET color = '#3b82f6' WHERE color IS NULL;
UPDATE leave_types SET default_days = max_days_per_year WHERE default_days IS NULL OR default_days = 0;
