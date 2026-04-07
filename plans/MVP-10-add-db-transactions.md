# MVP-10: Add DB Transaction for Critical Operations

## Prioritas: 🔴 CRITICAL — Data Integrity
## Estimasi: 4 jam
## Tipe: Bug Fix / Architecture

---

## Deskripsi Masalah

Beberapa operasi kritis **tidak menggunakan database transaction**:

1. **Payroll Generate Bulk** — Insert payroll + payroll items per employee dalam loop. Jika gagal di tengah, data sudah setengah masuk.
2. **Leave Approve** — Update balance (pending → used) + update status + create attendance records. Jika salah satu gagal, data inconsistent.
3. **Leave Request Create** — Add to pending balance + create request. Jika create request gagal, balance sudah terupdate.

## Solusi

Gunakan PostgreSQL transaction (`pgx` transaction) untuk membungkus operasi yang harus atomik.

## File yang Diubah

### 1. [MODIFY] `database/postgres.go`

**Tambah method untuk mendapatkan pool (jika belum exposed):**

Pastikan `Pool` dari `database.Postgres` bisa diakses untuk memulai transaction.
Pool sudah exposed sebagai `postgresDB.Pool` (tipe `*pgxpool.Pool`), jadi kita bisa langsung pakai `pool.Begin(ctx)`.

### 2. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Ubah `GenerateBulk` untuk menggunakan transaction.**

Karena service layer saat ini menerima repository (yang sudah ter-bind ke pool), kita perlu menambahkan pool ke service agar bisa begin transaction.

**Opsi A — Tambah `pool` di service struct (pragmatis):**

```go
import "github.com/jackc/pgx/v5/pgxpool"

type payrollServiceImpl struct {
    payrollRepo    payrollrepository.PayrollRepository
    employeeRepo   employeerepository.EmployeeRepository
    attendanceRepo attendanceRepo.AttendanceRepository
    pool           *pgxpool.Pool // TAMBAH
}

func NewPayrollService(
    payrollRepo payrollrepository.PayrollRepository,
    employeeRepo employeerepository.EmployeeRepository,
    attendanceRepo attendanceRepo.AttendanceRepository,
    pool *pgxpool.Pool, // TAMBAH
) PayrollService {
    return &payrollServiceImpl{
        payrollRepo:    payrollRepo,
        employeeRepo:   employeeRepo,
        attendanceRepo: attendanceRepo,
        pool:           pool,
    }
}
```

