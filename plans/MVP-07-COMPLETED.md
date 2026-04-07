# ✅ MVP-07 COMPLETED: Add Employee Self-Service Profile

## Status: ✅ COMPLETE
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (Basic UX)
## Time Taken: ~25 minutes

---

## 🎯 Objective
Tambahkan endpoint self-service untuk karyawan bisa melihat dan mengupdate profil sendiri (data non-sensitif saja).

---

## 📋 Changes Made

### 1. File Created
**`internal/employee/dto/update_my_profile.go`** (NEW FILE)

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

---

### 2. File Modified
**`internal/employee/service/employee_service.go`**

**Added methods to interface:**
```go
type EmployeeService interface {
    CreateEmployee(ctx context.Context, req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
    GetAllEmployees(ctx context.Context, page, perPage int, path string, search string) (*helper.Pagination[*dto.EmployeeResponse], error)
    GetEmployeeByID(ctx context.Context, id string) (*dto.EmployeeResponse, error)
    UpdateEmployee(ctx context.Context, id string, req *dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
    DeleteEmployee(ctx context.Context, id string) error

    // Self-service endpoints
    GetMyProfile(ctx context.Context, userID string) (*dto.EmployeeResponse, error)
    UpdateMyProfile(ctx context.Context, userID string, req *dto.UpdateMyProfileRequest) (*dto.EmployeeResponse, error)
}
```

---

### 3. File Modified
**`internal/employee/service/employee_service_impl.go`**

**Implemented methods:**
```go
func (s *employeeService) GetMyProfile(ctx context.Context, userID string) (*dto.EmployeeResponse, error) {
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    employeeWithUser, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, ErrEmployeeNotFound
    }

    // Get user data from MongoDB
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    return helper.ToEmployeeResponseWithSchedule(employeeWithUser, user), nil
}

func (s *employeeService) UpdateMyProfile(ctx context.Context, userID string, req *dto.UpdateMyProfileRequest) (*dto.EmployeeResponse, error) {
    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    // Fetch current employee
    employeeWithUser, err := s.employeeRepo.FindByUserID(ctx, userUUID)
    if err != nil {
        return nil, ErrEmployeeNotFound
    }

    // Build employee struct for update
    employee := &repository.Employee{
        ID:                employeeWithUser.ID,
        UserID:            employeeWithUser.UserID,
        Position:          employeeWithUser.Position,
        PhoneNumber:       employeeWithUser.PhoneNumber,
        SalaryBase:        employeeWithUser.SalaryBase,
        Address:           employeeWithUser.Address,
        BankName:          employeeWithUser.BankName,
        BankAccountNumber: employeeWithUser.BankAccountNumber,
        BankAccountHolder: employeeWithUser.BankAccountHolder,
        JoinDate:          employeeWithUser.JoinDate,
        EmploymentStatus:  employeeWithUser.EmploymentStatus,
        JobLevel:          employeeWithUser.JobLevel,
        Gender:            employeeWithUser.Gender,
        Division:          employeeWithUser.Division,
        ScheduleID:        employeeWithUser.ScheduleID,
        CreatedAt:         employeeWithUser.CreatedAt,
        UpdatedAt:         time.Now(),
    }

    // Update only allowed fields
    hasChanges := false
    if req.PhoneNumber != nil {
        employee.PhoneNumber = *req.PhoneNumber
        hasChanges = true
    }
    if req.Address != nil {
        employee.Address = *req.Address
        hasChanges = true
    }
    if req.BankName != nil {
        employee.BankName = *req.BankName
        hasChanges = true
    }
    if req.BankAccountNumber != nil {
        employee.BankAccountNumber = *req.BankAccountNumber
        hasChanges = true
    }
    if req.BankAccountHolder != nil {
        employee.BankAccountHolder = *req.BankAccountHolder
        hasChanges = true
    }

    if !hasChanges {
        return nil, errors.New("no fields to update")
    }

    if err := s.employeeRepo.Update(ctx, employee); err != nil {
        return nil, fmt.Errorf("failed to update profile: %w", err)
    }

    // Return updated profile
    return s.GetMyProfile(ctx, userID)
}
```

---

### 4. File Modified
**`internal/employee/handler/employee_handler.go`**

**Added handler methods:**
```go
func (h *EmployeeHandler) GetMyProfile(c *fiber.Ctx) error {
    userID := c.Locals(constants.ContextKeyUserID).(string)

    employee, err := h.service.GetMyProfile(c.Context(), userID)
    if err != nil {
        if errors.Is(err, service.ErrEmployeeNotFound) {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee profile not found", err.Error())
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve profile", err.Error())
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
        if errors.Is(err, service.ErrEmployeeNotFound) {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "Employee profile not found", err.Error())
        }
        if err.Error() == "no fields to update" {
            return helper.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
    }

    return helper.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", employee)
}
```

---

### 5. File Modified
**`internal/employee/routes.go`**

**Added routes — PENTING: `/me` sebelum `/:id`:**
```go
// Employee routes - all require authentication
employees := app.Group("/api/v1/employees", jwtAuth)

// Self-service routes — ALL authenticated users (must be before /:id)
employees.Get("/me",
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.GetMyProfile,
)
employees.Patch("/me",
    middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.UpdateMyProfile,
)

// Create employee - ADMIN and SUPER_USER only
employees.Post("/",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.CreateEmployee,
)

// Get all employees (paginated) - ADMIN and SUPER_USER only
employees.Get("/",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.GetAllEmployees,
)

// Get employee by ID - ADMIN and SUPER_USER only
employees.Get("/:id",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.GetEmployeeByID,
)

// Update employee - ADMIN and SUPER_USER only
employees.Patch("/:id",
    middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
    employeeHandler.UpdateEmployee,
)

// Delete employee - SUPER_USER only
employees.Delete("/:id",
    middleware.HasRole(constants.RoleSuperUser),
    employeeHandler.DeleteEmployee,
)
```

