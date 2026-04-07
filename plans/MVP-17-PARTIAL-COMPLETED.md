# 🔄 MVP-17: Dynamic Allowance & Deduction Config - PARTIAL COMPLETED

## Status: 🔄 PARTIAL COMPLETED
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (Business Flexibility)
## Time Taken: ~15 minutes

---

## 📊 Summary

Database migration dan entity untuk payroll configs telah dibuat. Namun, integrasi penuh memerlukan:
1. Repository implementation
2. Salary calculator refactor
3. Payroll service update
4. CRUD endpoints
5. Service dan handler untuk config management

---

## ✅ **Completed:**

### 1. **Database Migration**
- ✅ Created `000007_add_payroll_config.up.sql`
- ✅ Created `000007_add_payroll_config.down.sql`
- ✅ Seeded default configs (Transport, Meal, Late, Absent)

### 2. **Payroll Config Entity**
- ✅ Created `internal/payroll/entity/payroll_config.go`
- ✅ Defined `PayrollConfig` struct with all fields
- ✅ Added `PayrollItem` struct for breakdown

### 3. **Default Configs Seeded:**
```sql
- Transport Allowance: 500,000 (FIXED)
- Meal Allowance: 300,000 (FIXED)
- Late Deduction: 50,000 (PER_DAY)
- Absent Deduction: 0 (PER_DAY, calculated from daily salary)
```

---

## ⚠️ **Pending:**

### 1. **Payroll Config Repository**
Perlu buat file baru: `internal/payroll/repository/payroll_config_repository.go`

**Interface:**
```go
type PayrollConfigRepository interface {
    FindAll(ctx context.Context) ([]*entity.PayrollConfig, error)
    FindByCode(ctx context.Context, code string) (*entity.PayrollConfig, error)
    FindActiveByType(ctx context.Context, configType string) ([]*entity.PayrollConfig, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
}
```

### 2. **Salary Calculator Refactor**
Perlu update `internal/payroll/helper/salary_calculator.go`:

**Hapus Hardcoded Constants:**
```go
// REMOVE these lines:
const (
    TransportAllowance = 500000
    MealAllowance      = 300000
    LateDedPerDay      = 50000
)
```

**Tambah Configurable Function:**
```go
func CalculateSalaryFromConfig(baseSalary float64, lateDays, absentDays int,
    configs []*entity.PayrollConfig) (totalAllowance, totalDeduction, netSalary float64, items []*entity.PayrollItem) {

    for _, config := range configs {
        if !config.IsActive { continue }

        switch config.Type {
        case "EARNING":
            totalAllowance += config.Amount
            items = append(items, &entity.PayrollItem{
                Name: config.Name, Amount: config.Amount, Type: "EARNING",
            })
        case "DEDUCTION":
            var deductionAmount float64
            switch config.CalculationType {
            case "PER_DAY":
                if config.Code == "LATE_DEDUCTION" {
                    deductionAmount = config.Amount * float64(lateDays)
                } else if config.Code == "ABSENT_DEDUCTION" {
                    dailySalary := baseSalary / 22
                    deductionAmount = dailySalary * float64(absentDays)
                }
            case "FIXED":
                deductionAmount = config.Amount
            }
            if deductionAmount > 0 {
                totalDeduction += deductionAmount
                items = append(items, &entity.PayrollItem{
                    Name: config.Name, Amount: deductionAmount, Type: "DEDUCTION",
                })
            }
        }
    }

    netSalary = baseSalary + totalAllowance - totalDeduction
    return
}
```

### 3. **Payroll Service Update**
Perlu update `internal/payroll/service/payroll_service_impl.go`:

**Add Dependency:**
```go
type payrollServiceImpl struct {
    // ... existing fields
    payrollConfigRepo payrollconfigrepository.PayrollConfigRepository  // ADD THIS
}

func NewPayrollService(
    // ... existing params
    payrollConfigRepo payrollconfigrepository.PayrollConfigRepository,  // ADD THIS
    pool *pgxpool.Pool,
) PayrollService {
    return &payrollServiceImpl{
        // ... existing fields
        payrollConfigRepo: payrollConfigRepo,  // ADD THIS
        pool: pool,
    }
}
```

