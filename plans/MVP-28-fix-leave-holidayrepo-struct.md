# MVP-28: Fix Leave Service — Add Missing holidayRepo Struct Field

## Prioritas: 🔴 CRITICAL — Compile Error
## Estimasi: 15 menit
## Tipe: Bug Fix

---

## Deskripsi Masalah

Leave service constructor menerima `holidayRepo` parameter dan mengassign ke struct, tapi **field `holidayRepo` tidak ada di struct definition**. Code juga memanggil `s.holidayRepo.FindByDateRange()` di `CreateLeaveRequest`. Ini menyebabkan **compile error**.

## File yang Diubah

### [MODIFY] `internal/leave/service/leave_service.go`

**Perbaikan 1 — Tambah field ke struct (line ~41-49):**
```diff
 type leaveService struct {
     leaveTypeRepo    leaverepo.LeaveTypeRepository
     leaveBalanceRepo leaverepo.LeaveBalanceRepository
     leaveRequestRepo leaverepo.LeaveRequestRepository
     employeeRepo     employeerepo.EmployeeRepository
     userRepo         userRepo.UserRepository
     attendanceRepo   attendanceRepo.AttendanceRepository
+    holidayRepo      holidayrepo.HolidayRepository
     pool             *pgxpool.Pool
 }
```

**Perbaikan 2 — Tambah import (jika belum ada):**
```go
holidayrepo "example.com/hris/internal/holiday/repository"
```

## Verifikasi

```bash
go build ./internal/leave/...
# Harus: ✅ SUCCESS (bukan compile error)
```
