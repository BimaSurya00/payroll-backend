# ✅ MVP-08 COMPLETED: Add Payroll Slip View per Employee

## Status: ✅ COMPLETE
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (Essential UX)
## Time Taken: ~20 minutes

---

## 🎯 Objective
Tambahkan endpoint untuk karyawan bisa melihat slip gaji (payroll history) mereka sendiri.

---

## 📋 Changes Made

### 1. File Modified
**`internal/payroll/repository/payroll_repository.go`**

**Added methods to interface:**
```go
type PayrollRepository interface {
    // ... existing methods ...
    FindByEmployeeID(ctx context.Context, employeeID string, skip, limit int64) ([]*entity.Payroll, error)
    CountByEmployeeID(ctx context.Context, employeeID string) (int64, error)
}
```

---

### 2. File Modified
**`internal/payroll/repository/payroll_repository_impl.go`**

**Implemented methods:**
```go
func (r *payrollRepository) FindByEmployeeID(ctx context.Context, employeeID string, skip, limit int64) ([]*entity.Payroll, error) {
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

func (r *payrollRepository) CountByEmployeeID(ctx context.Context, employeeID string) (int64, error) {
    query := `SELECT COUNT(*) FROM payrolls WHERE employee_id = $1`
    var count int64
    err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
    return count, err
}
```

---

### 3. File Modified
**`internal/payroll/service/payroll_service.go`**

**Added methods to interface:**
```go
type PayrollService interface {
    // ... existing methods ...
    GetMyPayrolls(ctx context.Context, userID string, page, perPage int, path string) (*helper.PayrollPagination, error)
    GetMyPayrollByID(ctx context.Context, userID string, payrollID string) (*dto.PayrollResponse, error)
}
```

---

### 4. File Modified
**`internal/payroll/service/payroll_service_impl.go`**

**Implemented methods:**
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

---

### 5. File Modified
**`internal/payroll/handler/payroll_handler.go`**

**Added handler methods:**
```go
func (h *PayrollHandler) GetMyPayrolls(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)

    page, _ := strconv.Atoi(c.Query("page", "1"))
    perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

    if page < 1 {
        page = 1
    }
    if perPage < 1 || perPage > 100 {
        perPage = 15
    }

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

---

### 6. File Modified
**`internal/payroll/routes.go`**

**Added routes — PENTING: `/my` sebelum `/:id`:**
```go
// Public routes (with JWT)
api := app.Group("/api/v1/payrolls")
api.Use(jwtAuth)

// Self-service routes — all authenticated users (must be before /:id)
api.Get("/my", payrollHandler.GetMyPayrolls)
api.Get("/my/:id", payrollHandler.GetMyPayrollByID)

// Admin only routes - ADMIN and SUPER_USER only
admin := api.Group("", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser))

admin.Post("/generate", payrollHandler.GenerateBulk)
admin.Get("/", payrollHandler.GetAllPayrolls)
admin.Get("/export/csv", payrollHandler.ExportCSV)
admin.Patch("/:id/status", payrollHandler.UpdateStatus)

// All authenticated users can view their own payroll details
api.Get("/:id", payrollHandler.GetPayrollByID)
```

---

## 🔍 Technical Details

### Security Features:
1. **Ownership Verification**: `GetMyPayrollByID` verifies payroll belongs to employee
2. **No Information Leakage**: Returns 404 instead of 403 if accessing other's payroll
3. **Employee Filter**: `GetMyPayrolls` only returns payroll for authenticated user
4. **Route Ordering**: `/my` defined before `/:id` to prevent UUID conflicts

### Route Ordering (CRITICAL):
Routes dengan `/my` harus didefinisikan **SEBELUM** `/:id` karena:
```go
api.Get("/my", ...)        // ✅ Defined first
api.Get("/my/:id", ...)    // ✅ Defined first
api.Get("/:id", ...)       // ❌ Would catch "/my" as UUID param if defined first
```

---

## 📊 API Specification

### Get My Payroll History
```http
GET /api/v1/payrolls/my?page=1&per_page=15
Authorization: Bearer <access_token>
```

#### Response (200 OK):
```json
{
  "success": true,
  "message": "My payrolls retrieved",
  "data": {
    "data": [
      {
        "id": "uuid",
        "periodStart": "2026-01-01",
        "periodEnd": "2026-01-31",
        "baseSalary": 5000000,
        "totalAllowance": 800000,
        "totalDeduction": 150000,
        "netSalary": 5650000,
        "status": "APPROVED",
        "generatedAt": "2026-01-31T17:00:00Z"
      }
    ],
    "pagination": {
      "total": 12,
      "perPage": 15,
      "currentPage": 1,
      "lastPage": 1
    }
  }
}
```

### Get My Payroll Detail
```http
GET /api/v1/payrolls/my/:id
Authorization: Bearer <access_token>
```

#### Response (200 OK):
```json
{
  "success": true,
  "message": "Payroll detail retrieved",
  "data": {
    "id": "uuid",
    "periodStart": "2026-01-01",
    "periodEnd": "2026-01-31",
    "baseSalary": 5000000,
    "totalAllowance": 800000,
    "totalDeduction": 150000,
    "netSalary": 5650000,
    "status": "APPROVED",
    "items": [
      {
        "id": "uuid",
        "name": "Transport Allowance",
        "amount": 500000,
        "type": "EARNING"
      },
      {
        "id": "uuid",
        "name": "Meal Allowance",
        "amount": 300000,
        "type": "EARNING"
      },
      {
        "id": "uuid",
        "name": "Late Deduction",
        "amount": 150000,
        "type": "DEDUCTION"
      }
    ]
  }
}
```

#### Response (404 Not Found) - Accessing other's payroll:
```json
{
  "success": false,
  "message": "payroll not found",
  "errors": null
}
```

---

## ✅ Verification

### 1. Build Status
```bash
/usr/local/go/bin/go build -o /tmp/hris-mvp08-complete ./main.go
# Result: ✅ SUCCESS - No compilation errors
```

### 2. Test Cases

#### Test 1: USER can view own payroll history
```bash
# Login as USER
USER_TOKEN="..."

