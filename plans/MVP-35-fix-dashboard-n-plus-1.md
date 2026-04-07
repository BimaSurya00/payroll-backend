# MVP-35: Fix Dashboard N+1 Query Problem

**Estimasi**: 3 jam  
**Impact**: HIGH — Performance + Reliability  
**Prerequisite**: MVP-32 (compile errors fixed) ✅

---

## 1. Problem Statement

File [`internal/dashboard/service/dashboard_service.go`](file:///home/bima/Documents/example.com/hris/internal/dashboard/service/dashboard_service.go) membuat **14 individual `QueryRow` calls** untuk satu endpoint `GET /api/v1/dashboard/summary`. Setiap call adalah round-trip ke database.

### Masalah Spesifik

| # | Masalah | Lokasi | Impact |
|---|---------|--------|--------|
| 1 | **14 sequential queries** untuk 1 endpoint | `GetSummary()` lines 36–84 | Response time >500ms seiring data tumbuh |
| 2 | **Semua error di-swallow** (`_ = s.pool.QueryRow(...)`) | lines 36–84 | Dashboard return "semua 0" saat DB error, admin salah baca data |
| 3 | **Payroll counts tanpa period filter** | lines 68–70 | `WHERE status = 'DRAFT'` return data ALL TIME, bukan per-period |
| 4 | **Tidak ada caching** | N/A | Setiap page refresh = 14 DB round-trips |

---

## 2. Current Code (SEBELUM)

```go
// dashboard_service.go — GetSummary() 
// PROBLEM: 14 individual queries, all errors ignored

// Attendance: 5 queries
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'PRESENT'`, today).Scan(&todayPresent)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'LATE'`, today).Scan(&todayLate)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'ABSENT'`, today).Scan(&todayAbsent)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM attendances WHERE date = $1 AND status = 'LEAVE'`, today).Scan(&todayLeave)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees`).Scan(&totalEmp)

// Leave: 3 queries
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM leave_requests WHERE status = 'PENDING'`).Scan(&pending)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM leave_requests WHERE status = 'APPROVED' AND created_at >= $1`, monthStart).Scan(&approvedThisMonth)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM leave_requests WHERE status = 'REJECTED' AND created_at >= $1`, monthStart).Scan(&rejectedThisMonth)

// Payroll: 4 queries  
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payrolls WHERE status = 'DRAFT'`).Scan(&draftCount)       // ← NO period filter!
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payrolls WHERE status = 'APPROVED'`).Scan(&approvedCount) // ← NO period filter!
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM payrolls WHERE status = 'PAID'`).Scan(&paidCount)         // ← NO period filter!
_ = s.pool.QueryRow(ctx, `SELECT COALESCE(SUM(net_salary), 0) FROM payrolls WHERE period_start >= $1 AND period_start <= $2`, monthStart, today).Scan(&totalNetSalary)

// Employee: 3 queries
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION')`).Scan(&totalActive)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE employment_status NOT IN ('PERMANENT', 'CONTRACT', 'PROBATION')`).Scan(&totalInactive)
_ = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM employees WHERE created_at >= $1`, monthStart).Scan(&newThisMonth)
```

---

## 3. Target Code (SESUDAH)

### 3.1. Consolidated SQL Queries

14 queries → **4 queries** menggunakan `GROUP BY` dan conditional aggregation:

#### Query 1: Attendance Summary (replaces 5 queries → 1)

```sql
SELECT
    COALESCE(SUM(CASE WHEN a.status = 'PRESENT' THEN 1 ELSE 0 END), 0) AS today_present,
    COALESCE(SUM(CASE WHEN a.status = 'LATE'    THEN 1 ELSE 0 END), 0) AS today_late,
    COALESCE(SUM(CASE WHEN a.status = 'ABSENT'  THEN 1 ELSE 0 END), 0) AS today_absent,
    COALESCE(SUM(CASE WHEN a.status = 'LEAVE'   THEN 1 ELSE 0 END), 0) AS today_leave,
    (SELECT COUNT(*) FROM employees) AS total_employees
FROM attendances a
WHERE a.date = $1
```

#### Query 2: Leave Summary (replaces 3 queries → 1)

```sql
SELECT 
    COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending_requests,
    COALESCE(SUM(CASE WHEN status = 'APPROVED' AND created_at >= $1 THEN 1 ELSE 0 END), 0) AS approved_this_month,
    COALESCE(SUM(CASE WHEN status = 'REJECTED' AND created_at >= $1 THEN 1 ELSE 0 END), 0) AS rejected_this_month
FROM leave_requests
WHERE status IN ('PENDING', 'APPROVED', 'REJECTED')
    AND (status = 'PENDING' OR created_at >= $1)
```

