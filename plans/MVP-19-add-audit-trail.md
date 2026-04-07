# MVP-19: Add Audit Trail Module

## Prioritas: 🟡 IMPORTANT — Compliance & Security
## Estimasi: 4 jam
## Tipe: New Feature (Cross-cutting)

---

## Deskripsi Masalah

Tidak ada logging siapa melakukan apa. Kritis untuk:
- Payroll: siapa yang generate, siapa yang approve
- Leave: siapa yang approve/reject
- Employee: siapa yang ubah data karyawan
- Compliance: audit trail diperlukan jika ada dispute gaji/cuti

## Solusi

Buat audit trail module yang merekam setiap aksi penting ke tabel `audit_logs`.
Gunakan middleware/interceptor pattern agar bisa di-reuse di semua module.

## File yang Diubah

### 1. [NEW] Database migration: `000009_add_audit_logs.up.sql`

```sql
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,         -- Siapa yang melakukan aksi
    user_name VARCHAR(255),                 -- Nama user (denormalized untuk cepat query)
    action VARCHAR(50) NOT NULL,            -- CREATE, UPDATE, DELETE, APPROVE, REJECT, GENERATE
    resource_type VARCHAR(100) NOT NULL,    -- payroll, leave_request, employee, attendance
    resource_id VARCHAR(255),               -- ID resource yang diubah
    old_data JSONB,                         -- Data sebelum perubahan (nullable)
    new_data JSONB,                         -- Data setelah perubahan (nullable)
    metadata JSONB,                         -- Info tambahan (IP, user agent, etc.)
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);
```

### 2. [NEW] `internal/audit/entity/audit_log.go`

```go
package entity

import "time"

type AuditLog struct {
    ID           string    `json:"id" db:"id"`
    UserID       string    `json:"userId" db:"user_id"`
    UserName     string    `json:"userName" db:"user_name"`
    Action       string    `json:"action" db:"action"`         // CREATE, UPDATE, DELETE, APPROVE, REJECT, GENERATE
    ResourceType string    `json:"resourceType" db:"resource_type"` // payroll, leave, employee, etc.
    ResourceID   string    `json:"resourceId" db:"resource_id"`
    OldData      *string   `json:"oldData,omitempty" db:"old_data"` // JSON string
    NewData      *string   `json:"newData,omitempty" db:"new_data"` // JSON string
    Metadata     *string   `json:"metadata,omitempty" db:"metadata"`
    IPAddress    string    `json:"ipAddress" db:"ip_address"`
    CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}
```

### 3. [NEW] `internal/audit/repository/audit_repository.go`

```go
type AuditRepository interface {
    Create(ctx context.Context, log *entity.AuditLog) error
    FindByResource(ctx context.Context, resourceType, resourceID string, limit, offset int) ([]*entity.AuditLog, error)
    FindByUser(ctx context.Context, userID string, limit, offset int) ([]*entity.AuditLog, error)
    FindAll(ctx context.Context, filter AuditFilter, limit, offset int) ([]*entity.AuditLog, error)
    Count(ctx context.Context, filter AuditFilter) (int64, error)
}

type AuditFilter struct {
    UserID       *string
    Action       *string
    ResourceType *string
    DateFrom     *time.Time
    DateTo       *time.Time
}
```

### 4. [NEW] `internal/audit/service/audit_service.go`

```go
type AuditService interface {
    // Log an action — digunakan oleh service lain
    Log(ctx context.Context, entry AuditEntry) error

    // Query audit logs — untuk admin
    GetByResource(ctx context.Context, resourceType, resourceID string, page, perPage int) ([]dto.AuditLogResponse, int64, error)
    GetAll(ctx context.Context, filter AuditFilter, page, perPage int) ([]dto.AuditLogResponse, int64, error)
}

type AuditEntry struct {
    UserID       string
    UserName     string
    Action       string
    ResourceType string
    ResourceID   string
    OldData      interface{} // akan di-marshal ke JSON
    NewData      interface{}
    IPAddress    string
}
```

### 5. [NEW] `internal/audit/handler/audit_handler.go` + `routes.go`

**Routes (Admin/SuperUser only):**
```go
audit.Get("/logs", auditHandler.GetAll)                          // ?action=APPROVE&resource_type=leave
audit.Get("/logs/:resourceType/:resourceId", auditHandler.GetByResource) // History per resource
```

### 6. [MODIFY] Module yang perlu di-integrate

**Integrasi di payroll service — setelah GenerateBulk berhasil:**
```go
s.auditService.Log(ctx, audit.AuditEntry{
    UserID: adminID, Action: "GENERATE", ResourceType: "payroll",
    NewData: map[string]interface{}{
        "totalGenerated": generatedCount, "period": periodStartStr,
    },
})
```

**Integrasi di leave service — setelah approve/reject:**
```go
s.auditService.Log(ctx, audit.AuditEntry{
    UserID: approverID, Action: "APPROVE", ResourceType: "leave_request",
    ResourceID: id,
})
```

> **Catatan**: Audit logging harus async / fire-and-forget agar tidak memperlambat operasi utama. Gunakan goroutine: `go s.auditService.Log(...)`.

## Verifikasi

1. `go build ./...` — compile sukses
2. Run migration
3. Generate payroll → cek audit_logs ada entry dengan action="GENERATE"
4. Approve leave → cek audit_logs ada entry dengan action="APPROVE"
5. `GET /api/v1/audit/logs?action=GENERATE` → return log entries
6. `GET /api/v1/audit/logs/leave_request/{id}` → return history per leave request