**Update GenerateBulk Method:**
```go
// Before employee loop
configs, err := s.payrollConfigRepo.FindActiveByType(ctx, "")
if err != nil {
    return nil, fmt.Errorf("failed to get payroll configs: %w", err)
}

// Inside employee loop (replace hardcoded calculation)
totalAllowance, totalDeduction, netSalary, items := helper.CalculateSalaryFromConfig(
    emp.SalaryBase, lateDays, absentDays, configs)
```

### 4. **Config CRUD Module**
Perlu buat module baru: `internal/payroll_config/`

Structure:
```
internal/payroll_config/
├── dto/payroll_config_dto.go
├── handler/payroll_config_handler.go
├── service/payroll_config_service.go
└── routes.go
```

**DTOs:**
```go
type CreateConfigRequest struct {
    Name            string  `json:"name" validate:"required"`
    Code            string  `json:"code" validate:"required,uppercase"`
    Type            string  `json:"type" validate:"required,oneof=EARNING DEDUCTION"`
    Amount          float64 `json:"amount" validate:"required"`
    CalculationType string  `json:"calculationType" validate:"required,oneof=FIXED PER_DAY PERCENTAGE"`
    Description     string  `json:"description"`
}

type UpdateConfigRequest struct {
    Name            *string  `json:"name,omitempty"`
    Amount          *float64 `json:"amount,omitempty"`
    CalculationType *string `json:"calculationType,omitempty"`
    Description     *string  `json:"description,omitempty"`
    IsActive        *bool    `json:"isActive,omitempty"`
}
```

**Routes:**
```go
// Admin only
api.Get("/configs", payrollConfigHandler.GetAll)
api.Post("/configs", payrollConfigHandler.Create)
api.Patch("/configs/:id", payrollConfigHandler.Update)
```

### 5. **Register Routes**
Add to main.go:
```go
import "example.com/hris/internal/payroll_config"

payroll_config.RegisterRoutes(app, postgres, jwtAuth)
```

---

## 📁 **Files Created:**

1. **`database/migrations/000007_add_payroll_config.up.sql`**
2. **`database/migrations/000007_add_payroll_config.down.sql`**
3. **`internal/payroll/entity/payroll_config.go`**

---

## 🎯 **Payroll Config API:**

```http
GET    /api/v1/payroll/configs          # List all configs (Admin)
POST   /api/v1/payroll/configs          # Create new config (Admin)
PATCH  /api/v1/payroll/configs/:id      # Update config (Admin)
```

---

## 📊 **Config Types:**

### **EARNING Configs:**
- Transport Allowance (FIXED)
- Meal Allowance (FIXED)
- Health Allowance (FIXED or PERCENTAGE)

### **DEDUCTION Configs:**
- Late Deduction (PER_DAY)
- Absent Deduction (PER_DAY)
- BPJS Deduction (FIXED or PERCENTAGE)

---

## 🧪 **Testing Example:**

### Test 1: Update Transport Allowance
```bash
# Default: 500,000
curl -X PATCH "http://localhost:8080/api/v1/payroll/configs/transport-uuid" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 600000
  }'

# Generate payroll
# Verify: Transport allowance = 600,000
```

### Test 2: Add New Allowance
```bash
curl -X POST "http://localhost:8080/api/v1/payroll/configs" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Health Allowance",
    "code": "HEALTH_ALLOWANCE",
    "type": "EARNING",
    "amount": 200000,
    "calculationType": "FIXED",
    "description": "Health insurance allowance"
  }'

# Generate payroll
# Verify: Health allowance = 200,000 added to earnings
```

### Test 3: Change Late Deduction
```bash
# Default: 50,000 per day
curl -X PATCH "http://localhost:8080/api/v1/payroll/configs/late-uuid" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 75000
  }'

# Generate payroll with 2 late days
# Verify: Late deduction = 75,000 × 2 = 150,000
```

---

## 🎓 **Design Decisions:**

