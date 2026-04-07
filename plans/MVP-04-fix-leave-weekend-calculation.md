# MVP-04: Fix Leave Weekend Calculation

## Prioritas: 🟡 IMPORTANT — Data Accuracy
## Estimasi: 2 jam
## Tipe: Bug Fix

---

## Deskripsi Masalah

Di `internal/leave/service/leave_service.go` line 114:
```go
totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
```

Ini menghitung **semua hari** termasuk Sabtu dan Minggu. Contoh:
- Cuti Senin-Jumat = 5 hari → **benar**
- Cuti Senin-Minggu = 7 hari → **seharusnya 5 hari kerja**

Karyawan kehilangan saldo cuti lebih banyak dari seharusnya.

## Solusi

Buat helper function `CountWorkingDays` yang hanya menghitung hari kerja (Senin-Jumat).
Nanti ketika Holiday/Calendar management sudah ada, function ini bisa di-enhance untuk exclude hari libur.

## File yang Diubah

### 1. [NEW] `shared/helper/workday.go`

```go
package helper

import "time"

// CountWorkingDays menghitung jumlah hari kerja (Senin-Jumat) antara dua tanggal (inklusif).
// TODO: Integrasikan dengan holiday management setelah module tersebut dibuat.
func CountWorkingDays(startDate, endDate time.Time) int {
    if startDate.After(endDate) {
        return 0
    }

    workingDays := 0
    current := startDate

    for !current.After(endDate) {
        weekday := current.Weekday()
        if weekday != time.Saturday && weekday != time.Sunday {
            workingDays++
        }
        current = current.AddDate(0, 0, 1)
    }

    return workingDays
}
```

### 2. [MODIFY] `internal/leave/service/leave_service.go`

**Tambah import:**
```go
sharedHelper "hris/shared/helper"
```

**Ubah line 114 di `CreateLeaveRequest`:**
```go
// SEBELUM:
totalDays := int(endDate.Sub(startDate).Hours()/24) + 1

// SESUDAH:
totalDays := sharedHelper.CountWorkingDays(startDate, endDate)
```

**Tambah validasi — minimal 1 hari kerja:**
```go
if totalDays == 0 {
    return nil, errors.New("selected dates contain no working days")
}
```

### 3. [MODIFY] `internal/leave/service/leave_service.go` — `createLeaveAttendances`

**Ubah method (line 411-427) agar hanya buat attendance record untuk hari kerja:**
```go
func (s *leaveService) createLeaveAttendances(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time, leaveRequestID uuid.UUID) {
    currentDate := startDate
    for !currentDate.After(endDate) {
        // Skip weekend
        weekday := currentDate.Weekday()
        if weekday == time.Saturday || weekday == time.Sunday {
            currentDate = currentDate.AddDate(0, 0, 1)
            continue
        }

        attendance := &leaveattendance.Attendance{
            ID:         uuid.New().String(),
            EmployeeID: employeeID.String(),
            Date:       currentDate,
            Status:     "LEAVE",
            Notes:      fmt.Sprintf("Leave Request ID: %s", leaveRequestID.String()),
        }

        _ = s.attendanceRepo.Create(ctx, attendance)
        currentDate = currentDate.AddDate(0, 0, 1)
    }
}
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Test scenarios:
   - Cuti Senin (2026-02-09) s/d Jumat (2026-02-13) → totalDays = 5
   - Cuti Jumat (2026-02-13) s/d Senin (2026-02-16) → totalDays = 2 (Jumat + Senin)
   - Cuti Sabtu (2026-02-14) s/d Minggu (2026-02-15) → error "no working days"
3. Cek attendance records yang dibuat saat approve leave — harus hanya hari kerja
