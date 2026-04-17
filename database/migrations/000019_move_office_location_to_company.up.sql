ALTER TABLE companies ADD COLUMN IF NOT EXISTS office_lat DOUBLE PRECISION;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS office_long DOUBLE PRECISION;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS allowed_radius_meters INT;

INSERT INTO companies (id, name, slug, is_active, plan, max_employees, office_lat, office_long, allowed_radius_meters)
SELECT DISTINCT c.id, c.name, c.slug, c.is_active, c.plan, c.max_employees, s.office_lat, s.office_long, s.allowed_radius_meters
FROM companies c
INNER JOIN schedules s ON s.company_id = c.id
ON CONFLICT (id) DO UPDATE SET
    office_lat = EXCLUDED.office_lat,
    office_long = EXCLUDED.office_long,
    allowed_radius_meters = EXCLUDED.allowed_radius_meters;

ALTER TABLE schedules DROP COLUMN IF EXISTS office_lat;
ALTER TABLE schedules DROP COLUMN IF EXISTS office_long;
ALTER TABLE schedules DROP COLUMN IF EXISTS allowed_radius_meters;