### 1. **Config Types**
- **EARNING**: Allowances and bonuses (add to salary)
- **DEDUCTION**: Penalties and cuts (subtract from salary)

### 2. **Calculation Types**
- **FIXED**: Static amount (e.g., Transport 500K)
- **PER_DAY**: Multiplied by count (e.g., Late 50K × days)
- **PERCENTAGE**: Of base salary (future feature)

### 3. **Active/Inactive**
- Configs can be disabled without deletion
- Useful for temporary policies
- Easy to re-enable later

### 4. **Backward Compatibility**
- Default configs match hardcoded values
- Existing payrolls continue to work
- No breaking changes to current system

---

## 📋 **Next Steps to Complete:**

### Step 1: Create Payroll Config Repository
```bash
# Create file
touch internal/payroll/repository/payroll_config_repository.go
touch internal/payroll/repository/payroll_config_repository_impl.go
```

### Step 2: Implement Repository Methods
```go
func (r *payrollConfigRepository) FindAll(ctx context.Context) ([]*entity.PayrollConfig, error) {
    query := `SELECT * FROM payroll_configs ORDER BY type, name`
    rows, _ := r.pool.Query(ctx, query)
    // ... scan and return
}

func (r *payrollConfigRepository) FindActiveByType(ctx context.Context, configType string) ([]*entity.PayrollConfig, error) {
    query := `SELECT * FROM payroll_configs WHERE is_active = true`
    if configType != "" {
        query += ` AND type = $1`
    }
    // ... scan and return
}
```

### Step 3: Refactor Salary Calculator
- Remove hardcoded constants
- Add `CalculateSalaryFromConfig()` function
- Support FIXED, PER_DAY, PERCENTAGE types

### Step 4: Update Payroll Service
- Add `payrollConfigRepo` dependency
- Fetch configs in `GenerateBulk`
- Pass configs to calculator

### Step 5: Create Config CRUD Module
- Service with Create, Update, GetAll methods
- Handler with validation
- Routes (Admin only)

### Step 6: Register Routes
- Import payroll_config module
- Register in main.go

### Step 7: Run Migration
```bash
psql -U postgres -d example.com/hris -f database/migrations/000007_add_payroll_config.up.sql
```

---

## 🎯 **Benefits Once Completed:**

1. **Business Flexibility**
   - Each company can set their own allowances
   - Easy to adjust for inflation
   - Company-specific policies

2. **No Code Changes**
   - Update amounts via API
   - No deployment needed
   - Immediate effect

3. **Transparency**
   - All earnings/deductions visible
   - Clear breakdown in payslip
   - Audit trail of changes

4. **Multi-Tenancy Ready**
   - Foundation for per-company configs
   - Easy to extend for SaaS

---

## 🚧 **Known Limitations:**

1. **Not Yet Integrated**
   - Payroll service still uses hardcoded values
   - Salary calculator not refactored
   - Config CRUD not implemented

2. **No Validation**
   - Duplicate code prevention needed
   - Amount range validation
   - Type validation

3. **No History**
   - No audit trail for config changes
   - No versioning
   - No approval workflow

---

## 🔮 **Future Enhancements:**

1. **Versioning**
   - Track config changes history
   - Effective date ranges
   - Rollback capability

2. **Validation Rules**
   - Min/max amount constraints
   - Business rules (e.g., late deduction ≤ daily salary)
   - Dependency checks

3. **Advanced Calculations**
   - PERCENTAGE type (e.g., 5% of base salary)
   - Tiered calculations (e.g., overtime rates)
   - Conditional logic (e.g., probation vs permanent)

4. **Per-Company Configs**
   - Multi-tenancy support
   - Company templates
   - Industry standards

---

**Plan Status**: 🔄 **PARTIAL COMPLETED**
**Database Schema**: ✅ **CREATED**
**Entity**: ✅ **DEFINED**
**Integration**: ⚠️ **PENDING**
**CRUD Endpoints**: ⚠️ **PENDING**
**Migration**: ✅ **READY**
**Next Steps**: Complete repository, refactor calculator, update service, create CRUD endpoints
