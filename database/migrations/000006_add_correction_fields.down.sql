-- Rollback: Remove correction fields from attendances

ALTER TABLE attendances DROP COLUMN IF EXISTS corrected_by;
ALTER TABLE attendances DROP COLUMN IF EXISTS corrected_at;
ALTER TABLE attendances DROP COLUMN IF EXISTS correction_note;