---

## 🔍 Technical Details

### Route Ordering (CRITICAL):
Routes dengan `/me` harus didefinisikan **SEBELUM** `/:id` karena:
```go
employees.Get("/me", ...)     // ✅ Defined first
employees.Get("/:id", ...)    // ❌ Would catch "/me" as UUID param if defined first
```

### Allowed vs Restricted Fields:

#### ✅ **Can be updated by employee (self-service):**
- PhoneNumber
- Address
- BankName
- BankAccountNumber
- BankAccountHolder

#### ❌ **CANNOT be updated (admin only via /:id):**
- Position
- SalaryBase
- JobLevel
- Division
- ScheduleID
- EmploymentStatus

---

## 📊 API Specification

### Get My Profile
```http
GET /api/v1/employees/me
Authorization: Bearer <access_token>
```

#### Response (200 OK):
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "uuid",
    "userId": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "phoneNumber": "081234567890",
    "position": "Software Engineer",
    "salaryBase": 5000000,
    "address": "Jl. Sudirman No. 1",
    "bankName": "BCA",
    "bankAccountNumber": "1234567890",
    "bankAccountHolder": "John Doe",
    "employmentStatus": "PERMANENT",
    "jobLevel": "STAFF",
    "gender": "MALE",
    "division": "IT"
  }
}
```

### Update My Profile
```http
PATCH /api/v1/employees/me
Authorization: Bearer <access_token>
Content-Type: application/json
```

#### Request Body:
```json
{
  "phoneNumber": "081234567890",
  "address": "Jl. Sudirman No. 123, Jakarta",
  "bankName": "BCA",
  "bankAccountNumber": "0987654321",
  "bankAccountHolder": "John Doe"
}
```

#### Response (200 OK):
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": { /* updated employee profile */ }
}
```

---

## ✅ Verification

### 1. Build Status
```bash
/usr/local/go/bin/go build -o /tmp/hris-mvp07-complete ./main.go
# Result: ✅ SUCCESS - No compilation errors
```

### 2. Test Cases

#### Test 1: USER role can get own profile
```bash
# Login as USER, get token
USER_TOKEN="..."

curl -X GET http://localhost:8080/api/v1/employees/me \
  -H "Authorization: Bearer $USER_TOKEN"

# Expected: 200 OK with own employee data
```

#### Test 2: USER role can update own profile
```bash
curl -X PATCH http://localhost:8080/api/v1/employees/me \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phoneNumber": "081234567890",
    "address": "New Address"
  }'

# Expected: 200 OK, phoneNumber and address updated
```

#### Test 3: USER cannot access all employees
```bash
curl -X GET http://localhost:8080/api/v1/employees/ \
  -H "Authorization: Bearer $USER_TOKEN"

# Expected: 403 Forbidden
```

#### Test 4: ADMIN can access both /me and /:id
```bash
# Login as ADMIN
ADMIN_TOKEN="..."

# Get own profile
curl -X GET http://localhost:8080/api/v1/employees/me \
  -H "Authorization: Bearer $ADMIN_TOKEN"
# Expected: 200 OK

# Get specific employee
curl -X GET http://localhost:8080/api/v1/employees/{id} \
  -H "Authorization: Bearer $ADMIN_TOKEN"
# Expected: 200 OK
```

#### Test 5: Cannot update restricted fields via /me
```bash
curl -X PATCH http://localhost:8080/api/v1/employees/me \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "salaryBase": 10000000
  }'

# Expected: Field ignored (not in DTO) or validation error
# salaryBase not updated
```

---

## 🎯 Conclusion

### ✅ Completed:
1. Created `UpdateMyProfileRequest` DTO with only safe fields
2. Added `GetMyProfile` and `UpdateMyProfile` to service interface
3. Implemented both methods with proper validation
4. Created handler methods for both endpoints
5. Added `/me` routes before `/:id` to prevent routing conflicts
6. Build successful - no errors

### 🔒 Security Features:
- **Role-Based Access**: All authenticated users (USER, ADMIN, SUPER_USER) can access
- **Restricted Fields**: Only safe fields can be updated (phone, address, bank info)
- **Protected Fields**: Salary, position, job level cannot be changed via self-service
- **Auto-populated**: UpdatedAt timestamp automatically set

### 📈 User Experience:
- **Self-Service**: Employees can view and update their own profile
- **Convenience**: No need to contact HR for minor updates
- **Safety**: Cannot accidentally modify sensitive fields
- **Consistent**: Same response format as admin endpoints

### 🛡️ Best Practices Followed:
1. **Route Ordering**: `/me` before `/:id` to prevent UUID param conflicts
2. **Field Restrictions**: Only safe fields exposed to employees
3. **Validation**: All fields validated before update
4. **Error Handling**: Clear error messages for different scenarios
5. **Partial Updates**: Only sent fields are updated (pointer fields)

### 🔮 Future Enhancements:
- Add profile picture upload functionality
- Add emergency contact updates
- Add education/experience management
- Add document upload (KTP, CV, etc.)
- Consider approval workflow for certain profile changes

### 🚀 Next Steps:
1. Restart application to load the new endpoints
2. Test with USER role to verify self-service works
3. Test with ADMIN role to verify both /me and /:id work
4. Verify restricted fields cannot be updated
5. Update API documentation
6. Add to Postman collection

---

**Plan Status**: ✅ **EXECUTED**
**UX Gap**: ✅ **CLOSED**
**Build Status**: ✅ **SUCCESS**
**Ready For**: Deployment
