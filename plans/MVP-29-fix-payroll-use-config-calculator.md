# MVP-29: Fix Payroll GenerateBulk — Use Config-Based Calculator

## Prioritas: 🔴 CRITICAL — Feature Not Working
## Estimasi: 30 menit
## Tipe: Integration Fix (Completes MVP-22)

---

## Deskripsi Masalah

`payrollConfigRepo` sudah di-inject ke payroll service struct, `CalculateSalaryFromConfig()` sudah ada di helper, tapi `GenerateBulk()` **masih memanggil deprecated `CalculateSalary()`** yang pakai hardcoded values. Variable `absentDays` dihitung tapi **tidak dipakai** di kalkulasi.

## File yang Diubah

### [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Perubahan di `GenerateBulk()` — setelah line 55 (sebelum loop):**
```go
// Fetch payroll configs
configs, err := s.payrollConfigRepo.FindActive(ctx)
if err != nil {
    // Fallback ke config default jika gagal fetch
    zap.L().Warn("failed to fetch payroll configs, using defaults", zap.Error(err))
    configs = nil
}
```

**Perubahan di line 97-100 (di dalam loop, ganti panggilan kalkulasi):**
```diff
-   allowance, deduction, netSalary := helper.CalculateSalary(
-       emp.SalaryBase,
-       lateDays,
-   )
+   var allowance, deduction, netSalary float64
+   if len(configs) > 0 {
+       allowance, deduction, netSalary, _ = helper.CalculateSalaryFromConfig(
+           emp.SalaryBase, lateDays, absentDays, configs,
+       )
+   } else {
+       allowance, deduction, netSalary = helper.CalculateSalary(
+           emp.SalaryBase, lateDays,
+       )
+   }
```

## Verifikasi

```bash
go build ./internal/payroll/...
# Harus: ✅ SUCCESS
```

**Functional test:**
1. Insert payroll config untuk transport allowance = 600000 (bukan default 500000)
2. Generate payroll
3. Verify allowance di response = 600000 (bukan 500000)
