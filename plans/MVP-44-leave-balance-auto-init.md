# MVP-44: Leave Balance Auto-Initialization

**Estimasi**: 1 hari  
**Impact**: 🟡 SEDANG — Business Logic Completion

---

## 1. Problem

Leave balances harus dibuat manual per-karyawan per-tahun. Setiap awal tahun (1 Januari), admin harus membuat balance untuk setiap karyawan × setiap leave type. Untuk 100 karyawan × 5 leave types = 500 records manual.

### Current Leave Balance Entity:
```go
type LeaveBalance struct {
    ID          string    // UUID
    EmployeeID  string    // FK → employees
    LeaveTypeID string    // FK → leave_types
    Year        int       // 2026, 2027, etc.
    Balance     int       // total allowance (from leave_type.max_days_per_year)
    Used        int       // days already used
    Pending     int       // days in pending requests
}
```

### Current Leave Type Entity:
```go
type LeaveType struct {
    MaxDaysPerYear int  // e.g., 12 for annual leave
    IsActive       bool
}
```

## 2. Implementation

### Step 1: Create Admin Endpoint — `POST /api/v1/leave/balances/init`

This endpoint generates leave balances for ALL active employees × ALL active leave types for a given year.

**Flow:**
1. Get all active employees (with `deleted_at IS NULL`)
2. Get all active leave types (`is_active = true`)
3. For each employee × leave type combination:
   - Check if balance already exists for that year
   - If NOT exists → create with `balance = leave_type.max_days_per_year`, `used = 0`, `pending = 0`
   - If exists → skip (don't overwrite)

### Step 2: Add to Leave Service Interface

**File**: `internal/leave/service/leave_service_interface.go` (or wherever `LeaveService` interface lives)

```go
InitBalances(ctx context.Context, year int) (*dto.InitBalancesResponse, error)
```

### Step 3: Implement in Leave Service

```go
func (s *leaveService) InitBalances(ctx context.Context, year int) (*dto.InitBalancesResponse, error) {
    // 1. Get all active employees
    employees, err := s.employeeRepo.FindAllWithoutPagination(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get employees: %w", err)
    }

    // 2. Get all active leave types
    leaveTypes, err := s.leaveTypeRepo.FindAllActive(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get leave types: %w", err)
    }

    created := 0
    skipped := 0

    for _, emp := range employees {
        for _, lt := range leaveTypes {
            // 3. Check if balance exists
            existing, err := s.leaveBalanceRepo.FindByEmployeeTypeAndYear(ctx, emp.ID, lt.ID, year)
            if err == nil && existing != nil {
                skipped++
                continue  // already exists
            }

            // 4. Create new balance
            balance := &entity.LeaveBalance{
                ID:          uuid.New().String(),
                EmployeeID:  emp.ID.String(),
                LeaveTypeID: lt.ID,
                Year:        year,
                Balance:     lt.MaxDaysPerYear,
                Used:        0,
                Pending:     0,
            }
            if err := s.leaveBalanceRepo.Create(ctx, balance); err != nil {
                return nil, fmt.Errorf("failed to create balance for employee %s: %w", emp.ID, err)
            }
            created++
        }
    }

    return &dto.InitBalancesResponse{
        Year:    year,
        Created: created,
        Skipped: skipped,
        Total:   len(employees) * len(leaveTypes),
    }, nil
}
```

### Step 4: Create DTO

**File**: `internal/leave/dto/` (add or create)

```go
type InitBalancesRequest struct {
    Year int `json:"year" validate:"required,min=2020,max=2100"`
}

type InitBalancesResponse struct {
    Year    int `json:"year"`
    Created int `json:"created"`
    Skipped int `json:"skipped"`
    Total   int `json:"total"`
}
```

### Step 5: Add Handler and Route

**Handler:**
```go
func (h *LeaveRequestHandler) InitBalances(c *fiber.Ctx) error {
    var req dto.InitBalancesRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
    }
    // validate...
    result, err := h.service.InitBalances(c.Context(), req.Year)
    // respond...
}
```

**Route** (admin only):
```go
leave.Post("/balances/init", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveRequestHandler.InitBalances)
```

### Step 6: Add `FindAllActive` to Leave Type Repository (if not exists)

```go
func (r *leaveTypeRepository) FindAllActive(ctx context.Context) ([]entity.LeaveType, error) {
    query := `SELECT id, name, code, description, max_days_per_year, is_paid, requires_approval, is_active, created_at, updated_at
              FROM leave_types WHERE is_active = true ORDER BY name`
    // ...
}
```

## 3. Files Changed

| # | File | Change |
|---|------|--------|
| 1 | `internal/leave/service/leave_service.go` | Add `InitBalances` method to implementation + interface |
| 2 | `internal/leave/dto/` | Add `InitBalancesRequest` + `InitBalancesResponse` |
| 3 | `internal/leave/handler/` | Add `InitBalances` handler |
| 4 | `internal/leave/routes.go` | Add `POST /balances/init` route |
| 5 | `internal/leave/repository/leave_type_repository.go` | Add `FindAllActive` if missing |

## 4. Verification

```bash
go build ./...
# API test:
curl -X POST http://localhost:3000/api/v1/leave/balances/init \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"year": 2026}'
# Expected: {"year": 2026, "created": 500, "skipped": 0, "total": 500}
```
