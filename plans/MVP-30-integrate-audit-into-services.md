# MVP-30: Integrate Audit Trail Into Payroll & Leave Services

## Prioritas: 🟡 IMPORTANT
## Estimasi: 1.5 jam
## Tipe: Integration Fix (Completes MVP-24)

---

## Deskripsi Masalah

Audit module lengkap (`internal/audit/`), table `audit_logs` ada, API endpoints ada, tapi **tidak pernah dipanggil** — zero references ke `auditService` di payroll maupun leave service.

## File yang Diubah

### 1. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Inject auditService:**
```diff
 type payrollServiceImpl struct {
     payrollRepo       payrollrepository.PayrollRepository
     employeeRepo      employeerepository.EmployeeRepository
     attendanceRepo    attendancerepository.AttendanceRepository
     payrollConfigRepo payrollconfigrepository.PayrollConfigRepository
+    auditService      auditservice.AuditService
     pool              *pgxpool.Pool
 }
```

**Add logging setelah commit di GenerateBulk:**
```go
go s.auditService.Log(context.Background(), auditservice.AuditEntry{
    UserID:       "system",
    Action:       "GENERATE",
    ResourceType: "payroll",
    NewData: map[string]interface{}{
        "month": req.PeriodMonth, "year": req.PeriodYear,
        "totalGenerated": generatedCount,
    },
})
```

### 2. [MODIFY] `internal/leave/service/leave_service.go`

**Inject auditService, log pada approve dan reject.**

### 3. [MODIFY] Route files

Update constructors untuk pass auditService.

## Verifikasi

```bash
go build ./...
# Harus: ✅ SUCCESS

# Then:
# 1. Generate payroll → check SELECT * FROM audit_logs
# 2. Approve leave → check audit_logs
```