#### Query 3: Payroll Summary (replaces 4 queries → 1)

> **FIX: Add `period_start` filter** agar counts hanya untuk bulan ini

```sql
SELECT 
    COALESCE(SUM(CASE WHEN status = 'DRAFT'    THEN 1 ELSE 0 END), 0) AS draft_count,
    COALESCE(SUM(CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END), 0) AS approved_count,
    COALESCE(SUM(CASE WHEN status = 'PAID'     THEN 1 ELSE 0 END), 0) AS paid_count,
    COALESCE(SUM(net_salary), 0) AS total_net_salary
FROM payrolls
WHERE period_start >= $1 AND period_start <= $2
```

#### Query 4: Employee Summary (replaces 3 queries → 1)

```sql
SELECT 
    COALESCE(SUM(CASE WHEN employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0) AS total_active,
    COALESCE(SUM(CASE WHEN employment_status NOT IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0) AS total_inactive,
    COALESCE(SUM(CASE WHEN created_at >= $1 THEN 1 ELSE 0 END), 0) AS new_this_month
FROM employees
```

---

## 4. Implementation Steps (Step-by-Step)

### Step 1: Rewrite `GetSummary()` in `dashboard_service.go`

**File**: [`internal/dashboard/service/dashboard_service.go`](file:///home/bima/Documents/example.com/hris/internal/dashboard/service/dashboard_service.go)

**Replace the entire `GetSummary` function body** (lines 26–96) with the following implementation:

```go
func (s *dashboardService) GetSummary(ctx context.Context) (*dto.DashboardSummary, error) {
	today := sharedHelper.Today()
	now := sharedHelper.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, sharedHelper.GetLocation())

	// 1. Attendance summary — 1 query (was 5)
	attendanceSummary := dto.AttendanceSummary{}
	err := s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN a.status = 'PRESENT' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN a.status = 'LATE'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN a.status = 'ABSENT'  THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN a.status = 'LEAVE'   THEN 1 ELSE 0 END), 0),
			(SELECT COUNT(*) FROM employees)
		FROM attendances a
		WHERE a.date = $1
	`, today).Scan(
		&attendanceSummary.TodayPresent,
		&attendanceSummary.TodayLate,
		&attendanceSummary.TodayAbsent,
		&attendanceSummary.TodayLeave,
		&attendanceSummary.TotalEmployees,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query attendance summary: %w", err)
	}

	// 2. Leave summary — 1 query (was 3)
	leaveSummary := dto.LeaveSummary{}
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'APPROVED' AND created_at >= $1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'REJECTED' AND created_at >= $1 THEN 1 ELSE 0 END), 0)
		FROM leave_requests
		WHERE status IN ('PENDING', 'APPROVED', 'REJECTED')
			AND (status = 'PENDING' OR created_at >= $1)
	`, monthStart).Scan(
		&leaveSummary.PendingRequests,
		&leaveSummary.ApprovedThisMonth,
		&leaveSummary.RejectedThisMonth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query leave summary: %w", err)
	}

	// 3. Payroll summary — 1 query (was 4)
	// FIX: Added period filter so counts are for THIS MONTH only
	payrollSummary := dto.PayrollSummary{
		CurrentPeriod: today.Format("January 2006"),
	}
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'DRAFT'    THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'PAID'     THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(net_salary), 0)
		FROM payrolls
		WHERE period_start >= $1 AND period_start <= $2
	`, monthStart, today).Scan(
		&payrollSummary.DraftCount,
		&payrollSummary.ApprovedCount,
		&payrollSummary.PaidCount,
		&payrollSummary.TotalNetSalary,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query payroll summary: %w", err)
	}

	// 4. Employee summary — 1 query (was 3)
	employeeSummary := dto.EmployeeSummary{}
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN employment_status IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN employment_status NOT IN ('PERMANENT', 'CONTRACT', 'PROBATION') THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN created_at >= $1 THEN 1 ELSE 0 END), 0)
		FROM employees
	`, monthStart).Scan(
		&employeeSummary.TotalActive,
		&employeeSummary.TotalInactive,
		&employeeSummary.NewThisMonth,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query employee summary: %w", err)
	}

	return &dto.DashboardSummary{
		Attendance: attendanceSummary,
		Leave:      leaveSummary,
		Payroll:    payrollSummary,
		Employee:   employeeSummary,
	}, nil
}
```

### Step 2: Add `fmt` to Imports

The current import block in `dashboard_service.go` is:

