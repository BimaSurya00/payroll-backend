# MVP-32: Fix All Compile Errors

## Prioritas: 🔴 CRITICAL BLOCKER
## Estimasi: 30 menit
## Tipe: Bug Fix

---

## Deskripsi

App **tidak bisa di-build** karena 6 compile error di 5 file. Semua harus diperbaiki agar app bisa jalan.

## File yang Diubah

### 1. [MODIFY] `internal/attendance/repository/attendance_repository.go`

**Line 43: Hapus extra closing brace**
```diff
 type AttendanceMonthlySummary struct {
     EmployeeID   uuid.UUID
     TotalPresent int
     TotalLate    int
     TotalAbsent  int
     TotalLeave   int
     TotalDays    int
 }
-}
```

### 2. [MODIFY] `internal/payroll/entity/payroll_config.go`

**Line 18-22: Hapus duplicate `PayrollItem`** (keep yang di `payroll.go` yang punya ID/PayrollID/CreatedAt):
```diff
-type PayrollItem struct {
-    Name   string  `json:"name"`
-    Amount float64 `json:"amount"`
-    Type   string  `json:"type"` // EARNING or DEDUCTION
-}
```

> Note: `CalculateSalaryFromConfig()` di `salary_calculator.go` menggunakan `PayrollItem` tanpa ID — perlu di-update untuk menggunakan struct dari `payroll.go` atau buat type alias.

### 3. [MODIFY] `internal/attendance/dto/report_response.go`

**Line 3: Hapus unused import**
```diff
-import "time"
```

### 4. [MODIFY] `internal/audit/repository/audit_repository.go`

**Line 7: Hapus unused import**
```diff
 import (
     "context"
     "time"
-
-    "github.com/google/uuid"
     "example.com/hris/internal/audit/entity"
 )
```

### 5. [MODIFY] `internal/audit/repository/audit_repository_impl.go`

**Line 5-9: Hapus 4 unused imports**
```diff
 import (
     "context"
-    "encoding/json"
-    "time"
-
-    "github.com/google/uuid"
-    "github.com/jackc/pgx/v5"
     "github.com/jackc/pgx/v5/pgxpool"
     "example.com/hris/internal/audit/entity"
 )
```

### 6. [MODIFY] `internal/holiday/repository/holiday_repository.go`

**Line 7: Hapus unused import**
```diff
 import (
     "context"
     "time"
-
-    "github.com/google/uuid"
     "example.com/hris/internal/holiday/entity"
 )
```

### 7. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Tambah missing imports:**
```diff
 import (
     ...
+    payrollconfigrepository "example.com/hris/internal/payroll/repository"
+    "go.uber.org/zap"
 )
```

> Note: Jika `PayrollConfigRepository` ada di package terpisah, adjust import path-nya.

## Verifikasi

```bash
go build ./...
# Harus: ✅ SUCCESS — zero errors
```
