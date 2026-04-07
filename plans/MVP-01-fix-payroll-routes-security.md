# MVP-01: Fix Payroll Routes Security

## Prioritas: 🔴 CRITICAL — Security Vulnerability
## Estimasi: 10 menit
## Tipe: Bug Fix / Security

---

## Deskripsi Masalah

File `internal/payroll/routes.go` **tidak menggunakan `middleware.HasRole`** pada admin routes.
Semua authenticated user (termasuk role `USER` biasa) bisa:
- Generate payroll bulk (`POST /api/v1/payrolls/generate`)
- Get semua payroll (`GET /api/v1/payrolls/`)
- Export CSV (`GET /api/v1/payrolls/export/csv`)
- Update status payroll (`PATCH /api/v1/payrolls/:id/status`)

Hanya `GET /api/v1/payrolls/:id` yang boleh diakses semua user (lihat slip sendiri nanti).

## File yang Diubah

### [MODIFY] `internal/payroll/routes.go`

**Tambah import:**
```go
"hris/middleware"
"hris/shared/constants"
```

**Ubah admin route group (line 28-33) dari:**
```go
// Admin only routes
admin := api.Group("")
// Note: Add middleware.RoleAuth("ADMIN", "SUPER_USER") when available

admin.Post("/generate", payrollHandler.GenerateBulk)
admin.Get("/", payrollHandler.GetAllPayrolls)
admin.Get("/export/csv", payrollHandler.ExportCSV)
```

**Menjadi:**
```go
// Admin only routes - ADMIN and SUPER_USER only
admin := api.Group("", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser))

admin.Post("/generate", payrollHandler.GenerateBulk)
admin.Get("/", payrollHandler.GetAllPayrolls)
admin.Get("/export/csv", payrollHandler.ExportCSV)
```

**Ubah juga line 37 — `UpdateStatus` harus pakai role check:**
```go
admin.Patch("/:id/status", payrollHandler.UpdateStatus)
```

## Verifikasi

1. Jalankan `go build ./...` — pastikan compile sukses
2. Test dengan user role `USER`:
   - `POST /api/v1/payrolls/generate` → harus return `403 Forbidden`
   - `GET /api/v1/payrolls/` → harus return `403 Forbidden`
   - `GET /api/v1/payrolls/export/csv` → harus return `403 Forbidden`
   - `PATCH /api/v1/payrolls/:id/status` → harus return `403 Forbidden`
3. Test dengan user role `ADMIN`:
   - Semua endpoint di atas harus return `200 OK`
4. Test `GET /api/v1/payrolls/:id` tetap bisa diakses semua role
