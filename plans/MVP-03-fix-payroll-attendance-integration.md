# MVP-03: Fix Payroll Attendance Integration (Late Days)

## Prioritas: 🔴 CRITICAL — Business Logic Error
## Estimasi: 4 jam
## Tipe: Bug Fix + Feature Enhancement

---

## Deskripsi Masalah

Di `internal/payroll/service/payroll_service_impl.go` line 66-69, saat generate payroll:

```go
allowance, deduction, netSalary := helper.CalculateSalary(
    emp.SalaryBase,
    0, // lateDays - can be fetched from attendance if needed  ← HARDCODED!
)
```

`lateDays` selalu 0, artinya **deduction keterlambatan tidak pernah dihitung**. Late deduction (Rp 50.000/hari) tidak pernah terpotong.

## Solusi

Integrasikan data attendance ke payroll generation — hitung jumlah hari terlambat dari tabel `attendances` berdasarkan periode payroll.

## File yang Diubah

### 1. [MODIFY] `internal/attendance/repository/attendance_repository.go`

**Tambah method di interface `AttendanceRepository`:**
```go
// CountLateByEmployeeAndPeriod menghitung jumlah hari dengan status LATE
// untuk employee tertentu dalam rentang tanggal.
CountLateByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error)

// CountAbsentByEmployeeAndPeriod menghitung jumlah hari ABSENT.
CountAbsentByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error)

// GetAttendanceSummaryByPeriod mengembalikan ringkasan kehadiran.
GetAttendanceSummaryByPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*AttendanceSummary, error)
```

**Tambah struct:**
```go
type AttendanceSummary struct {
    TotalPresent int
    TotalLate    int
    TotalAbsent  int
    TotalLeave   int
    TotalDays    int
}
```

### 2. [MODIFY] `internal/attendance/repository/attendance_repository_impl.go`

**Implementasi method baru:**

```go
func (r *attendanceRepositoryImpl) CountLateByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error) {
    query := `SELECT COUNT(*) FROM attendances 
              WHERE employee_id = $1 AND date >= $2 AND date <= $3 AND status = 'LATE'`
    var count int
    err := r.pool.QueryRow(ctx, query, employeeID, startDate, endDate).Scan(&count)
    return count, err
}

func (r *attendanceRepositoryImpl) CountAbsentByEmployeeAndPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (int, error) {
    query := `SELECT COUNT(*) FROM attendances 
              WHERE employee_id = $1 AND date >= $2 AND date <= $3 AND status = 'ABSENT'`
    var count int
    err := r.pool.QueryRow(ctx, query, employeeID, startDate, endDate).Scan(&count)
    return count, err
}

func (r *attendanceRepositoryImpl) GetAttendanceSummaryByPeriod(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) (*AttendanceSummary, error) {
    query := `SELECT 
                COALESCE(SUM(CASE WHEN status = 'PRESENT' THEN 1 ELSE 0 END), 0) as total_present,
                COALESCE(SUM(CASE WHEN status = 'LATE' THEN 1 ELSE 0 END), 0) as total_late,
                COALESCE(SUM(CASE WHEN status = 'ABSENT' THEN 1 ELSE 0 END), 0) as total_absent,
                COALESCE(SUM(CASE WHEN status = 'LEAVE' THEN 1 ELSE 0 END), 0) as total_leave,
                COUNT(*) as total_days
              FROM attendances 
              WHERE employee_id = $1 AND date >= $2 AND date <= $3`
    
    summary := &AttendanceSummary{}
    err := r.pool.QueryRow(ctx, query, employeeID, startDate, endDate).Scan(
        &summary.TotalPresent, &summary.TotalLate, &summary.TotalAbsent, 
        &summary.TotalLeave, &summary.TotalDays,
    )
    return summary, err
}
```

### 3. [MODIFY] `internal/payroll/routes.go`

**Tambah dependency `attendanceRepo`:**
```go
import (
    // ... existing imports
    attendanceRepo "example.com/hris/internal/attendance/repository"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
    payrollRepo := payrollrepository.NewPayrollRepository(postgresDB.Pool)
    employeeRepo := employeerepository.NewEmployeeRepository(postgresDB.Pool)
    attendRepo := attendanceRepo.NewAttendanceRepository(postgresDB.Pool) // TAMBAH

    payrollService := payrollService.NewPayrollService(payrollRepo, employeeRepo, attendRepo) // UBAH
    // ... rest sama
}
```

### 4. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Ubah struct dan constructor — tambah `attendanceRepo`:**
```go
type payrollServiceImpl struct {
    payrollRepo    payrollrepository.PayrollRepository
    employeeRepo   employeerepository.EmployeeRepository
    attendanceRepo attendanceRepo.AttendanceRepository // TAMBAH
}

func NewPayrollService(
    payrollRepo payrollrepository.PayrollRepository,
    employeeRepo employeerepository.EmployeeRepository,
    attendanceRepo attendanceRepo.AttendanceRepository, // TAMBAH
) PayrollService {
    return &payrollServiceImpl{
        payrollRepo:    payrollRepo,
        employeeRepo:   employeeRepo,
        attendanceRepo: attendanceRepo, // TAMBAH
    }
}
```

**Ubah `GenerateBulk` — hitung lateDays dari attendance (line 64-69):**
```go
for _, emp := range employees {
    // Hitung attendance summary dari database
    summary, err := s.attendanceRepo.GetAttendanceSummaryByPeriod(ctx, emp.ID, periodStart, periodEnd)
    lateDays := 0
    absentDays := 0
    if err == nil && summary != nil {
        lateDays = summary.TotalLate
        absentDays = summary.TotalAbsent
    }

    // Calculate salary dengan data aktual
    allowance, deduction, netSalary := helper.CalculateSalary(
        emp.SalaryBase,
        lateDays,
    )
    
    // ... rest sama (create payroll, items, etc.)
```

**Juga tambah payroll item untuk absent deduction jika diperlukan:**
```go
if absentDays > 0 {
    dailySalary := emp.SalaryBase / 22 // asumsi 22 hari kerja
    absentDeduction := float64(absentDays) * dailySalary
    deduction += absentDeduction
    netSalary -= absentDeduction

    items = append(items, &entity.PayrollItem{
        ID:        uuid.New().String(),
        PayrollID: payrollID,
        Name:      fmt.Sprintf("Absent Deduction (%d days)", absentDays),
        Amount:    absentDeduction,
        Type:      "DEDUCTION",
        CreatedAt: now,
    })
}
```

### 5. [MODIFY] `internal/payroll/service/payroll_service.go`

**Update interface jika ada.**

## Verifikasi

1. `go build ./...` — compile sukses
2. Buat data test:
   - Create employee dengan salary_base = 5.000.000
   - Create attendance records: 3 hari LATE, 1 hari ABSENT dalam periode Januari 2026
3. Generate payroll untuk Januari 2026
4. Verifikasi:
   - Late deduction = 3 × Rp 50.000 = Rp 150.000
   - Payroll items harus mengandung "Late Deduction" dengan amount Rp 150.000
   - Net salary = 5.000.000 + 800.000 (allowance) - 150.000 (late) = 5.650.000
