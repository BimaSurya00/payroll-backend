# MVP-26: Fix Leave Error Handling

## Prioritas: 🟢 IMPROVEMENT
## Estimasi: 30 menit
## Tipe: Bug Fix

---

## Deskripsi Masalah

`createLeaveAttendances()` di leave service men-ignore error dari `attendanceRepo.Create(...)`. Jika attendance gagal dibuat saat leave di-approve, tidak ada log apapun.

## File yang Diubah

### [MODIFY] `internal/leave/service/leave_service.go`

**Di `createLeaveAttendances()`:**
```go
for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
    if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
        continue
    }
    err := s.attendanceRepo.Create(ctx, &attendanceEntity.Attendance{
        // ... fields ...
    })
    if err != nil {
        // Log error instead of ignoring
        zap.L().Error("failed to create leave attendance",
            zap.String("employeeID", employeeID),
            zap.Time("date", date),
            zap.Error(err),
        )
    }
}
```

## Verifikasi

1. `go build ./...`
2. Approve leave → verify no error logs in normal case
3. Simulate failure → verify error is logged (not silently ignored)
