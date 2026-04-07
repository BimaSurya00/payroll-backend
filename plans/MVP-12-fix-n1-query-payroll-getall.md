# MVP-12: Fix N+1 Query in Payroll GetAll

## Prioritas: 🟡 IMPORTANT — Performance
## Estimasi: 2 jam
## Tipe: Performance Optimization

---

## Deskripsi Masalah

Di `internal/payroll/service/payroll_service_impl.go` line 329-369, `GetAll`:
```go
for i, payroll := range payrolls {
    employeeUUID, _ := uuid.Parse(payroll.EmployeeID)
    employee, _ := s.employeeRepo.FindByID(ctx, employeeUUID) // ← N+1!
    // ...
}
```

100 payroll records = **101 SQL queries** (1 FindAll + 100 FindByID). Ini sangat lambat.

## Solusi

Gunakan **batch fetch** — kumpulkan semua EmployeeID unik, query sekali, lalu map ke payroll.

## File yang Diubah

### 1. [MODIFY] `internal/employee/repository/employee_repository.go`

**Tambah method di interface:**
```go
FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Employee, error)
```

### 2. [MODIFY] Implementasi employee repository

```go
func (r *employeeRepositoryImpl) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Employee, error) {
    if len(ids) == 0 {
        return nil, nil
    }

    query := `SELECT id, user_id, position, salary_base, phone_number, address,
              employment_status, join_date, schedule_id, bank_name, bank_account_number,
              bank_account_holder, division, job_level, user_name, created_at, updated_at
              FROM employees WHERE id = ANY($1)`

    rows, err := r.pool.Query(ctx, query, ids)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var employees []*entity.Employee
    for rows.Next() {
        emp := &entity.Employee{}
        err := rows.Scan(&emp.ID, &emp.UserID, &emp.Position, &emp.SalaryBase,
            &emp.PhoneNumber, &emp.Address, &emp.EmploymentStatus, &emp.JoinDate,
            &emp.ScheduleID, &emp.BankName, &emp.BankAccountNumber, &emp.BankAccountHolder,
            &emp.Division, &emp.JobLevel, &emp.UserName, &emp.CreatedAt, &emp.UpdatedAt)
        if err != nil {
            return nil, err
        }
        employees = append(employees, emp)
    }
    return employees, nil
}
```

### 3. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Replace `GetAll` method (line 329-369):**
```go
func (s *payrollServiceImpl) GetAll(ctx context.Context, page, perPage int, path string) (*helper.PayrollPagination, error) {
    skip := int64((page - 1) * perPage)
    limit := int64(perPage)

    payrolls, err := s.payrollRepo.FindAll(ctx, skip, limit)
    if err != nil {
        return nil, err
    }

    total, err := s.payrollRepo.Count(ctx)
    if err != nil {
        return nil, err
    }

    // Collect unique employee IDs
    employeeIDSet := make(map[uuid.UUID]bool)
    for _, p := range payrolls {
        empUUID, err := uuid.Parse(p.EmployeeID)
        if err == nil {
            employeeIDSet[empUUID] = true
        }
    }

    employeeIDs := make([]uuid.UUID, 0, len(employeeIDSet))
    for id := range employeeIDSet {
        employeeIDs = append(employeeIDs, id)
    }

    // Batch fetch employees (1 query instead of N)
    employees, err := s.employeeRepo.FindByIDs(ctx, employeeIDs)
    employeeMap := make(map[string]*employeerepository.Employee)
    if err == nil {
        for _, emp := range employees {
            employeeMap[emp.ID.String()] = emp
        }
    }

    // Build response using map
    data := make([]dto.PayrollListResponse, len(payrolls))
    for i, payroll := range payrolls {
        employeeName := payroll.EmployeeID // fallback
        if emp, ok := employeeMap[payroll.EmployeeID]; ok {
            employeeName = emp.UserName
            if employeeName == "" {
                employeeName = emp.Position
            }
        }
        data[i] = *helper.PayrollToListResponse(payroll, employeeName)
    }

    return helper.BuildPayrollPagination(data, page, perPage, total, path), nil
}
```

> **Catatan**: Perlu adjust type import — bisa pakai employee entity langsung atau import alias sesuai kode existing.

## Verifikasi

1. `go build ./...` — compile sukses
2. Buat 50+ payroll records
3. `GET /api/v1/payrolls/?page=1&per_page=50`
4. Verifikasi: response ada employee name yang benar
5. Benchmark: response time harus signifikan lebih cepat vs sebelumnya
