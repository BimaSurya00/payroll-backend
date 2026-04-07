# MVP-16: Add Attendance Correction Flow

## Prioritas: 🟡 IMPORTANT — Operational Need
## Estimasi: 3 jam
## Tipe: New Feature

---

## Deskripsi Masalah

Jika karyawan lupa clock in/out, tidak ada cara koreksi selain langsung edit database.
Admin butuh fitur untuk:
1. Membuat attendance manual untuk karyawan
2. Mengoreksi clock in/out time yang salah
3. Menambahkan catatan koreksi

## Solusi

Tambah endpoint:
1. `POST /api/v1/attendances/correction` — Admin buat attendance manual
2. `PATCH /api/v1/attendances/:id/correction` — Admin koreksi attendance yang existing

Semua koreksi harus mencatat siapa yang melakukan koreksi (audit).

## File yang Diubah

### 1. [NEW] `internal/attendance/dto/correction_request.go`

```go
package dto

type CreateCorrectionRequest struct {
    EmployeeID string `json:"employeeId" validate:"required,uuid"`
    Date       string `json:"date" validate:"required"` // format: 2006-01-02
    ClockIn    string `json:"clockIn" validate:"required"` // format: 15:04
    ClockOut   string `json:"clockOut,omitempty"` // format: 15:04
    Status     string `json:"status" validate:"required,oneof=PRESENT LATE ABSENT"`
    Notes      string `json:"notes" validate:"required,min=5"` // Alasan koreksi wajib diisi
}

type UpdateCorrectionRequest struct {
    ClockIn  *string `json:"clockIn,omitempty"` // format: 15:04
    ClockOut *string `json:"clockOut,omitempty"` // format: 15:04
    Status   *string `json:"status,omitempty" validate:"omitempty,oneof=PRESENT LATE ABSENT"`
    Notes    *string `json:"notes,omitempty" validate:"omitempty,min=5"` // Alasan koreksi
}
```

### 2. [MODIFY] `internal/attendance/entity/attendance.go`

**Tambah field untuk tracking koreksi:**
```go
type Attendance struct {
    // ... existing fields ...
    CorrectedBy   *string    `json:"correctedBy,omitempty" db:"corrected_by"` // User ID yang mengoreksi
    CorrectedAt   *time.Time `json:"correctedAt,omitempty" db:"corrected_at"`
    CorrectionNote *string   `json:"correctionNote,omitempty" db:"correction_note"`
}
```

### 3. [NEW] Database migration: `000006_add_correction_fields.up.sql`

```sql
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS corrected_by VARCHAR(255);
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS corrected_at TIMESTAMP;
ALTER TABLE attendances ADD COLUMN IF NOT EXISTS correction_note TEXT;
```

### 4. [NEW] `000006_add_correction_fields.down.sql`

```sql
ALTER TABLE attendances DROP COLUMN IF EXISTS corrected_by;
ALTER TABLE attendances DROP COLUMN IF EXISTS corrected_at;
ALTER TABLE attendances DROP COLUMN IF EXISTS correction_note;
```

### 5. [MODIFY] `internal/attendance/service/` — Tambah correction methods

```go
// Tambah di interface
CreateCorrection(ctx context.Context, adminID string, req *dto.CreateCorrectionRequest) (*dto.AttendanceResponse, error)
UpdateCorrection(ctx context.Context, adminID, attendanceID string, req *dto.UpdateCorrectionRequest) (*dto.AttendanceResponse, error)
```

**Implementasi `CreateCorrection`:**
- Parse date → cek apakah sudah ada attendance di tanggal itu untuk employee
- Jika sudah ada → return error "attendance already exists, use update correction"
- Jika belum → create attendance baru dengan field `corrected_by`, `corrected_at`, `correction_note`
- Notes wajib diisi (untuk audit trail)

**Implementasi `UpdateCorrection`:**
- Cari attendance by ID
- Update field yang dikirim
- Set `corrected_by = adminID`, `corrected_at = now`, `correction_note = notes`

### 6. [MODIFY] `internal/attendance/routes.go`

```go
// Admin correction routes
api.Post("/correction",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    attendanceHandler.CreateCorrection)
api.Patch("/:id/correction",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    attendanceHandler.UpdateCorrection)
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Run migration
3. Admin create correction untuk employee yang lupa clock in → 201 Created
4. Admin update clock out yang salah → 200 OK, data terupdate dengan `corrected_by`
5. Verify `corrected_by`, `corrected_at`, `correction_note` terisi di database
6. User biasa tidak bisa akses endpoint correction → 403 Forbidden
