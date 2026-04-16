ALTER TABLE schedules ADD COLUMN IF NOT EXISTS office_lat DOUBLE PRECISION;
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS office_long DOUBLE PRECISION;
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS allowed_radius_meters INT;

UPDATE schedules s
SET office_lat = c.office_lat,
    office_long = c.office_long,
    allowed_radius_meters = c.allowed_radius_meters
FROM companies c
WHERE s.company_id = c.id;

ALTER TABLE companies DROP COLUMN IF EXISTS office_lat;
ALTER TABLE companies DROP COLUMN IF EXISTS office_long;
ALTER TABLE companies DROP COLUMN IF EXISTS allowed_radius_meters;
