# 🔍 Verification Report — All 20 MVPs (Deep Source-Level Audit)

**Date**: 10 February 2026
**Scope**: MVP-01 through MVP-20

---

## 🚨 CRITICAL BLOCKER: `main.go` Regressed

> [!CAUTION]
> `main.go` hanya meregistrasi **3 dari 12 modul** yang seharusnya ada — semua modul PostgreSQL-based (employee, attendance, payroll, leave, schedule, overtime, dashboard, holiday, audit) **TIDAK TERDAFTAR**.

**Saat ini terdaftar:**
```go
auth.RegisterRoutes(app, mongoDB, keydb, cfg, jwtAuth)
user.RegisterRoutes(app, mongoDB, minioRepo, jwtAuth)
department.RegisterRoutes(app, postgres, jwtAuth)
```

**HILANG 9 module:**
| Module | Impact |
|--------|--------|
| `employee` | Semua employee CRUD tidak accessible |
| `attendance` | Clock in/out, report, correction tidak bisa diakses |
| `payroll` | Generate payroll, slip, export tidak bisa diakses |
| `leave` | Request, approve, balance tidak bisa diakses |
| `schedule` | Jadwal kerja tidak bisa diakses |
| `overtime` | Overtime request tidak bisa diakses |
| `dashboard` | Dashboard summary crash lagi (not registered) |
| `holiday` | Holiday CRUD tidak bisa diakses |
| `audit` | Audit logs tidak bisa diakses |

**Juga hilang dari `main.go`:**
- ❌ `sharedHelper.InitTimezone()` — timezone tetap UTC
- ❌ `middleware.GlobalRateLimiter()` — global rate limit tidak aktif

---

## 📊 Summary per MVP

### Round 1 (MVP-01 — MVP-10)

| MVP | Status | Issue |
|-----|--------|-------|
| 01 | ✅ OK | Security routes benar |
| 02 | ⚠️ | `InitTimezone()` hilang di main.go |
| 03 | ✅ OK | Attendance data digunakan di payroll |
| 04 | ✅ OK | `CountWorkingDays` exclude weekends |
| 05 | ✅ OK | Pagination count benar |
| 06 | ✅ OK | Scan bug fixed (MVP-11) |
| 07 | ✅ OK | Self-service routes ada |
| 08 | ✅ OK | Payroll slip routes ada |
| 09 | ✅ OK | Change password implemented |
| 10 | ✅ OK | DB transaction di payroll + leave |

### Round 2 (MVP-11 — MVP-20)

| MVP | Status | Issue |
|-----|--------|-------|
| 11 | ✅ OK | Dashboard Scan fixed, tapi module tidak registered di main.go |
| 12 | ✅ OK | N+1 payroll fixed — batch fetch `FindByIDs` benar |
| 13 | ✅ OK | N+1 leave fixed — batch fetch user/employee/leaveType |
| 14 | ⚠️ PARTIAL | `rate_limiter.go` ada, auth routes wired, **tapi `GlobalRateLimiter` hilang di main.go** |
| 15 | ✅ OK | Attendance report routes + service ada, **tapi module belum registered** |
| 16 | ✅ OK | Correction routes + service ada, **tapi module belum registered** |
| 17 | 🔴 PARTIAL | Hanya migration + entity. **Tidak ada**: repo impl, salary refactor, CRUD endpoint, payroll integration |
| 18 | 🔴 PARTIAL | Module holiday lengkap, workday.go enhanced. **Tapi leave service belum pakai `CountWorkingDaysExcluding`** |
| 19 | 🔴 PARTIAL | Module audit lengkap. **Tapi belum diintegrasikan ke payroll/leave service, belum registered** |
| 20 | 🔴 PARTIAL | Module department lengkap. **Tapi employee query belum di-JOIN, employee service belum accept departmentID** |

---

## 📋 Detail Temuan Kritis

### 1. `salary_calculator.go` Masih Hardcoded (MVP-17)
```go
// FILE: internal/payroll/helper/salary_calculator.go
const (
    TransportAllowance  = 500000  // ❌ HARDCODED
    MealAllowance       = 300000  // ❌ HARDCODED
    LateDeductionPerDay = 50000   // ❌ HARDCODED
)
```
Tabel `payroll_configs` sudah ada di migration tapi **tidak pernah di-query**. Payroll masih pakai `CalculateSalary()` hardcoded.

