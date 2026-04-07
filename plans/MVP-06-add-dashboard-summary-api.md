# MVP-06: Add Dashboard Summary API

## Prioritas: 🟡 IMPORTANT — Core Feature
## Estimasi: 4 jam
## Tipe: New Feature

---

## Deskripsi Masalah

Admin tidak memiliki overview/ringkasan data. Saat ini harus buka satu per satu:
- Berapa karyawan yang hadir hari ini?
- Berapa leave request yang pending?
- Berapa total payroll bulan ini?

Ini fitur paling kritis untuk kepuasan admin UMKM.

## Solusi

Buat module `dashboard` baru dengan satu endpoint yang mengembalikan ringkasan semua data.

## File yang Dibuat/Diubah

### 1. [NEW] `internal/dashboard/dto/dashboard_response.go`

```go
package dto

type DashboardSummary struct {
    Attendance  AttendanceSummary `json:"attendance"`
    Leave       LeaveSummary     `json:"leave"`
    Payroll     PayrollSummary   `json:"payroll"`
    Employee    EmployeeSummary  `json:"employee"`
}

type AttendanceSummary struct {
    TodayPresent int `json:"todayPresent"`
    TodayLate    int `json:"todayLate"`
    TodayAbsent  int `json:"todayAbsent"`
    TodayLeave   int `json:"todayLeave"`
    TotalEmployees int `json:"totalEmployees"`
}

type LeaveSummary struct {
    PendingRequests  int `json:"pendingRequests"`
    ApprovedThisMonth int `json:"approvedThisMonth"`
    RejectedThisMonth int `json:"rejectedThisMonth"`
}

type PayrollSummary struct {
    DraftCount    int     `json:"draftCount"`
    ApprovedCount int     `json:"approvedCount"`
    PaidCount     int     `json:"paidCount"`
    TotalNetSalary float64 `json:"totalNetSalary"`
    CurrentPeriod  string  `json:"currentPeriod"`
}

type EmployeeSummary struct {
    TotalActive   int `json:"totalActive"`
    TotalInactive int `json:"totalInactive"`
    NewThisMonth  int `json:"newThisMonth"`
}
```

### 2. [NEW] `internal/dashboard/service/dashboard_service.go`

```go
package service

import (
    "context"
    "time"

    attendanceRepo "example.com/hris/internal/attendance/repository"
    employeeRepo "example.com/hris/internal/employee/repository"
    leaveRepo "example.com/hris/internal/leave/repository"
    payrollRepo "example.com/hris/internal/payroll/repository"
    "example.com/hris/internal/dashboard/dto"
    sharedHelper "example.com/hris/shared/helper"
)

type DashboardService interface {
    GetSummary(ctx context.Context) (*dto.DashboardSummary, error)
}

type dashboardService struct {
    attendanceRepo attendanceRepo.AttendanceRepository
    employeeRepo   employeeRepo.EmployeeRepository
    leaveRepo      leaveRepo.LeaveRequestRepository
    payrollRepo    payrollRepo.PayrollRepository
}

func NewDashboardService(
    attendanceRepo attendanceRepo.AttendanceRepository,
    employeeRepo employeeRepo.EmployeeRepository,
    leaveRepo leaveRepo.LeaveRequestRepository,
    payrollRepo payrollRepo.PayrollRepository,
) DashboardService {
    return &dashboardService{
        attendanceRepo: attendanceRepo,
        employeeRepo:   employeeRepo,
        leaveRepo:      leaveRepo,
        payrollRepo:    payrollRepo,
    }
}

func (s *dashboardService) GetSummary(ctx context.Context) (*dto.DashboardSummary, error) {
    today := sharedHelper.Today()

    // 1. Attendance summary hari ini
    attendanceSummary := dto.AttendanceSummary{}
    // Query: COUNT per status WHERE date = today
    // Implementation: gunakan AttendanceFilter yang sudah ada di repo

    // 2. Leave summary
    leaveSummary := dto.LeaveSummary{}
    // Query: COUNT pending leave requests
    // Query: COUNT approved/rejected THIS MONTH

    // 3. Payroll summary bulan ini
    payrollSummary := dto.PayrollSummary{
        CurrentPeriod: today.Format("January 2006"),
    }
    // Query: COUNT per status, SUM net_salary

    // 4. Employee summary
    employeeSummary := dto.EmployeeSummary{}
    // Query: COUNT employees, COUNT new this month

    return &dto.DashboardSummary{
        Attendance: attendanceSummary,
        Leave:      leaveSummary,
        Payroll:    payrollSummary,
        Employee:   employeeSummary,
    }, nil
}
```

> **Catatan untuk agent**: Implementasi detail query di service harus disesuaikan dengan method yang tersedia di masing-masing repository. Jika method count belum ada, tambahkan di repository yang bersangkutan. Gunakan pattern yang sama seperti `CountByEmployeeID` di attendance repo.

### 3. [NEW] `internal/dashboard/handler/dashboard_handler.go`

```go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "example.com/hris/internal/dashboard/service"
    "example.com/hris/shared/helper"
)

type DashboardHandler struct {
    service service.DashboardService
}

func NewDashboardHandler(service service.DashboardService) *DashboardHandler {
    return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
    summary, err := h.service.GetSummary(c.Context())
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get dashboard summary", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "Dashboard summary retrieved", summary)
}
```

### 4. [NEW] `internal/dashboard/routes.go`

```go
package dashboard

import (
    "github.com/gofiber/fiber/v2"
    "example.com/hris/database"
    attendanceRepo "example.com/hris/internal/attendance/repository"
    employeeRepo "example.com/hris/internal/employee/repository"
    leaveRepo "example.com/hris/internal/leave/repository"
    payrollRepo "example.com/hris/internal/payroll/repository"
    "example.com/hris/internal/dashboard/handler"
    "example.com/hris/internal/dashboard/service"
    "example.com/hris/middleware"
    "example.com/hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
    attRepo := attendanceRepo.NewAttendanceRepository(postgresDB.Pool)
    empRepo := employeeRepo.NewEmployeeRepository(postgresDB.Pool)
    lvRepo := leaveRepo.NewLeaveRequestRepository(postgresDB.Pool)
    prRepo := payrollRepo.NewPayrollRepository(postgresDB.Pool)

    dashboardService := service.NewDashboardService(attRepo, empRepo, lvRepo, prRepo)
    dashboardHandler := handler.NewDashboardHandler(dashboardService)

    // Dashboard routes - ADMIN and SUPER_USER only
    dash := app.Group("/api/v1/dashboard", jwtAuth, middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser))
    dash.Get("/summary", dashboardHandler.GetSummary)
}
```

### 5. [MODIFY] `main.go`

**Tambah import dan register routes:**
```go
import (
    // ... existing imports
    "example.com/hris/internal/dashboard"
)

// Di dalam main(), setelah register module routes yang sudah ada:
dashboard.RegisterRoutes(app, postgres, jwtAuth)
```

## Catatan Implementasi

- Method count yang belum ada di repository harus ditambahkan dulu
- Gunakan single query per summary section jika memungkinkan (bukan N+1)
- Response harus cepat — target < 200ms karena ini halaman pertama yang dilihat admin

## Verifikasi

1. `go build ./...` — compile sukses
2. `GET /api/v1/dashboard/summary` dengan ADMIN token → 200 OK
3. Response harus mengandung semua 4 section: attendance, leave, payroll, employee
4. `GET /api/v1/dashboard/summary` dengan USER token → 403 Forbidden
