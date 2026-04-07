# MVP-35: Integrate Audit Trail Into Services

## Prioritas: 🟡 IMPORTANT
## Estimasi: 1.5 jam
## Tipe: Integration Fix (re-plan dari MVP-30 yang belum dikerjakan)

---

## Deskripsi

Audit module ada tapi tidak diintegrasikan ke service apapun. Zero references ke `auditService` di payroll/leave.

## Detail Perubahan

Sama persis seperti [MVP-30](./MVP-30-integrate-audit-into-services.md).

- Inject `auditService` ke payroll service struct + constructor
- Inject `auditService` ke leave service struct + constructor
- Add `go s.auditService.Log(...)` di GenerateBulk, ApproveLeave, RejectLeave
- Update `payroll/routes.go` + `leave/routes.go` untuk pass auditService

## Verifikasi

```bash
go build ./...
# Generate payroll → SELECT * FROM audit_logs → entry exists
# Approve leave → SELECT * FROM audit_logs → entry exists
```
