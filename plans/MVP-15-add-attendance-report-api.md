# MVP-15: Add Attendance Report/Summary API

## Prioritas: 🟡 IMPORTANT — Core Feature
## Estimasi: 4 jam
## Tipe: New Feature

---

## Deskripsi Masalah

Tidak ada endpoint untuk melihat rekap kehadiran bulanan per karyawan. Admin harus menghitung manual.
Ini fitur yang paling sering diakses oleh admin UMKM.

Endpoint yang dibutuhkan:
1. **Monthly summary per employee** — berapa hari hadir, telat, absent, cuti dalam bulan X
2. **Monthly summary all employees** — overview semua karyawan dalam bulan X
3. **My attendance summary** — karyawan lihat rekap kehadiran sendiri

## File yang Diubah

### 1. [NEW] `internal/attendance/dto/report_response.go`

```go
package dto

type AttendanceReportItem struct {
    EmployeeID   string  `json:"employeeId"`
    EmployeeName string  `json:"employeeName"`
    Position     string  `json:"position"`
    Division     string  `json:"division"`
    TotalPresent int     `json:"totalPresent"`
    TotalLate    int     `json:"totalLate"`
    TotalAbsent  int     `json:"totalAbsent"`
    TotalLeave   int     `json:"totalLeave"`
    TotalDays    int     `json:"totalDays"`
    AttendanceRate float64 `json:"attendanceRate"` // Persentase kehadiran
    Period       string  `json:"period"`
}

type MonthlyAttendanceReport struct {
    Period     string                `json:"period"`
    Month      int                   `json:"month"`
    Year       int                   `json:"year"`
    TotalEmployees int               `json:"totalEmployees"`
    Summary    AttendanceReportItem  `json:"summary"` // Agregat semua employee
    Items      []AttendanceReportItem `json:"items"`
}

type MyAttendanceSummary struct {
    Period       string `json:"period"`
    TotalPresent int    `json:"totalPresent"`
    TotalLate    int    `json:"totalLate"`
    TotalAbsent  int    `json:"totalAbsent"`
    TotalLeave   int    `json:"totalLeave"`
    TotalDays    int    `json:"totalDays"`
    AttendanceRate float64 `json:"attendanceRate"`
}
```

### 2. [MODIFY] `internal/attendance/repository/attendance_repository.go`

**Tambah method di interface:**
```go
// GetMonthlySummaryAll mengembalikan ringkasan kehadiran semua karyawan dalam bulan tertentu
GetMonthlySummaryAll(ctx context.Context, startDate, endDate time.Time) ([]AttendanceMonthlySummary, error)
```

**Tambah struct:**
```go
type AttendanceMonthlySummary struct {
    EmployeeID   uuid.UUID
    TotalPresent int
    TotalLate    int
    TotalAbsent  int
    TotalLeave   int
    TotalDays    int
}
```

### 3. [MODIFY] Implementasi attendance repository

```go
func (r *attendanceRepositoryImpl) GetMonthlySummaryAll(ctx context.Context, startDate, endDate time.Time) ([]AttendanceMonthlySummary, error) {
    query := `SELECT employee_id,
              COALESCE(SUM(CASE WHEN status = 'PRESENT' THEN 1 ELSE 0 END), 0),
              COALESCE(SUM(CASE WHEN status = 'LATE' THEN 1 ELSE 0 END), 0),
              COALESCE(SUM(CASE WHEN status = 'ABSENT' THEN 1 ELSE 0 END), 0),
              COALESCE(SUM(CASE WHEN status = 'LEAVE' THEN 1 ELSE 0 END), 0),
              COUNT(*)
              FROM attendances
              WHERE date >= $1 AND date <= $2
              GROUP BY employee_id`

    rows, err := r.pool.Query(ctx, query, startDate, endDate)
    if err != nil { return nil, err }
    defer rows.Close()

    var summaries []AttendanceMonthlySummary
    for rows.Next() {
        s := AttendanceMonthlySummary{}
        rows.Scan(&s.EmployeeID, &s.TotalPresent, &s.TotalLate,
            &s.TotalAbsent, &s.TotalLeave, &s.TotalDays)
        summaries = append(summaries, s)
    }
    return summaries, nil
}
```

### 4. [MODIFY] `internal/attendance/service/` — Tambah report methods

**Tambah method di AttendanceService interface:**
```go
GetMonthlyReport(ctx context.Context, month, year int) (*dto.MonthlyAttendanceReport, error)
GetMyMonthlySummary(ctx context.Context, userID string, month, year int) (*dto.MyAttendanceSummary, error)
```

### 5. [MODIFY] `internal/attendance/handler/` — Tambah report handlers

```go
func (h *AttendanceHandler) GetMonthlyReport(c *fiber.Ctx) error {
    month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(time.Now().Month()))))
    year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(time.Now().Year())))
    // Validate month 1-12, year reasonable range
    // Call service, return response
}

func (h *AttendanceHandler) GetMyMonthlySummary(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)
    month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(time.Now().Month()))))
    year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(time.Now().Year())))
    // Call service, return response
}
```

### 6. [MODIFY] `internal/attendance/routes.go`

**Tambah routes:**
```go
// Report endpoints
api.Get("/report/monthly", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    attendanceHandler.GetMonthlyReport)

api.Get("/report/my",
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    attendanceHandler.GetMyMonthlySummary)
```

## Verifikasi

1. `go build ./...` — compile sukses
2. `GET /api/v1/attendances/report/monthly?month=2&year=2026` (admin) → return report per employee
3. `GET /api/v1/attendances/report/my?month=2&year=2026` (user) → return summary sendiri
4. Cek attendance rate calculation benar: (present+late) / totalDays × 100