### 2. Leave Service Tanpa Holiday (MVP-18)
```go
// FILE: internal/leave/service/leave_service.go, line 121
totalDays := sharedHelper.CountWorkingDays(startDate, endDate)  // ❌ Tidak exclude holiday
```
`CountWorkingDaysExcluding()` sudah ada di workday.go tapi **leave service tidak punya `holidayRepo`** dan belum memanggil fungsi baru.

### 3. Audit Trail Tidak Terintegrasi (MVP-19)
Module audit sudah lengkap (entity, repo, service, handler, routes) tapi:
- ❌ Payroll service tidak inject `auditService`
- ❌ Leave service tidak inject `auditService`
- ❌ Tidak ada `go s.auditService.Log(...)` di operasi apapun
- Checklist di COMPLETED report sendiri pun mark integration sebagai ❌

### 4. Employee Tidak Pakai Department (MVP-20)
- ❌ Employee repository `Create()` tidak insert `department_id`
- ❌ `FindByID()` tidak JOIN departments
- ❌ `FindAll()` tidak include `department_name`
- ❌ Employee service tidak validate department exists

---

## 🛠️ New MVP Plans (Round 3)

Berdasarkan temuan di atas, berikut MVP baru yang harus dikerjakan:

### 🔴 CRITICAL — App tidak akan pernah jalan tanpa ini

**MVP-21: Fix main.go — Register All Modules**
- Daftar semua 12 module ke main.go
- Tambah `sharedHelper.InitTimezone()`
- Tambah `middleware.GlobalRateLimiter()`
- Import semua module packages
- **Estimasi: 1 jam**

### 🟡 IMPORTANT — Complete Partial Implementations

**MVP-22: Complete MVP-17 — Payroll Config Repository + Integration**
- Buat `payroll_config_repository_impl.go`
- Buat `CalculateSalaryFromConfig()` di salary_calculator.go
- Update payroll service GenerateBulk → baca config dari DB
- Buat CRUD endpoint (`GET/POST/PATCH /payroll/configs`)
- **Estimasi: 3 jam**

**MVP-23: Complete MVP-18 — Integrate Holiday into Leave**
- Inject `holidayRepo` ke leave service constructor
- Update `CreateLeaveRequest` → pakai `CountWorkingDaysExcluding`
- Update `leave/routes.go` → pass holidayRepo
- **Estimasi: 1 jam**

**MVP-24: Complete MVP-19 — Integrate Audit Trail**
- Inject `auditService` ke payroll service
- Inject `auditService` ke leave service
- Add `go s.auditService.Log(...)` di GenerateBulk, ApproveLeave, RejectLeave
- Register audit routes di main.go
- **Estimasi: 2 jam**

**MVP-25: Complete MVP-20 — Integrate Department into Employee**
- Update employee repository queries (JOIN departments)
- Update employee service → accept/validate departmentID
- Update employee response → include departmentName
- **Estimasi: 2 jam**

### 🟢 IMPROVEMENT — Performance & Reliability

**MVP-26: Add Missing Error Handling in Leave Service**
- `createLeaveAttendances` masih ignore errors (`_ = attendanceRepo.Create(...)`)
- Sebaiknya log error meskipun fire-and-forget
- **Estimasi: 30 menit**

**MVP-27: Add Employee Name to Payroll/Attendance Reports**
- MVP-12 payroll GetAll uses `emp.Position` as fallback (bukan nama)
- MVP-15 attendance report uses employee_id instead of name
- Perlu JOIN ke user collection atau denormalize name
- **Estimasi: 2 jam**

---

## 📊 Execution Order

```
1. MVP-21 (1 jam)  ← BLOCKER — app tidak jalan tanpa ini
2. MVP-22 (3 jam)  ← Payroll config integration
3. MVP-23 (1 jam)  ← Holiday + leave integration
4. MVP-24 (2 jam)  ← Audit trail integration
5. MVP-25 (2 jam)  ← Department + employee integration
6. MVP-26 (30 min) ← Error handling improvement
7. MVP-27 (2 jam)  ← Name resolution
```

**Total: ~11.5 jam**

> [!IMPORTANT]
> **MVP-21 HARUS dikerjakan pertama.** Tanpa fix main.go, SEMUA endpoint PostgreSQL-based (employee, attendance, payroll, leave, dashboard, dll) **TIDAK BISA DIAKSES** sama sekali. Aplikasi saat ini hanya bisa: register, login, dan manage department.
