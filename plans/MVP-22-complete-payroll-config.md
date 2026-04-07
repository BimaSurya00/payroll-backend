# MVP-22: Complete Payroll Config — Repository + Integration

## Prioritas: 🟡 IMPORTANT
## Estimasi: 3 jam
## Tipe: Complete Partial (MVP-17)
## Dependency: MVP-21 (main.go fix)

---

## Deskripsi Masalah

MVP-17 hanya buat migration + entity. Yang masih hilang:
1. Repository implementation (query ke `payroll_configs`)
2. Salary calculator refactor (hapus hardcoded constants)
3. Payroll service update (baca config dari DB)
4. CRUD endpoint untuk admin manage configs

## File yang Diubah

### 1. [NEW] `internal/payroll/repository/payroll_config_repository.go`

Interface:
```go
type PayrollConfigRepository interface {
    FindAll(ctx context.Context) ([]*entity.PayrollConfig, error)
    FindByCode(ctx context.Context, code string) (*entity.PayrollConfig, error)
    FindActiveByType(ctx context.Context, configType string) ([]*entity.PayrollConfig, error)
    Update(ctx context.Context, id string, updates map[string]interface{}) error
    Create(ctx context.Context, config *entity.PayrollConfig) error
}
```

### 2. [NEW] `internal/payroll/repository/payroll_config_repository_impl.go`

Standard PostgreSQL implementation using `pgxpool`.

### 3. [MODIFY] `internal/payroll/helper/salary_calculator.go`

- **Hapus** hardcoded constants (`TransportAllowance`, `MealAllowance`, `LateDeductionPerDay`)
- **Tambah** `CalculateSalaryFromConfig()` yang baca dari `[]*entity.PayrollConfig`
- **Keep** `CalculateSalary()` sebagai deprecated fallback

### 4. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

- Inject `payrollConfigRepo` ke struct dan constructor
- Di `GenerateBulk`: fetch configs lalu pakai `CalculateSalaryFromConfig()`

### 5. [NEW] CRUD handler + routes untuk `/api/v1/payroll/configs`

- `GET /configs` — list semua configs
- `POST /configs` — create new config (admin)
- `PATCH /configs/:id` — update config (admin)

### 6. [MODIFY] `internal/payroll/routes.go`

Register config CRUD routes.

## Verifikasi

1. `go build ./...`
2. Generate payroll → verify allowance/deduction match DB config
3. Admin update transport ke 600K → generate lagi → verify 600K
4. `GET /api/v1/payroll/configs` → list configs
