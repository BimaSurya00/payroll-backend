# MVP-18: Add Holiday/Calendar Management

## Prioritas: 🟡 IMPORTANT — Data Accuracy
## Estimasi: 3 jam
## Tipe: New Feature
## Dependency: MVP-04 (CountWorkingDays) — enhance untuk exclude hari libur

---

## Deskripsi Masalah

Saat ini `CountWorkingDays` hanya exclude Sabtu-Minggu. Tidak ada konsep hari libur nasional.
Dampak:
- Leave calculation menghitung hari libur sebagai hari kerja — karyawan rugi saldo cuti
- Payroll attendance summary tidak aware hari libur
- Tidak ada calendar reference untuk admin

## Solusi

Buat master data `holidays` + enhance `CountWorkingDays` untuk exclude hari libur.

## File yang Diubah

### 1. [NEW] Database migration: `000008_add_holidays_table.up.sql`

```sql
CREATE TABLE IF NOT EXISTS holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'NATIONAL', -- 'NATIONAL', 'COMPANY', 'OPTIONAL'
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,    -- e.g. Hari Kemerdekaan tiap tahun
    year INT,                                        -- NULL if recurring
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_holidays_date ON holidays(date);
CREATE INDEX idx_holidays_year ON holidays(year);

-- Seed hari libur nasional Indonesia 2026 (contoh)
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
```

### 2. [NEW] `internal/holiday/` — Module baru

Structure:
```
internal/holiday/
├── entity/holiday.go
├── repository/holiday_repository.go
├── repository/holiday_repository_impl.go
├── handler/holiday_handler.go
├── dto/holiday_dto.go
└── routes.go
```

**Key interface:**
```go
type HolidayRepository interface {
    Create(ctx context.Context, holiday *entity.Holiday) error
    FindAll(ctx context.Context, year int) ([]*entity.Holiday, error)
    FindByDateRange(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error)
    Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
    Delete(ctx context.Context, id uuid.UUID) error
    IsHoliday(ctx context.Context, date time.Time) (bool, error)
}
```

**Routes:**
```go
// Admin/SuperUser only
holidays.Get("/", holidayHandler.GetAllByYear)    // ?year=2026
holidays.Post("/", holidayHandler.Create)
holidays.Patch("/:id", holidayHandler.Update)
holidays.Delete("/:id", holidayHandler.Delete)
```

### 3. [MODIFY] `shared/helper/workday.go`

**Enhance `CountWorkingDays` untuk menerima list holidays:**
```go
// CountWorkingDays menghitung hari kerja (Senin-Jumat) yang bukan hari libur.
func CountWorkingDays(startDate, endDate time.Time) int {
    // Keep existing implementation as fallback
    return CountWorkingDaysExcluding(startDate, endDate, nil)
}

// CountWorkingDaysExcluding menghitung hari kerja excluding holidays.
func CountWorkingDaysExcluding(startDate, endDate time.Time, holidays map[string]bool) int {
    if startDate.After(endDate) { return 0 }

    workingDays := 0
    current := startDate
    for !current.After(endDate) {
        weekday := current.Weekday()
        if weekday != time.Saturday && weekday != time.Sunday {
            dateStr := current.Format("2006-01-02")
            if holidays == nil || !holidays[dateStr] {
                workingDays++
            }
        }
        current = current.AddDate(0, 0, 1)
    }
    return workingDays
}
```

### 4. [MODIFY] `internal/leave/service/leave_service.go`

**Update `CreateLeaveRequest` — fetch holidays lalu pakai `CountWorkingDaysExcluding`:**
```go
// Fetch holidays in date range
holidays, _ := s.holidayRepo.FindByDateRange(ctx, startDate, endDate)
holidayMap := make(map[string]bool)
for _, h := range holidays {
    holidayMap[h.Date.Format("2006-01-02")] = true
}

totalDays := sharedHelper.CountWorkingDaysExcluding(startDate, endDate, holidayMap)
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Run migration — seed holidays
3. Test: buat cuti 1-8 Feb 2026 (8 Februari = Isra Mi'raj)
   - Seharusnya total days = 5 (1 weekend Sab-Min + 1 holiday)
4. Admin CRUD holidays: create, update, delete → semua 200 OK
5. `GET /api/v1/holidays?year=2026` → return list holidays
