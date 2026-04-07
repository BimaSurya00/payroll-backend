# MVP-36: Integrate Audit Trail Into Services

**Estimasi**: 2–3 jam  
**Impact**: MEDIUM-HIGH — Compliance + Traceability  
**Prerequisite**: MVP-32 (compile errors fixed) ✅

---

## 1. Problem Statement

The audit module exists and is fully functional ([`internal/audit/`](file:///home/bima/Documents/example.com/hris/internal/audit)), but **zero services call `auditService.Log()`**. The audit_logs table will always be empty. This means:

- No traceability for payroll generation/approval
- No traceability for leave approval/rejection
- No traceability for user CRUD operations
- No compliance with Depnaker audit requirements

Additionally, the audit repository uses `SELECT *` in 2 methods which is fragile against schema changes.

---

## 2. Scope

### In Scope
1. Wire `auditService` into **payroll** service (GenerateBulk, UpdateStatus)
2. Wire `auditService` into **leave** service (ApproveLeaveRequest, RejectLeaveRequest)
3. Wire `auditService` into **user** service (CreateUser, UpdateUser, DeleteUser)
4. Fix `SELECT *` in audit repository → use explicit column list
5. Update routes to pass `auditService` to services that need it

### Out of Scope
- Audit retention policy (future MVP)
- Audit for employee CRUD (future MVP)
- Audit dashboard/analytics (future MVP)

---

## 3. Current State

### Audit Module (Ready — No Changes Needed)

| File | Status |
|------|--------|
| [`entity/audit_log.go`](file:///home/bima/Documents/example.com/hris/internal/audit/entity/audit_log.go) | ✅ Complete — 11 fields |
| [`repository/audit_repository.go`](file:///home/bima/Documents/example.com/hris/internal/audit/repository/audit_repository.go) | ✅ Complete — interface with `Create`, `FindAll`, `Count`, etc. |
| [`repository/audit_repository_impl.go`](file:///home/bima/Documents/example.com/hris/internal/audit/repository/audit_repository_impl.go) | ⚠️ `SELECT *` in `FindByResource` (line 43) and `FindByUser` (line 82) |
| [`service/audit_service.go`](file:///home/bima/Documents/example.com/hris/internal/audit/service/audit_service.go) | ✅ Complete — `Log()` method ready, `AuditEntry` struct defined |
| [`handler/audit_handler.go`](file:///home/bima/Documents/example.com/hris/internal/audit/handler/audit_handler.go) | ✅ Complete |
| [`routes.go`](file:///home/bima/Documents/example.com/hris/internal/audit/routes.go) | ✅ Complete |

### `AuditEntry` struct (already defined in audit service):
```go
type AuditEntry struct {
    UserID       string
    UserName     string
    Action       string                 // CREATE, UPDATE, DELETE, APPROVE, REJECT, GENERATE
    ResourceType string                 // payroll, leave_request, user, etc.
    ResourceID   string
    OldData      interface{}
    NewData      interface{}
    Metadata     map[string]interface{}
    IPAddress    string
}
```

---

## 4. Implementation Steps

### Step 1: Fix `SELECT *` in Audit Repository

**File**: [`internal/audit/repository/audit_repository_impl.go`](file:///home/bima/Documents/example.com/hris/internal/audit/repository/audit_repository_impl.go)

Replace `SELECT *` in two methods with explicit column list.

#### `FindByResource` (line 43):

```diff
- SELECT * FROM audit_logs
+ SELECT id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address, created_at FROM audit_logs
```

#### `FindByUser` (line 82):

```diff
- SELECT * FROM audit_logs
+ SELECT id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address, created_at FROM audit_logs
```

---

### Step 2: Add `auditService` to Payroll Service

#### 2a. Modify `payrollServiceImpl` struct

**File**: [`internal/payroll/service/payroll_service_impl.go`](file:///home/bima/Documents/example.com/hris/internal/payroll/service/payroll_service_impl.go)

**Current struct** (lines 30–36):
```go
type payrollServiceImpl struct {
    payrollRepo       payrollrepository.PayrollRepository
    employeeRepo      employeerepository.EmployeeRepository
    attendanceRepo    attendancerepository.AttendanceRepository
    payrollConfigRepo payrollrepository.PayrollConfigRepository
    pool              *pgxpool.Pool
}
```

**Add import and field:**
```go
import (
    // ... existing imports ...
    auditservice "example.com/hris/internal/audit/service"
)

type payrollServiceImpl struct {
    payrollRepo       payrollrepository.PayrollRepository
    employeeRepo      employeerepository.EmployeeRepository
    attendanceRepo    attendancerepository.AttendanceRepository
    payrollConfigRepo payrollrepository.PayrollConfigRepository
    auditService      auditservice.AuditService    // ← ADD
    pool              *pgxpool.Pool
}
```

#### 2b. Modify `NewPayrollService` constructor

**Current** (check lines ~38–55 for current signature):
```go
func NewPayrollService(
    payrollRepo payrollrepository.PayrollRepository,
    employeeRepo employeerepository.EmployeeRepository,
    attendanceRepo attendancerepository.AttendanceRepository,
    payrollConfigRepo payrollrepository.PayrollConfigRepository,
    pool *pgxpool.Pool,
) PayrollService {
```

**Change to:**
```go
func NewPayrollService(
    payrollRepo payrollrepository.PayrollRepository,
    employeeRepo employeerepository.EmployeeRepository,
    attendanceRepo attendancerepository.AttendanceRepository,
    payrollConfigRepo payrollrepository.PayrollConfigRepository,
    auditService auditservice.AuditService,    // ← ADD
    pool *pgxpool.Pool,
) PayrollService {
    return &payrollServiceImpl{
        payrollRepo:       payrollRepo,
        employeeRepo:      employeeRepo,
        attendanceRepo:    attendanceRepo,
        payrollConfigRepo: payrollConfigRepo,
        auditService:      auditService,       // ← ADD
        pool:              pool,
    }
}
```

#### 2c. Add audit logging to `GenerateBulk`

Find the end of the `GenerateBulk` method (after successful generation). **Add audit log before the return statement:**

```go
// Audit log — payroll batch generated
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    UserID:       req.GeneratedBy, // or however the caller ID is available
    Action:       "GENERATE",
    ResourceType: "payroll",
    ResourceID:   "", // batch — no single resource ID
    NewData: map[string]interface{}{
        "period_start": req.PeriodStart,
        "period_end":   req.PeriodEnd,
        "count":        len(response.Payrolls), // adjust based on actual response field name
    },
})
```

> **Note**: Check what fields `GeneratePayrollRequest` has to get `GeneratedBy`. If it doesn't have a user context, you may need to check how the handler passes user info. The key is to log it as best-effort (ignoring the error with `_ =`).

#### 2d. Add audit logging to `UpdateStatus`

Find the `UpdateStatus` method. **Add audit log after successful status update:**

```go
// Audit log — payroll status updated
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    Action:       "UPDATE",
    ResourceType: "payroll",
    ResourceID:   id,
    OldData:      map[string]interface{}{"status": oldStatus},
    NewData:      map[string]interface{}{"status": req.Status},
})
```

> **Note**: You'll need to capture the old status before the update. Check the method implementation —it likely already fetches the payroll before updating it.

---

### Step 3: Add `auditService` to Leave Service

#### 3a. Modify `leaveService` struct

**File**: [`internal/leave/service/leave_service.go`](file:///home/bima/Documents/example.com/hris/internal/leave/service/leave_service.go)

**Current struct** (lines ~42–51):
```go
type leaveService struct {
    leaveTypeRepo    leaverepo.LeaveTypeRepository
    leaveBalanceRepo leaverepo.LeaveBalanceRepository
    leaveRequestRepo leaverepo.LeaveRequestRepository
    employeeRepo     employeerepo.EmployeeRepository
    userRepo         userRepo.UserRepository
    attendanceRepo   attendanceRepo.AttendanceRepository
    holidayRepo      holidayrepo.HolidayRepository
    pool             *pgxpool.Pool
}
```

**Add import and field:**
```go
import (
    // ... existing imports ...
    auditservice "example.com/hris/internal/audit/service"
)

type leaveService struct {
    // ... existing fields ...
    auditService     auditservice.AuditService    // ← ADD
    pool             *pgxpool.Pool
}
```

#### 3b. Modify `NewLeaveService` constructor

Add `auditService auditservice.AuditService` parameter and assignment.

**Current signature:**
```go
func NewLeaveService(
    leaveTypeRepo leaverepo.LeaveTypeRepository,
    leaveBalanceRepo leaverepo.LeaveBalanceRepository,
    leaveRequestRepo leaverepo.LeaveRequestRepository,
    employeeRepo employeerepo.EmployeeRepository,
    userRepo userRepo.UserRepository,
    attendanceRepo attendanceRepo.AttendanceRepository,
    holidayRepo holidayrepo.HolidayRepository,
    pool *pgxpool.Pool,
) LeaveService {
```

**Change to (add `auditService` parameter after `holidayRepo`):**
```go
func NewLeaveService(
    leaveTypeRepo leaverepo.LeaveTypeRepository,
    leaveBalanceRepo leaverepo.LeaveBalanceRepository,
    leaveRequestRepo leaverepo.LeaveRequestRepository,
    employeeRepo employeerepo.EmployeeRepository,
    userRepo userRepo.UserRepository,
    attendanceRepo attendanceRepo.AttendanceRepository,
    holidayRepo holidayrepo.HolidayRepository,
    auditService auditservice.AuditService,           // ← ADD
    pool *pgxpool.Pool,
) LeaveService {
```

Don't forget to assign it: `auditService: auditService,` in the struct literal.

#### 3c. Add audit logging to `ApproveLeaveRequest`

**Location**: end of `ApproveLeaveRequest` (lines ~345–431), after successful approval.

```go
// Audit log — leave request approved
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    UserID:       approverID,
    Action:       "APPROVE",
    ResourceType: "leave_request",
    ResourceID:   id,
    NewData: map[string]interface{}{
        "status": "APPROVED",
        "notes":  req.Notes,
    },
})
```

#### 3d. Add audit logging to `RejectLeaveRequest`

**Location**: end of `RejectLeaveRequest` (lines ~433–468), after successful rejection.

```go
// Audit log — leave request rejected
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    UserID:       approverID,
    Action:       "REJECT",
    ResourceType: "leave_request",
    ResourceID:   id,
    NewData: map[string]interface{}{
        "status":  "REJECTED",
        "reason":  req.Reason,
    },
})
```

---

### Step 4: Add `auditService` to User Service

#### 4a. Modify `service` struct in user service impl

**File**: [`internal/user/service/user_service_impl.go`](file:///home/bima/Documents/example.com/hris/internal/user/service/user_service_impl.go)

Add `auditService auditservice.AuditService` field to the struct.

Add import:
```go
auditservice "example.com/hris/internal/audit/service"
```

#### 4b. Modify constructor(s)

The user service has constructors (`NewUserService` or `NewUserServiceWithEmployee`). Add `auditService` as a parameter and assign it.

#### 4c. Add audit logging to `CreateUser`

After successful user creation:
```go
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    Action:       "CREATE",
    ResourceType: "user",
    ResourceID:   user.ID,
    NewData: map[string]interface{}{
        "email": user.Email,
        "name":  user.Name,
        "role":  user.Role,
    },
})
```

#### 4d. Add audit logging to `DeleteUser`

After successful deletion:
```go
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    Action:       "DELETE",
    ResourceType: "user",
    ResourceID:   id,
})
```

#### 4e. Add audit logging to `UpdateUser`

After successful update:
```go
_ = s.auditService.Log(ctx, auditservice.AuditEntry{
    Action:       "UPDATE",
    ResourceType: "user",
    ResourceID:   id,
    OldData:      oldUser,  // the user before update (already fetched in the method)
    NewData:      updatedUser,
})
```

---

### Step 5: Update Routes to Pass `auditService`

#### 5a. Payroll routes

**File**: [`internal/payroll/routes.go`](file:///home/bima/Documents/example.com/hris/internal/payroll/routes.go)

```diff
+ import (
+     auditrepository "example.com/hris/internal/audit/repository"
+     auditservice "example.com/hris/internal/audit/service"
+ )

  func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
      payrollRepo := payrollrepository.NewPayrollRepository(postgresDB.Pool)
      payrollConfigRepo := payrollrepository.NewPayrollConfigRepository(postgresDB.Pool)
      employeeRepo := employeerepository.NewEmployeeRepository(postgresDB.Pool)
      attendanceRepo := attendancerepository.NewAttendanceRepository(postgresDB.Pool)
+     auditRepo := auditrepository.NewAuditRepository(postgresDB.Pool)
+     auditService := auditservice.NewAuditService(auditRepo)

-     payrollService := payrollService.NewPayrollService(payrollRepo, employeeRepo, attendanceRepo, payrollConfigRepo, postgresDB.Pool)
+     payrollService := payrollService.NewPayrollService(payrollRepo, employeeRepo, attendanceRepo, payrollConfigRepo, auditService, postgresDB.Pool)
```

#### 5b. Leave routes

**File**: [`internal/leave/routes.go`](file:///home/bima/Documents/example.com/hris/internal/leave/routes.go)

```diff
+ import (
+     auditrepository "example.com/hris/internal/audit/repository"
+     auditservice "example.com/hris/internal/audit/service"
+ )

  func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
      // ... existing repos ...
+     auditRepo := auditrepository.NewAuditRepository(postgresDB.Pool)
+     auditService := auditservice.NewAuditService(auditRepo)

-     leaveService := service.NewLeaveService(leaveTypeRepo, leaveBalanceRepo, leaveRequestRepo, employeeRepo, userRepo, attendanceRepo, holidayRepo, postgresDB.Pool)
+     leaveService := service.NewLeaveService(leaveTypeRepo, leaveBalanceRepo, leaveRequestRepo, employeeRepo, userRepo, attendanceRepo, holidayRepo, auditService, postgresDB.Pool)
```

#### 5c. User routes

**File**: [`internal/user/routes.go`](file:///home/bima/Documents/example.com/hris/internal/user/routes.go)

Same pattern: create `auditRepo`, `auditService`, pass to user service constructor.

---

## 5. Files Changed Summary

| # | File | Change |
|---|------|--------|
| 1 | `internal/audit/repository/audit_repository_impl.go` | Replace `SELECT *` → explicit column list in `FindByResource` and `FindByUser` |
| 2 | `internal/payroll/service/payroll_service_impl.go` | Add `auditService` field, update constructor, add `Log()` calls in `GenerateBulk` and `UpdateStatus` |
| 3 | `internal/leave/service/leave_service.go` | Add `auditService` field to struct, update constructor |
| 4 | `internal/user/service/user_service_impl.go` | Add `auditService` field, update constructors, add `Log()` calls in `CreateUser`, `UpdateUser`, `DeleteUser` |
| 5 | `internal/payroll/routes.go` | Create `auditRepo`+`auditService`, pass to `NewPayrollService` |
| 6 | `internal/leave/routes.go` | Create `auditRepo`+`auditService`, pass to `NewLeaveService` |
| 7 | `internal/user/routes.go` | Create `auditRepo`+`auditService`, pass to user service constructor |

---

## 6. Design Decisions

### Why `_ = auditService.Log(...)` (fire-and-forget)?

Audit logging should **NOT** block or fail the main operation. If audit logging fails:
- The main operation (payroll generation, leave approval) should still succeed
- The error should be logged (the audit service could add internal logging)
- This is standard practice for audit trails — they are secondary to the main business operation

### Why not pass `auditService` from `main.go`?

Each module's `RegisterRoutes` creates its own `auditService` instance since they share the same stateless service backed by the same `pgxpool.Pool`. This avoids changing the `RegisterRoutes` signature of every module and keeps modules self-contained.

---

## 7. Verification Plan

### Build Verification
```bash
cd /home/bima/Documents/example.com/hris
go build ./...
go vet ./...
```

### Check No `SELECT *` Remaining in Audit Repo
```bash
grep -n "SELECT \*" internal/audit/repository/audit_repository_impl.go
# Should return NO matches
```

### Check Audit Log Calls Exist
```bash
grep -rn "auditService.Log\|s.auditService.Log" internal/payroll/ internal/leave/ internal/user/
# Should return matches in payroll, leave, and user services
```
