# MVP-17: Dynamic Allowance & Deduction Config

## Prioritas: 🟡 IMPORTANT — Business Flexibility
## Estimasi: 4 jam
## Tipe: New Feature + Refactor
## Dependency: Sebelum multi-tenancy, karena nanti config ini per-company

---

## Deskripsi Masalah

Di `internal/payroll/helper/salary_calculator.go`:
```go
const (
    TransportAllowance = 500000
    MealAllowance      = 300000
    LateDedPerDay      = 50000
)
```

Hardcoded. Setiap UMKM punya struktur gaji berbeda. UMKM A mungkin kasih transport 200K, UMKM B kasih 1M.

## Solusi

Buat tabel `payroll_configs` di PostgreSQL dan ubah salary calculator untuk membaca dari config, bukan constants.

## File yang Diubah

### 1. [NEW] Database migration: `000007_add_payroll_config.up.sql`

```sql
CREATE TABLE IF NOT EXISTS payroll_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,           -- 'Transport Allowance', 'Meal Allowance', etc.
    code VARCHAR(50) UNIQUE NOT NULL,     -- 'TRANSPORT_ALLOWANCE', 'MEAL_ALLOWANCE', etc.
    type VARCHAR(20) NOT NULL,            -- 'EARNING' or 'DEDUCTION'
    amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    calculation_type VARCHAR(20) NOT NULL DEFAULT 'FIXED', -- 'FIXED', 'PER_DAY', 'PERCENTAGE'
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed default configs (migrasi dari hardcoded values)
INSERT INTO payroll_configs (name, code, type, amount, calculation_type, description) VALUES
    ('Transport Allowance', 'TRANSPORT_ALLOWANCE', 'EARNING', 500000, 'FIXED', 'Tunjangan transport bulanan'),
    ('Meal Allowance', 'MEAL_ALLOWANCE', 'EARNING', 300000, 'FIXED', 'Tunjangan makan bulanan'),
    ('Late Deduction', 'LATE_DEDUCTION', 'DEDUCTION', 50000, 'PER_DAY', 'Potongan keterlambatan per hari'),
    ('Absent Deduction', 'ABSENT_DEDUCTION', 'DEDUCTION', 0, 'PER_DAY', 'Potongan absen per hari (hitung dari gaji harian)');
```

### 2. [NEW] `000007_add_payroll_config.down.sql`

```sql
DROP TABLE IF EXISTS payroll_configs;
```

### 3. [NEW] `internal/payroll/entity/payroll_config.go`

```go
package entity

import "time"

type PayrollConfig struct {
    ID              string    `json:"id" db:"id"`
    Name            string    `json:"name" db:"name"`
    Code            string    `json:"code" db:"code"`
    Type            string    `json:"type" db:"type"`         // EARNING | DEDUCTION
    Amount          float64   `json:"amount" db:"amount"`
    CalculationType string    `json:"calculationType" db:"calculation_type"` // FIXED | PER_DAY | PERCENTAGE
    IsActive        bool      `json:"isActive" db:"is_active"`
    Description     string    `json:"description" db:"description"`
    CreatedAt       time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}
```

### 4. [NEW] `internal/payroll/repository/payroll_config_repository.go`

```go
type PayrollConfigRepository interface {
    FindAll(ctx context.Context) ([]*entity.PayrollConfig, error)
    FindByCode(ctx context.Context, code string) (*entity.PayrollConfig, error)
    FindActiveByType(ctx context.Context, configType string) ([]*entity.PayrollConfig, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
}
```

### 5. [MODIFY] `internal/payroll/helper/salary_calculator.go`

**Refactor — buat configurable version:**
```go
package helper

// CalculateSalaryFromConfig menghitung gaji berdasarkan config dari database
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

### 6. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Update `GenerateBulk` — gunakan config dari DB alih-alih hardcoded:**
```go
// Sebelum loop employees, fetch payroll configs
configs, err := s.payrollConfigRepo.FindActiveByType(ctx, "")
if err != nil { return nil, fmt.Errorf("failed to get payroll configs: %w", err) }

for _, emp := range employees {
    // ... attendance summary ...
    totalAllowance, totalDeduction, netSalary, items := helper.CalculateSalaryFromConfig(
        emp.SalaryBase, lateDays, absentDays, configs)
    // ... create payroll with these values ...
}
```

### 7. [NEW] CRUD endpoint untuk admin manage payroll configs

Tambah routes: `GET /api/v1/payroll-configs`, `PATCH /api/v1/payroll-configs/:id`
(Admin/SuperUser only)

## Verifikasi

1. `go build ./...` — compile sukses
2. Run migration — seed default configs
3. Generate payroll → verify allowance/deduction sama seperti sebelumnya (backward compatible)
4. Admin update transport allowance ke 600.000 via API
5. Generate payroll lagi → verify transport allowance sekarang 600.000
