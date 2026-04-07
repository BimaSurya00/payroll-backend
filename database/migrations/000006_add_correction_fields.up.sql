-- Migration: Add correction fields to attendances
-- Description: Add audit trail fields for attendance corrections

ALTER TABLE attendances ADD COLUMN IF NOT EXISTS corrected_by VARCHAR(255);
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS corrected_at TIMESTAMP;
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS correction_note TEXT;
