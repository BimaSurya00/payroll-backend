-- Migration: Add holidays table
-- Description: Create holidays table for national/company holidays and seed Indonesia 2026 holidays

CREATE TABLE IF NOT EXISTS holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'NATIONAL',
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    year INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_holidays_date ON holidays(date);
CREATE INDEX idx_holidays_year ON holidays(year);
CREATE INDEX idx_holidays_type ON holidays(type);

-- Seed hari libur nasional Indonesia 2026
INSERT INTO holidays (name, date, type, year) VALUES
    ('Tahun Baru', '2026-01-01', 'NATIONAL', 2026),
    ('Isra Mi''raj', '2026-02-08', 'NATIONAL', 2026),
    ('Hari Raya Nyepi', '2026-03-19', 'NATIONAL', 2026),
    ('Wafat Isa Al Masih', '2026-04-03', 'NATIONAL', 2026),
    ('Hari Buruh', '2026-05-01', 'NATIONAL', 2026),
    ('Kenaikan Isa Al Masih', '2026-05-14', 'NATIONAL', 2026),
    ('Hari Lahir Pancasila', '2026-06-01', 'NATIONAL', 2026),
    ('Hari Kemerdekaan RI', '2026-08-17', 'NATIONAL', 2026),
    ('Maulid Nabi', '2026-08-28', 'NATIONAL', 2026),
    ('Natal', '2026-12-25', 'NATIONAL', 2026);
