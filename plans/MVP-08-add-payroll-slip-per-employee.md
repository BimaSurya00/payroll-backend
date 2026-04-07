# MVP-08: Add Payroll Slip View per Employee

## Prioritas: 🟡 IMPORTANT — Essential UX
## Estimasi: 2 jam
## Tipe: New Feature

---

## Deskripsi Masalah

Karyawan **tidak bisa melihat slip gaji sendiri**. Semua payroll endpoint hanya untuk admin.
Ini fitur fundamental — setiap karyawan harus bisa melihat rincian gaji mereka.

## Solusi

Tambah endpoint `/api/v1/payrolls/my` yang menampilkan payroll history milik user yang login.

## File yang Diubah

### 1. [MODIFY] `internal/payroll/service/payroll_service.go`

**Tambah method di interface:**
```go
type PayrollService interface {
    // ... existing methods

    // Self-service - karyawan lihat payroll sendiri
    GetMyPayrolls(ctx context.Context, userID string, page, perPage int, path string) (*helper.PayrollPagination, error)
    GetMyPayrollByID(ctx context.Context, userID string, payrollID string) (*dto.PayrollResponse, error)
}
```

### 2. [MODIFY] `internal/payroll/repository/payroll_repository.go`

**Tambah method di interface (jika belum ada):**
```go
FindByEmployeeID(ctx context.Context, employeeID string, skip, limit int64) ([]*entity.Payroll, error)
CountByEmployeeID(ctx context.Context, employeeID string) (int64, error)
```

### 3. [MODIFY] Implementasi payroll repository

```go
func (r *payrollRepositoryImpl) FindByEmployeeID(ctx context.Context, employeeID string, skip, limit int64) ([]*entity.Payroll, error) {
    query := `SELECT id, employee_id, period_start, period_end, base_salary, 
              total_allowance, total_deduction, net_salary, status, generated_at, created_at, updated_at
              FROM payrolls WHERE employee_id = $1 
              ORDER BY period_start DESC LIMIT $2 OFFSET $3`
    
    rows, err := r.pool.Query(ctx, query, employeeID, limit, skip)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var payrolls []*entity.Payroll
    for rows.Next() {
        p := &entity.Payroll{}
        err := rows.Scan(&p.ID, &p.EmployeeID, &p.PeriodStart, &p.PeriodEnd,
            &p.BaseSalary, &p.TotalAllowance, &p.TotalDeduction, &p.NetSalary,
            &p.Status, &p.GeneratedAt, &p.CreatedAt, &p.UpdatedAt)
        if err != nil {
            return nil, err
        }
        payrolls = append(payrolls, p)
    }
    return payrolls, nil
}

func (r *payrollRepositoryImpl) CountByEmployeeID(ctx context.Context, employeeID string) (int64, error) {
    query := `SELECT COUNT(*) FROM payrolls WHERE employee_id = $1`
    var count int64
    err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
    return count, err
}
```

### 4. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Tambah implementasi:**

```go
func (s *payrollServiceImpl) GetMyPayrolls(ctx context.Context, userID string, page, perPage int, path string) (*helper.PayrollPagination, error) {
    // Find employee by userID
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, fmt.Errorf("employee not found: %w", err)
    }

    skip := int64((page - 1) * perPage)
    limit := int64(perPage)

    payrolls, err := s.payrollRepo.FindByEmployeeID(ctx, employee.ID.String(), skip, limit)
    if err != nil {
        return nil, err
    }

    total, err := s.payrollRepo.CountByEmployeeID(ctx, employee.ID.String())
    if err != nil {
        return nil, err
    }

    // Convert to list response
    employeeName := employee.UserName
    if employeeName == "" {
        employeeName = employee.Position
    }

    data := make([]dto.PayrollListResponse, len(payrolls))
    for i, p := range payrolls {
        data[i] = *helper.PayrollToListResponse(p, employeeName)
    }

    return helper.BuildPayrollPagination(data, page, perPage, total, path), nil
}

func (s *payrollServiceImpl) GetMyPayrollByID(ctx context.Context, userID string, payrollID string) (*dto.PayrollResponse, error) {
    // Verify ownership
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, ErrPayrollNotFound
    }

    payrollUUID, err := uuid.Parse(payrollID)
    if err != nil {
        return nil, ErrPayrollNotFound
    }

    payrollWithItems, err := s.payrollRepo.FindByIDWithItems(ctx, payrollUUID)
    if err != nil {
        return nil, ErrPayrollNotFound
    }

    // SECURITY: Pastikan payroll milik employee ini
    if payrollWithItems.Payroll.EmployeeID != employee.ID.String() {
        return nil, ErrPayrollNotFound // Jangan expose info bahwa payroll ada tapi bukan miliknya
    }

    employeeName := employee.UserName
    if employeeName == "" {
        employeeName = employee.Position
    }

    return helper.PayrollToResponse(payrollWithItems, employeeName, 
        employee.BankName, employee.BankAccountNumber, employee.BankAccountHolder), nil
}
```

### 5. [MODIFY] `internal/payroll/handler/payroll_handler.go`

**Tambah handler:**
```go
func (h *PayrollHandler) GetMyPayrolls(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)

    page, _ := strconv.Atoi(c.Query("page", "1"))
    perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

    if page < 1 { page = 1 }
    if perPage < 1 || perPage > 100 { perPage = 15 }

    path := c.Protocol() + "://" + c.Hostname() + c.Path()

    result, err := h.service.GetMyPayrolls(c.Context(), userID, page, perPage, path)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch payrolls", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "My payrolls retrieved", result)
}

func (h *PayrollHandler) GetMyPayrollByID(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)
    payrollID := c.Params("id")

    result, err := h.service.GetMyPayrollByID(c.Context(), userID, payrollID)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusNotFound, "Payroll not found", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "Payroll detail retrieved", result)
}
```

### 6. [MODIFY] `internal/payroll/routes.go`

**Tambah self-service routes SEBELUM admin routes:**
```go
// Self-service routes — all authenticated users
api.Get("/my", payrollHandler.GetMyPayrolls)
api.Get("/my/:id", payrollHandler.GetMyPayrollByID)

// Admin only routes (existing)
admin := api.Group("", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser))
// ...
```

> **PENTING**: Route `/my` harus sebelum `/:id` agar tidak ditangkap sebagai UUID.

## Verifikasi

1. `go build ./...` — compile sukses
2. Generate payroll untuk sebulan
3. Login sebagai karyawan (role USER)
4. `GET /api/v1/payrolls/my` → return list payroll miliknya
5. `GET /api/v1/payrolls/my/:id` → return detail payroll dengan items
6. `GET /api/v1/payrolls/my/:id` dengan payroll orang lain → 404 (bukan 403, untuk security)
