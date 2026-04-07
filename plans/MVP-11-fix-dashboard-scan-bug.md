# MVP-11: Fix Dashboard Scan Bug

## Prioritas: 🔴 CRITICAL — Runtime Crash
## Estimasi: 30 menit
## Tipe: Bug Fix
## Dependency: MVP-06 (Dashboard) sudah ada tapi broken

---

## Deskripsi Masalah

Di `internal/dashboard/service/dashboard_service.go`, semua query menggunakan pattern:
```go
todayPresent, _ := s.pool.QueryRow(ctx, ...).Scan(new(int64))
```

**Masalah:** `QueryRow(...).Scan(...)` mengembalikan `error`, bukan value.
Variable `todayPresent` berisi `error` (bukan `int64`), lalu di-cast ke `int(*todayPresent)` → **PANIC at runtime**.

Dashboard endpoint **pasti crash** saat dipanggil karena nil pointer dereference.

## Solusi

Ganti semua pattern `Scan(new(int64))` menjadi proper variable scanning.

## File yang Diubah

### [MODIFY] `internal/dashboard/service/dashboard_service.go`

**Ganti SELURUH isi fungsi `GetSummary` (line 46-121):**

```go
func (s *dashboardService) GetSummary(ctx context.Context) (*dto.DashboardSummary, error) {
	today := sharedHelper.Today()
	now := sharedHelper.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	// 1. Attendance summary hari ini
	attendanceSummary := dto.AttendanceSummary{}

	var todayPresent, todayLate, todayAbsent, todayLeave, totalEmp int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'PRESENT'`, today).Scan(&todayPresent)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'LATE'`, today).Scan(&todayLate)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'ABSENT'`, today).Scan(&todayAbsent)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'LEAVE'`, today).Scan(&todayLeave)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees`).Scan(&totalEmp)

	attendanceSummary.TodayPresent = int(todayPresent)
	attendanceSummary.TodayLate = int(todayLate)
	attendanceSummary.TodayAbsent = int(todayAbsent)
	attendanceSummary.TodayLeave = int(todayLeave)
	attendanceSummary.TotalEmployees = int(totalEmp)

	// 2. Leave summary
	leaveSummary := dto.LeaveSummary{}
	monthStart := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, sharedHelper.GetLocation())

	var pending, approvedThisMonth, rejectedThisMonth int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM leave_requests WHERE status = 'PENDING'`).Scan(&pending)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM leave_requests WHERE status = 'APPROVED' AND created_at >= $1`, monthStart).Scan(&approvedThisMonth)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM leave_requests WHERE status = 'REJECTED' AND created_at >= $1`, monthStart).Scan(&rejectedThisMonth)

	leaveSummary.PendingRequests = int(pending)
	leaveSummary.ApprovedThisMonth = int(approvedThisMonth)
	leaveSummary.RejectedThisMonth = int(rejectedThisMonth)

	// 3. Payroll summary
	payrollSummary := dto.PayrollSummary{
		CurrentPeriod: today.Format("January 2006"),
	}

	var draftCount, approvedCount, paidCount int64
	var totalNetSalary float64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payrolls WHERE status = 'DRAFT'`).Scan(&draftCount)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payrolls WHERE status = 'APPROVED'`).Scan(&approvedCount)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payrolls WHERE status = 'PAID'`).Scan(&paidCount)
	_ = s.pool.QueryRow(ctx, `SELECT COALESCE(SUM(net_salary), 0) FROM payrolls WHERE period_start >= $1 AND period_start <= $2`, monthStart, today).Scan(&totalNetSalary)

	payrollSummary.DraftCount = int(draftCount)
	payrollSummary.ApprovedCount = int(approvedCount)
	payrollSummary.PaidCount = int(paidCount)
	payrollSummary.TotalNetSalary = totalNetSalary

	// 4. Employee summary
	employeeSummary := dto.EmployeeSummary{}

	var totalActive, totalInactive, newThisMonth int64
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION')`).Scan(&totalActive)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE employment_status NOT IN ('PERMANENT', 'CONTRACT', 'PROBATION')`).Scan(&totalInactive)
	_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE created_at >= $1`, monthStart).Scan(&newThisMonth)

	employeeSummary.TotalActive = int(totalActive)
	employeeSummary.TotalInactive = int(totalInactive)
	employeeSummary.NewThisMonth = int(newThisMonth)

	return &dto.DashboardSummary{
		Attendance: attendanceSummary,
		Leave:      leaveSummary,
		Payroll:    payrollSummary,
		Employee:   employeeSummary,
	}, nil
}
```

**Hapus import `uuid` jika tidak digunakan di tempat lain, dan pastikan import `time` ada:**
```go
import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"example.com/hris/internal/dashboard/dto"
	sharedHelper "example.com/hris/shared/helper"
)
```

> **Catatan**: import untuk repository bisa dihapus juga jika service hanya menggunakan `pool` untuk direct query dan tidak memanggil method dari repo.

## Verifikasi

1. `go build ./...` — compile sukses
2. `GET /api/v1/dashboard/summary` — harus return 200 tanpa crash
3. Response harus berisi angka yang valid (bukan nil/null)
4. Cek semua section: attendance, leave, payroll, employee — semua angka ≥ 0
