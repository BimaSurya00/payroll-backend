# MVP-23: Complete Holiday — Integrate into Leave Service

## Prioritas: 🟡 IMPORTANT
## Estimasi: 1 jam
## Tipe: Complete Partial (MVP-18)
## Dependency: MVP-21 (main.go fix)

---

## Deskripsi Masalah

Holiday module sudah lengkap (`internal/holiday/`), `CountWorkingDaysExcluding()` sudah ada di `workday.go`, tapi leave service **masih memanggil `CountWorkingDays()` tanpa holiday exclusion**.

```go
// internal/leave/service/leave_service.go, line 121
totalDays := sharedHelper.CountWorkingDays(startDate, endDate) // ❌ Tidak exclude holiday
```

## File yang Diubah

### 1. [MODIFY] `internal/leave/service/leave_service.go`

**Tambah dependency:**
```go
type leaveService struct {
    // ... existing fields ...
    holidayRepo holidayrepo.HolidayRepository  // ADD
}
```

**Update constructor:**
```go
func NewLeaveService(
    // ... existing params ...
    holidayRepo holidayrepo.HolidayRepository,  // ADD
    pool *pgxpool.Pool,
) LeaveService {
```

**Update `CreateLeaveRequest` (line ~121):**
```go
// Fetch holidays in date range
holidays, _ := s.holidayRepo.FindByDateRange(ctx, startDate, endDate)
holidayMap := make(map[string]bool)
for _, h := range holidays {
    holidayMap[h.Date.Format("2006-01-02")] = true
}
totalDays := sharedHelper.CountWorkingDaysExcluding(startDate, endDate, holidayMap)
```

### 2. [MODIFY] `internal/leave/routes.go`

Pass `holidayRepo` ke `NewLeaveService`:
```go
holidayRepo := holidayrepo.NewHolidayRepository(postgresDB.Pool)
leaveService := service.NewLeaveService(
    // ... existing params ...
    holidayRepo,  // ADD
    postgresDB.Pool,
)
```

## Verifikasi

1. `go build ./...`
2. Buat cuti 1-8 Feb 2026 (8 Feb = Isra Mi'raj)
   - Harus: totalDays = 5 (bukan 6)
3. Buat cuti di minggu tanpa holiday → hasil sama seperti sebelumnya
