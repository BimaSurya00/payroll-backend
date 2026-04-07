# MVP-24: Complete Audit Trail — Integrate into Services

## Prioritas: 🟡 IMPORTANT
## Estimasi: 2 jam
## Tipe: Complete Partial (MVP-19)
## Dependency: MVP-21 (main.go fix)

---

## Deskripsi Masalah

Audit module sudah lengkap (`internal/audit/`), tapi **belum diintegrasikan** ke service mana pun. Tidak ada aksi yang tercatat di audit_logs.

## File yang Diubah

### 1. [MODIFY] `internal/payroll/service/payroll_service_impl.go`

**Inject auditService:**
```go
type payrollServiceImpl struct {
    // ... existing ...
    auditService auditservice.AuditService  // ADD
}
```

**Log di GenerateBulk (setelah commit):**
```go
go s.auditService.Log(context.Background(), auditservice.AuditEntry{
    UserID:       "system", // atau dari context
    Action:       "GENERATE",
    ResourceType: "payroll",
    NewData: map[string]interface{}{
        "period":         periodStartStr,
        "totalGenerated": generatedCount,
    },
})
```

### 2. [MODIFY] `internal/leave/service/leave_service.go`

**Inject auditService dan log di approve/reject:**
```go
// After approve
go s.auditService.Log(context.Background(), auditservice.AuditEntry{
    UserID: approverID, Action: "APPROVE",
    ResourceType: "leave_request", ResourceID: id,
})

// After reject
go s.auditService.Log(context.Background(), auditservice.AuditEntry{
    UserID: approverID, Action: "REJECT",
    ResourceType: "leave_request", ResourceID: id,
})
```

### 3. [MODIFY] Route files

Update `payroll/routes.go` dan `leave/routes.go` untuk pass auditService ke constructor.

## Verifikasi

1. `go build ./...`
2. Generate payroll → cek `SELECT * FROM audit_logs WHERE action='GENERATE'`
3. Approve leave → cek `SELECT * FROM audit_logs WHERE action='APPROVE'`
4. `GET /api/v1/audit/logs` → return entries