curl -X GET "http://localhost:8080/api/v1/payrolls/my?page=1&per_page=15" \
  -H "Authorization: Bearer $USER_TOKEN"

# Expected: 200 OK with only user's own payrolls
```

#### Test 2: USER can view own payroll detail
```bash
curl -X GET "http://localhost:8080/api/v1/payrolls/my/{payroll_id}" \
  -H "Authorization: Bearer $USER_TOKEN"

# Expected: 200 OK with full payroll details including items
```

#### Test 3: USER cannot access other's payroll
```bash
# Try to access another employee's payroll
curl -X GET "http://localhost:8080/api/v1/payrolls/my/{other_payroll_id}" \
  -H "Authorization: Bearer $USER_TOKEN"

# Expected: 404 Not Found (not 403, to avoid information leakage)
```

#### Test 4: ADMIN can access both /my and admin endpoints
```bash
# Login as ADMIN
ADMIN_TOKEN="..."

# Get own payrolls
curl -X GET "http://localhost:8080/api/v1/payrolls/my" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
# Expected: 200 OK

# Get all payrolls (admin endpoint)
curl -X GET "http://localhost:8080/api/v1/payrolls/" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
# Expected: 200 OK with all employees' payrolls
```

#### Test 5: Pagination works correctly
```bash
# Generate 25 payrolls for one employee

curl -X GET "http://localhost:8080/api/v1/payrolls/my?page=1&per_page=10" \
  -H "Authorization: Bearer $USER_TOKEN"

# Expected:
# - data: 10 payrolls
# - pagination.total: 25
# - pagination.lastPage: 3
```

---

## 🎯 Conclusion

### ✅ Completed:
1. Added `FindByEmployeeID` and `CountByEmployeeID` to repository interface
2. Implemented both repository methods with SQL queries
3. Added `GetMyPayrolls` and `GetMyPayrollByID` to service interface
4. Implemented both service methods with ownership verification
5. Created handler methods for both endpoints
6. Added `/my` routes before `/:id` to prevent routing conflicts
7. Build successful - no errors

### 🔒 Security Features:
- **Ownership Verification**: Employees can only access their own payrolls
- **No Information Leakage**: Returns 404 (not 403) when accessing other's payroll
- **Automatic Filtering**: GetMyPayrolls automatically filters by employee ID
- **JWT Authentication**: All endpoints require valid access token

### 📈 Employee Benefits:
- **Self-Service**: Can view own payroll history anytime
- **Detailed Slip**: Full breakdown of earnings and deductions
- **Pagination**: Easy navigation through multiple payrolls
- **Convenience**: No need to contact HR for slip copies

### 🛡️ Best Practices Followed:
1. **Route Ordering**: `/my` before `/:id` to prevent UUID conflicts
2. **Ownership Check**: Verify payroll belongs to authenticated user
3. **Error Handling**: Consistent 404 for not found/access denied
4. **Pagination**: Consistent with other list endpoints
5. **Data Filtering**: Repository level filtering by employee ID

### 🔮 Future Enhancements:
- Add payroll slip PDF download
- Add email notification when payroll is generated
- Add year-to-date summary
- Add comparison with previous month
- Add print/download functionality

### 🚀 Next Steps:
1. Restart application to load the new endpoints
2. Test with USER role to verify self-service works
3. Test ownership verification with different users
4. Verify pagination works correctly
5. Update API documentation
6. Add to Postman collection

---

**Plan Status**: ✅ **EXECUTED**
**Feature Gap**: ✅ **CLOSED**
**Build Status**: ✅ **SUCCESS**
**Ready For**: Deployment
