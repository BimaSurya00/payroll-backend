# MVP-05: Fix Leave Pagination Count

## Prioritas: 🟡 IMPORTANT — UX Bug
## Estimasi: 30 menit
## Tipe: Bug Fix

---

## Deskripsi Masalah

Di `internal/leave/service/leave_service.go` line 212:
```go
// Simplified total count
total := int64(len(responses))
```

Ini mengembalikan jumlah item di **halaman saat ini**, bukan total keseluruhan.
Akibatnya: pagination metadata salah, frontend tidak tahu total halaman.

Contoh: 50 leave requests, page 1 per_page 10 → `total = 10` (salah, seharusnya 50).

## Solusi

Tambah method `Count` di leave request repository dan gunakan di service.

## File yang Diubah

### 1. [MODIFY] `internal/leave/repository/leave_request_repository.go`

**Tambah method di interface:**
```go
CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error)
```

### 2. [MODIFY] Implementasi repository (cari file implementation-nya di `internal/leave/repository/`)

**Implementasi:**
```go
func (r *leaveRequestRepositoryImpl) CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error) {
    query := `SELECT COUNT(*) FROM leave_requests WHERE employee_id = $1`
    var count int64
    err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
    return count, err
}
```

### 3. [MODIFY] `internal/leave/service/leave_service.go`

**Ubah `GetMyLeaveRequests` (line 181-214):**
```go
func (s *leaveService) GetMyLeaveRequests(ctx context.Context, userID string, page, perPage int) ([]dto.LeaveRequestResponse, int64, error) {
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, 0, fmt.Errorf("invalid user ID: %w", err)
    }

    employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, 0, fmt.Errorf("employee not found: %w", err)
    }

    offset := (page - 1) * perPage
    requests, err := s.leaveRequestRepo.FindByEmployeeID(ctx, employee.ID, perPage, offset)
    if err != nil {
        return nil, 0, err
    }

    // PERBAIKAN: Gunakan count query yang benar
    total, err := s.leaveRequestRepo.CountByEmployeeID(ctx, employee.ID)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count leave requests: %w", err)
    }

    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, 0, fmt.Errorf("user not found: %w", err)
    }

    responses := make([]dto.LeaveRequestResponse, len(requests))
    for i, req := range requests {
        leaveType, _ := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(req.LeaveTypeID))
        responses[i] = *s.toLeaveRequestResponse(&req, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid)
    }

    return responses, total, nil
}
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Buat 15 leave requests untuk satu employee
3. Query `GET /api/v1/leave/requests/my?page=1&per_page=5`
4. Verifikasi response:
   - `data` = 5 items
   - `pagination.total` = 15
   - `pagination.lastPage` = 3
