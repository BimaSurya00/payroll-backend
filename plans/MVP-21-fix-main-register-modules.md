# MVP-21: Fix main.go — Register All Modules

## Prioritas: 🔴 CRITICAL BLOCKER
## Estimasi: 1 jam
## Tipe: Critical Fix
## Dependency: NONE — harus dikerjakan PERTAMA

---

## Deskripsi Masalah

`main.go` hanya meregistrasi 3 dari 12 modul. Semua endpoint PostgreSQL-based (employee, attendance, payroll, leave, schedule, overtime, dashboard, holiday, audit) **tidak bisa diakses**.

Juga hilang:
- `sharedHelper.InitTimezone()` → timezone tetap UTC
- `middleware.GlobalRateLimiter()` → no protection

## File yang Diubah

### [MODIFY] `main.go`

**1. Tambah import semua module:**
```go
import (
    // ... existing imports ...
    "example.com/hris/internal/attendance"
    "example.com/hris/internal/audit"
    "example.com/hris/internal/dashboard"
    "example.com/hris/internal/employee"
    "example.com/hris/internal/holiday"
    "example.com/hris/internal/leave"
    "example.com/hris/internal/overtime"
    "example.com/hris/internal/payroll"
    "example.com/hris/internal/schedule"
    sharedHelper "example.com/hris/shared/helper"
)
```

**2. Tambah InitTimezone sebelum Fiber init:**
```go
// Initialize timezone
sharedHelper.InitTimezone()
```

**3. Tambah GlobalRateLimiter ke middleware stack:**
```go
app.Use(recover.New())
app.Use(middleware.Logger())
app.Use(middleware.GlobalRateLimiter())  // ADD
app.Use(cors.New(...))
```

**4. Register semua module routes:**
```go
// Register module routes
auth.RegisterRoutes(app, mongoDB, keydb, cfg, jwtAuth)
user.RegisterRoutes(app, mongoDB, minioRepo, jwtAuth)
department.RegisterRoutes(app, postgres, jwtAuth)

// PostgreSQL-based modules — SEMUA INI HILANG
employee.RegisterRoutes(app, postgres, mongoDB, jwtAuth)
attendance.RegisterRoutes(app, postgres, mongoDB, jwtAuth)
schedule.RegisterRoutes(app, postgres, jwtAuth)
payroll.RegisterRoutes(app, postgres, jwtAuth)
leave.RegisterRoutes(app, postgres, jwtAuth)
overtime.RegisterRoutes(app, postgres, jwtAuth)
dashboard.RegisterRoutes(app, postgres, jwtAuth)
holiday.RegisterRoutes(app, postgres, jwtAuth)
audit.RegisterRoutes(app, postgres, jwtAuth)
```

> **PENTING**: Cek signature `RegisterRoutes()` di setiap module's `routes.go` untuk parameter yang benar. Beberapa module mungkin butuh `*database.MongoDB` tambahan (employee, attendance).

## Verifikasi

1. `go build ./...` — HARUS compile sukses
2. `make run` — server start tanpa panic
3. `curl http://localhost:8080/api/v1/attendance/history` → 401 (bukan 404)
4. `curl http://localhost:8080/api/v1/payrolls` → 401 (bukan 404)
5. `curl http://localhost:8080/api/v1/leave/types` → response (bukan 404)
6. `curl http://localhost:8080/api/v1/dashboard/summary` → 401 (bukan 404)
