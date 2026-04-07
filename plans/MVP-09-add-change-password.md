# MVP-09: Add Change Password API

## Prioritas: 🟡 IMPORTANT — Basic Auth Feature
## Estimasi: 2 jam
## Tipe: New Feature

---

## Deskripsi Masalah

Tidak ada fitur change password. Jika karyawan ingin ganti password, tidak ada caranya.
Ini fitur auth paling basic yang HARUS ada sebelum production.

## Solusi

Tambah endpoint `PUT /api/v1/auth/change-password` yang membutuhkan old password dan new password.

## File yang Diubah

### 1. [NEW] `internal/auth/dto/change_password.go`

```go
package dto

type ChangePasswordRequest struct {
    OldPassword string `json:"oldPassword" validate:"required"`
    NewPassword string `json:"newPassword" validate:"required,min=8,password_strength"`
}
```

### 2. [MODIFY] `internal/auth/service/auth_service.go`

**Tambah method di interface:**
```go
type AuthService interface {
    // ... existing methods
    ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error
}
```

**Tambah error variable:**
```go
var (
    // ... existing errors
    ErrInvalidOldPassword = errors.New("old password is incorrect")
    ErrSamePassword       = errors.New("new password must be different from old password")
)
```

**Tambah implementasi:**
```go
func (s *authService) ChangePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
    // 1. Get user
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }

    // 2. Verify old password
    if !sharedHelper.CheckPassword(user.Password, req.OldPassword) {
        return ErrInvalidOldPassword
    }

    // 3. Check new password is different
    if sharedHelper.CheckPassword(user.Password, req.NewPassword) {
        return ErrSamePassword
    }

    // 4. Hash new password
    hashedPassword, err := sharedHelper.HashPassword(req.NewPassword)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // 5. Update password di MongoDB
    if err := s.userRepo.Update(ctx, userID, bson.M{"password": hashedPassword}); err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }

    // 6. Logout semua device (invalidate semua refresh token) — security best practice
    if err := s.tokenRepo.DeleteAllUserRefreshTokens(ctx, userID); err != nil {
        // Log tapi jangan gagalkan — password sudah berubah
        fmt.Printf("Warning: failed to revoke tokens after password change: %v\n", err)
    }

    return nil
}
```

**Tambah import jika belum ada:**
```go
"go.mongodb.org/mongo-driver/bson"
```

### 3. [MODIFY] `internal/auth/handler/auth_handler.go`

**Tambah handler:**
```go
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)

    var req dto.ChangePasswordRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }

    if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
        return helper.ValidationErrorResponse(c, validationErrors)
    }

    err := h.service.ChangePassword(c.Context(), userID, &req)
    if err != nil {
        if errors.Is(err, service.ErrInvalidOldPassword) {
            return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
        }
        if errors.Is(err, service.ErrSamePassword) {
            return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "Password changed successfully. Please login again.", nil)
}
```

### 4. [MODIFY] `internal/auth/routes.go`

**Tambah route di protected routes:**
```go
// Protected routes
auth.Post("/logout", jwtAuth, authHandler.Logout)
auth.Post("/logout-all", jwtAuth, authHandler.LogoutAll)
auth.Put("/change-password", jwtAuth, authHandler.ChangePassword) // TAMBAH
```

## Catatan Implementasi

- Setelah change password, semua refresh token dihapus → user harus login ulang di semua device
- Validator `password_strength` sudah ada (digunakan di RegisterRequest), jadi new password otomatis divalidasi
- JANGAN kirim password lama atau baru di response

## Verifikasi

1. `go build ./...` — compile sukses
2. Test cases:
   - `PUT /api/v1/auth/change-password` dengan old password benar → 200 OK
   - Cek login dengan password baru → berhasil
   - Cek login dengan password lama → gagal
   - `PUT /api/v1/auth/change-password` dengan old password salah → 400 Bad Request
   - `PUT /api/v1/auth/change-password` dengan new password = old password → 400 Bad Request
   - `PUT /api/v1/auth/change-password` dengan new password < 8 char → 422 Validation Error
3. Pastikan refresh token lama sudah tidak valid setelah change password
