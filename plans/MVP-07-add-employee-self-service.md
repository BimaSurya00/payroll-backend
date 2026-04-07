# MVP-07: Add Employee Self-Service Profile

## Prioritas: 🟡 IMPORTANT — Basic UX
## Estimasi: 2 jam
## Tipe: New Feature

---

## Deskripsi Masalah

Karyawan dengan role `USER` **tidak bisa melihat data employee-nya sendiri**.
Semua endpoint di `/api/v1/employees` hanya bisa diakses `ADMIN` dan `SUPER_USER`.

Karyawan perlu bisa:
1. Lihat profil employee sendiri
2. Update data non-sensitif sendiri (phone, address, bank info)

## Solusi

Tambah endpoint `/api/v1/employees/me` dan `/api/v1/employees/me` (PATCH) untuk self-service.

## File yang Diubah

### 1. [MODIFY] `internal/employee/service/employee_service.go`

**Tambah method di interface:**
```go
type EmployeeService interface {
    // ... existing methods

    // Self-service endpoints
    GetMyProfile(ctx context.Context, userID string) (*dto.EmployeeResponse, error)
    UpdateMyProfile(ctx context.Context, userID string, req *dto.UpdateMyProfileRequest) (*dto.EmployeeResponse, error)
}
```

### 2. [NEW] `internal/employee/dto/update_my_profile.go`

```go
package dto

// UpdateMyProfileRequest — hanya field yang boleh diubah sendiri oleh karyawan.
// Tidak termasuk: position, salaryBase, jobLevel, division, scheduleId, employmentStatus.
type UpdateMyProfileRequest struct {
    PhoneNumber       *string `json:"phoneNumber" validate:"omitempty,min=10,max=15"`
    Address           *string `json:"address" validate:"omitempty,min=5"`
    BankName          *string `json:"bankName" validate:"omitempty"`
    BankAccountNumber *string `json:"bankAccountNumber" validate:"omitempty"`
    BankAccountHolder *string `json:"bankAccountHolder" validate:"omitempty"`
}
```

### 3. [MODIFY] `internal/employee/service/employee_service_impl.go`

**Tambah implementasi:**

```go
func (s *employeeServiceImpl) GetMyProfile(ctx context.Context, userID string) (*dto.EmployeeResponse, error) {
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, ErrEmployeeNotFound
    }

    // Get user data from MongoDB
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    return helper.ToEmployeeResponse(employee, user), nil
}

func (s *employeeServiceImpl) UpdateMyProfile(ctx context.Context, userID string, req *dto.UpdateMyProfileRequest) (*dto.EmployeeResponse, error) {
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    employee, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, ErrEmployeeNotFound
    }

    // Build update map — hanya field yang dikirim
    updates := make(map[string]interface{})
    if req.PhoneNumber != nil {
        updates["phone_number"] = *req.PhoneNumber
    }
    if req.Address != nil {
        updates["address"] = *req.Address
    }
    if req.BankName != nil {
        updates["bank_name"] = *req.BankName
    }
    if req.BankAccountNumber != nil {
        updates["bank_account_number"] = *req.BankAccountNumber
    }
    if req.BankAccountHolder != nil {
        updates["bank_account_holder"] = *req.BankAccountHolder
    }

    if len(updates) == 0 {
        return nil, errors.New("no fields to update")
    }

    updates["updated_at"] = time.Now()

    if err := s.employeeRepo.Update(ctx, employee.ID, updates); err != nil {
        return nil, fmt.Errorf("failed to update profile: %w", err)
    }

    // Return updated profile
    return s.GetMyProfile(ctx, userID)
}
```

### 4. [MODIFY] `internal/employee/handler/employee_handler.go`

**Tambah handler methods:**
```go
func (h *EmployeeHandler) GetMyProfile(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)

    employee, err := h.service.GetMyProfile(c.Context(), userID)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee profile not found", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", employee)
}

func (h *EmployeeHandler) UpdateMyProfile(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)

    var req dto.UpdateMyProfileRequest
    if err := c.BodyParser(&req); err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
    }

    if validationErrors := customValidator.ValidateStruct(&req); len(validationErrors) > 0 {
        return helper.ValidationErrorResponse(c, validationErrors)
    }

    employee, err := h.service.UpdateMyProfile(c.Context(), userID, &req)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", employee)
}
```

### 5. [MODIFY] `internal/employee/routes.go`

**Tambah routes — PENTING: `/me` harus sebelum `/:id` agar tidak tertangkap sebagai UUID param:**

```go
employees := app.Group("/api/v1/employees", jwtAuth)

// Self-service routes — ALL authenticated users
employees.Get("/me",
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.GetMyProfile,
)
employees.Patch("/me",
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.UpdateMyProfile,
)

// Admin CRUD routes — existing (setelah /me)
employees.Post("/", ...)
employees.Get("/", ...)
employees.Get("/:id", ...)
employees.Patch("/:id", ...)
employees.Delete("/:id", ...)
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Login sebagai user role `USER`
3. `GET /api/v1/employees/me` → return profil employee sendiri (200 OK)
4. `PATCH /api/v1/employees/me` dengan `{"phoneNumber": "081234567890"}` → update berhasil
5. Pastikan field sensitif (salary, position) tetap **tidak bisa diubah** via `/me`
6. `GET /api/v1/employees/` dengan role USER → tetap 403 Forbidden