```go
import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"example.com/hris/internal/dashboard/dto"
	sharedHelper "example.com/hris/shared/helper"
)
```

**Add `"fmt"` to the import block:**

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"example.com/hris/internal/dashboard/dto"
	sharedHelper "example.com/hris/shared/helper"
)
```

### Step 3: No Changes Needed to Other Files

- **`dto/dashboard_response.go`** — ✅ No changes needed. The same DTO struct works with the new queries.
- **`handler/dashboard_handler.go`** — ✅ No changes needed. It already returns an error response if `GetSummary` returns an error. Previously this could never happen since all errors were swallowed.
- **`routes.go`** — ✅ No changes needed.

---

## 5. Changes Summary

### Files to Modify

| File | What Changes | Lines Affected |
|------|-------------|----------------|
| [`dashboard_service.go`](file:///home/bima/Documents/example.com/hris/internal/dashboard/service/dashboard_service.go) | Rewrite `GetSummary()`: consolidate 14→4 queries, add error handling, add period filter on payroll, add `fmt` import | Lines 1–97 (entire file) |

### Files NOT Changed

| File | Reason |
|------|--------|
| `dto/dashboard_response.go` | Struct is unchanged |
| `handler/dashboard_handler.go` | Already handles errors from service |
| `routes.go` | No dependency changes |

---

## 6. Before vs After Comparison

| Metric | Before (Current) | After (Target) |
|--------|-------------------|----------------|
| DB round-trips per request | **14** | **4** |
| Error handling | **None** (`_ = err`) — silent zeros on failure | **Full** — returns error to caller |
| Payroll period filter | **None** — counts ALL TIME data | **Monthly** — scoped to current month |
| Response time (estimate) | ~50-100ms (14 queries) | ~15-30ms (4 queries) |

---

## 7. Verification Plan

### Build Verification

```bash
cd /home/bima/Documents/example.com/hris
go build ./...
go vet ./internal/dashboard/...
```

### Manual Testing (if DB available)

```bash
# Start the server
make run

# Test the endpoint
curl -s -H "Authorization: Bearer <admin_token>" \
  http://localhost:8080/api/v1/dashboard/summary | jq .
```

**Expected Response Structure** (unchanged from before):
```json
{
  "status": "success",
  "message": "Dashboard summary retrieved",
  "data": {
    "attendance": {
      "todayPresent": 0,
      "todayLate": 0,
      "todayAbsent": 0,
      "todayLeave": 0,
      "totalEmployees": 0
    },
    "leave": {
      "pendingRequests": 0,
      "approvedThisMonth": 0,
      "rejectedThisMonth": 0
    },
    "payroll": {
      "draftCount": 0,
      "approvedCount": 0,
      "paidCount": 0,
      "totalNetSalary": 0,
      "currentPeriod": "February 2026"
    },
    "employee": {
      "totalActive": 0,
      "totalInactive": 0,
      "newThisMonth": 0
    }
  }
}
```

### Edge Case: Attendance Query When No Rows

The attendance query uses a subquery `(SELECT COUNT(*) FROM employees)` inside a `FROM attendances a WHERE a.date = $1`. If there are **zero attendance records** for today, the main `SUM(CASE ...)` returns `NULL` which is handled by `COALESCE`.

**However**, the subquery for `total_employees` is a scalar subquery and always returns a value. But the outer `SELECT` may return **zero rows** if there's no matching attendance row. To handle this correctly, the query wraps everything with `COALESCE` and the `SUM` aggregation ensures at least one row is always returned (aggregate without `GROUP BY` always returns exactly one row, even on empty input).

✅ This is safe — aggregate functions without `GROUP BY` always return exactly 1 row with `NULL` values, which `COALESCE` converts to `0`.

---

## 8. Known Limitations (Out of Scope for This MVP)

These are explicitly **NOT** addressed here:

| Item | Reason | Future MVP |
|------|--------|------------|
| KeyDB caching (30s TTL) | Adds complexity; query consolidation gives sufficient performance boost for now | MVP-53 (Caching Layer) |
| Time-series / trending data | Separate feature requiring new queries and DTOs | MVP-51 (Dashboard Analytics) |
| Multi-tenancy filtering (`company_id`) | Requires MVP-33 completion + migration to filter all queries | After MVP-33 integration |

---

## 9. Rollback Plan

If the new queries produce incorrect results:
1. Revert `dashboard_service.go` to the previous version (the old code is shown in Section 2)
2. Run `go build ./...` to verify
3. The only behavioral difference is error handling — reverting also reverts to silent error swallowing