**Ubah `GenerateBulk`:**
```go
func (s *payrollServiceImpl) GenerateBulk(ctx context.Context, req *dto.GeneratePayrollRequest) (*dto.GeneratePayrollResponse, error) {
    periodStart, periodEnd := helper.GetPeriodRange(req.PeriodMonth, req.PeriodYear)
    periodStartStr := periodStart.Format("2006-01-02")
    periodEndStr := periodEnd.Format("2006-01-02")

    // Check existing
    existingPayrolls, err := s.payrollRepo.FindByPeriod(ctx, periodStartStr, periodEndStr)
    if err == nil && len(existingPayrolls) > 0 {
        return nil, ErrPayrollAlreadyExists
    }

    // Fetch employees
    employees, err := s.employeeRepo.FindAllWithoutPagination(ctx)
    if err != nil || len(employees) == 0 {
        return nil, ErrNoEmployeesFound
    }

    // === BEGIN TRANSACTION ===
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Rollback jika belum commit

    generatedCount := 0
    now := time.Now()

    for _, emp := range employees {
        // Hitung attendance summary
        summary, _ := s.attendanceRepo.GetAttendanceSummaryByPeriod(ctx, emp.ID, periodStart, periodEnd)
        lateDays := 0
        if summary != nil {
            lateDays = summary.TotalLate
        }

        allowance, deduction, netSalary := helper.CalculateSalary(emp.SalaryBase, lateDays)

        payrollID := uuid.New().String()
        payroll := &entity.Payroll{
            ID: payrollID, EmployeeID: emp.ID.String(),
            PeriodStart: periodStart, PeriodEnd: periodEnd,
            BaseSalary: emp.SalaryBase, TotalAllowance: allowance,
            TotalDeduction: deduction, NetSalary: netSalary,
            Status: "DRAFT", GeneratedAt: now, CreatedAt: now, UpdatedAt: now,
        }

        items := []*entity.PayrollItem{
            {ID: uuid.New().String(), PayrollID: payrollID, Name: "Transport Allowance", Amount: helper.TransportAllowance, Type: "EARNING", CreatedAt: now},
            {ID: uuid.New().String(), PayrollID: payrollID, Name: "Meal Allowance", Amount: helper.MealAllowance, Type: "EARNING", CreatedAt: now},
        }
        if deduction > 0 {
            items = append(items, &entity.PayrollItem{
                ID: uuid.New().String(), PayrollID: payrollID,
                Name: "Late Deduction", Amount: deduction, Type: "DEDUCTION", CreatedAt: now,
            })
        }

        // INSERT menggunakan transaction (tx, bukan pool)
        if err := s.insertPayrollWithTx(ctx, tx, payroll, items); err != nil {
            return nil, fmt.Errorf("%w: %v", ErrGenerateFailed, err)
        }

        generatedCount++
    }

    // === COMMIT TRANSACTION ===
    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

    return &dto.GeneratePayrollResponse{
        TotalGenerated: generatedCount,
        PeriodStart: periodStartStr, PeriodEnd: periodEndStr,
        Message: fmt.Sprintf("Successfully generated %d payrolls", generatedCount),
    }, nil
}

// insertPayrollWithTx inserts payroll dan items menggunakan transaction
func (s *payrollServiceImpl) insertPayrollWithTx(ctx context.Context, tx pgx.Tx, payroll *entity.Payroll, items []*entity.PayrollItem) error {
    // Insert payroll
    _, err := tx.Exec(ctx,
        `INSERT INTO payrolls (id, employee_id, period_start, period_end, base_salary, 
         total_allowance, total_deduction, net_salary, status, generated_at, created_at, updated_at) 
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
        payroll.ID, payroll.EmployeeID, payroll.PeriodStart, payroll.PeriodEnd,
        payroll.BaseSalary, payroll.TotalAllowance, payroll.TotalDeduction,
        payroll.NetSalary, payroll.Status, payroll.GeneratedAt, payroll.CreatedAt, payroll.UpdatedAt,
    )
    if err != nil {
        return err
    }

    // Insert items
    for _, item := range items {
        _, err := tx.Exec(ctx,
            `INSERT INTO payroll_items (id, payroll_id, name, amount, type) VALUES ($1,$2,$3,$4,$5)`,
            item.ID, item.PayrollID, item.Name, item.Amount, item.Type,
        )
        if err != nil {
            return err
        }
    }

    return nil
}
```

### 3. [MODIFY] `internal/payroll/routes.go`

**Pass pool ke service:**
```go
payrollService := payrollService.NewPayrollService(payrollRepo, employeeRepo, attendRepo, postgresDB.Pool)
```

### 4. [MODIFY] `internal/leave/service/leave_service.go` — Wrap Leave Approve

**Sama konsepnya — tambah `pool` ke leaveService struct:**

```go
type leaveService struct {
    // ... existing fields
    pool *pgxpool.Pool // TAMBAH
}
```

**Wrap `ApproveLeaveRequest` dalam transaction:**
```go
func (s *leaveService) ApproveLeaveRequest(ctx context.Context, id, approverID string, req *dto.ApproveLeaveRequest) (*dto.LeaveRequestResponse, error) {
    // ... validation code yang sudah ada ...

    // === BEGIN TRANSACTION ===
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Move balance from pending to used (gunakan tx)
    // Update request status (gunakan tx)
    
    // === COMMIT ===
    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("failed to commit: %w", err)
    }

    // Create attendance records (di luar transaction, non-critical)
    s.createLeaveAttendances(ctx, ...)

    // return updated response
}
```

> **Catatan untuk agent**: Pattern yang sama bisa di-apply ke `CreateLeaveRequest` (add pending + create request) dan `RejectLeaveRequest` (return pending + update status).

### 5. Update routes yang terlibat untuk pass `pool`

Sesuaikan constructor di routes.go yang berubah signature.

## Verifikasi

1. `go build ./...` — compile sukses
2. Test payroll generate:
   - Simulate error di tengah insert (misalnya constraint violation) → semua payroll yang sudah diinsert harus rollback
   - Generate normal → semua payroll harus masuk semua
3. Test leave approve:
   - Approve leave — balance dan status harus konsisten
   - Jika update status gagal, balance tidak boleh berubah
