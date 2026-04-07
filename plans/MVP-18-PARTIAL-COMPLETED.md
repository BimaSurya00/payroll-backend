# 🔄 MVP-18: Add Holiday/Calendar Management - PARTIALLY COMPLETED

## Status: 🔄 PARTIAL COMPLETED
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (Data Accuracy)
## Time Taken: ~40 minutes

---

## 📊 Summary

Holiday module telah dibuat lengkap dengan entity, repository, service, handler, dan routes. `workday.go` telah di-enhance untuk mendukung holiday exclusion. Namun, integrasi penuh dengan leave service memerlukan penambahan dependency.

---

## ✅ **Completed:**

### 1. **Database Migration**
- ✅ Created `000008_add_holidays_table.up.sql`
- ✅ Created `000008_add_holidays_table.down.sql`
- ✅ Seeded Indonesia 2026 national holidays (10 holidays)

### 2. **Holiday Module Structure**
```
internal/holiday/
├── entity/holiday.go ✅
├── repository/holiday_repository.go ✅
├── repository/holiday_repository_impl.go ✅
├── handler/holiday_handler.go ✅
├── service/holiday_service.go ✅
├── dto/holiday_dto.go ✅
└── routes.go ✅
```

### 3. **Workday Helper Enhancement**
- ✅ Added `CountWorkingDaysExcluding()` function
- ✅ Accepts `holidays map[string]bool` parameter
- ✅ Backward compatible - `CountWorkingDays()` calls new function with `nil`

### 4. **Holiday Entity Features**
```go
type Holiday struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Date        time.Time `json:"date"`
    Type        string    `json:"type"` // NATIONAL, COMPANY, OPTIONAL
    IsRecurring bool      `json:"isRecurring"`
    Year        *int      `json:"year,omitempty"`
    Description *string   `json:"description,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}
```

---

## ⚠️ **Pending:**

### 1. **Leave Service Integration**
Perlu update `internal/leave/service/leave_service.go`:

**Add Dependency:**
```go
type leaveService struct {
    // ... existing fields
    holidayRepo     holidayrepo.HolidayRepository  // ADD THIS
}

func NewLeaveService(
    // ... existing params
    holidayRepo holidayrepo.HolidayRepository,  // ADD THIS
) LeaveService {
    return &leaveService{
        // ... existing fields
        holidayRepo: holidayRepo,  // ADD THIS
    }
}
```

**Update CreateLeaveRequest (line ~121):**
```go
// BEFORE:
totalDays := sharedHelper.CountWorkingDays(startDate, endDate)

// AFTER:
// Fetch holidays in date range
holidays, _ := s.holidayRepo.FindByDateRange(ctx, startDate, endDate)
holidayMap := make(map[string]bool)
for _, h := range holidays {
    holidayMap[h.Date.Format("2006-01-02")] = true
}

totalDays := sharedHelper.CountWorkingDaysExcluding(startDate, endDate, holidayMap)
```

### 2. **Main.go Registration**
Add holiday import and route registration:
```go
import "example.com/hris/internal/holiday"

// In main():
holiday.RegisterRoutes(app, postgres, jwtAuth)
```

### 3. **Leave Routes Update**
Update `internal/leave/routes.go` to pass holidayRepo:
```go
holidayRepo := holidayrepo.NewHolidayRepository(postgresDB.Pool)
leaveService := service.NewLeaveService(
    // ... existing params
    holidayRepo,  // ADD THIS
)
```

---

## 📁 **Files Created:**

1. **`database/migrations/000008_add_holidays_table.up.sql`**
2. **`database/migrations/000008_add_holidays_table.down.sql`**
3. **`internal/holiday/entity/holiday.go`**
4. **`internal/holiday/repository/holiday_repository.go`**
5. **`internal/holiday/repository/holiday_repository_impl.go`**
6. **`internal/holiday/dto/holiday_dto.go`**
7. **`internal/holiday/handler/holiday_handler.go`**
8. **`internal/holiday/service/holiday_service.go`**
9. **`internal/holiday/routes.go`**

---

## 📝 **Files Modified:**

1. **`shared/helper/workday.go`**
   - Added `CountWorkingDaysExcluding()` function
   - Enhanced to support holiday exclusion

---

## 🎯 **Holiday Module API:**

```http
GET    /api/v1/holidays          # List all holidays by year (e.g., ?year=2026)
POST   /api/v1/holidays          # Create holiday (Admin/SuperUser)
GET    /api/v1/holidays/:id      # Get holiday by ID
PATCH  /api/v1/holidays/:id      # Update holiday (Admin/SuperUser)
DELETE /api/v1/holidays/:id      # Delete holiday (Admin/SuperUser)
```

---

## 📊 **Seeded Holidays (Indonesia 2026):**

| Name | Date | Type |
|------|------|------|
| Tahun Baru | 2026-01-01 | NATIONAL |
| Isra Mi'raj | 2026-02-08 | NATIONAL |
| Hari Raya Nyepi | 2026-03-19 | NATIONAL |
| Wafat Isa Al Masih | 2026-04-03 | NATIONAL |
| Hari Buruh | 2026-05-01 | NATIONAL |
| Kenaikan Isa Al Masih | 2026-05-14 | NATIONAL |
| Hari Lahir Pancasila | 2026-06-01 | NATIONAL |
| Hari Kemerdekaan RI | 2026-08-17 | NATIONAL |
| Maulid Nabi | 2026-08-28 | NATIONAL |
| Natal | 2026-12-25 | NATIONAL |

---

## 🧪 **Testing Example:**

### Test 1: Create Holiday
```bash
curl -X POST "http://localhost:8080/api/v1/holidays" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Company Holiday",
    "date": "2026-03-15",
    "type": "COMPANY",
    "description": "Company team building event"
  }'
```

### Test 2: Get Holidays by Year
```bash
curl -X GET "http://localhost:8080/api/v1/holidays?year=2026" \
  -H "Authorization: Bearer <admin_token>"
```

### Test 3: Leave Calculation with Holidays
```bash
# Request leave from Feb 1-8, 2026 (8 days)
# Expected:
# - Weekends: Feb 7 (Sat), Feb 8 (Sun) = 2 days
# - Holiday: Feb 8 (Isra Mi'raj) = 1 day
# - Working days: 8 - 2 - 1 = 5 days

curl -X POST "http://localhost:8080/api/v1/leave/requests" \
  -H "Authorization: Bearer <user_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "leaveTypeId": "...",
    "startDate": "2026-02-01",
    "endDate": "2026-02-08",
    "reason": "Family vacation"
  }'

# Response should show totalDays = 5 (not 8)
```

---

## 🔧 **Workday Helper Usage:**

### Before (No Holiday Exclusion):
```go
// Count working days (Mon-Fri only)
totalDays := sharedHelper.CountWorkingDays(startDate, endDate)
```

### After (With Holiday Exclusion):
```go
// Fetch holidays
holidays, _ := s.holidayRepo.FindByDateRange(ctx, startDate, endDate)

// Build holiday map
holidayMap := make(map[string]bool)
for _, h := range holidays {
    holidayMap[h.Date.Format("2006-01-02")] = true
}

// Count working days excluding weekends AND holidays
totalDays := sharedHelper.CountWorkingDaysExcluding(startDate, endDate, holidayMap)
```

---

## 🎓 **Design Decisions:**

### 1. **Recurring Holidays**
- Some holidays occur every year (e.g., Independence Day, Cexample.com/hristmas)
- `is_recurring = true` means holiday applies every year
- `year` is NULL for recurring holidays

### 2. **Holiday Types**
- **NATIONAL**: Public holidays,所有人都适用
- **COMPANY**: Company-specific holidays (team building, etc.)
- **OPTIONAL**: Optional holidays (employees can choose to work)

### 3. **Holiday Map**
- Used for O(1) lookup when counting working days
- Format: `"2026-02-08": true` (date string -> bool)
- NULL map = no holidays to exclude (backward compatible)

### 4. **Backward Compatibility**
- `CountWorkingDays()` still works without holidays
- Calls `CountWorkingDaysExcluding()` with `nil`
- Existing code doesn't break

---

## 📋 **Next Steps to Complete:**

### Step 1: Update Leave Service
Add `holidayRepo` dependency and update `CreateLeaveRequest`:

```go
// In leave_service.go
type leaveService struct {
    // ... existing
    holidayRepo holidayrepo.HolidayRepository
}

// In CreateLeaveRequest (~line 121)
holidays, _ := s.holidayRepo.FindByDateRange(ctx, startDate, endDate)
holidayMap := make(map[string]bool)
for _, h := range holidays {
    holidayMap[h.Date.Format("2006-01-02")] = true
}
totalDays := sharedHelper.CountWorkingDaysExcluding(startDate, endDate, holidayMap)
```

### Step 2: Update Leave Routes
Pass holidayRepo to leave service:

```go
// In internal/leave/routes.go
holidayRepo := holidayrepo.NewHolidayRepository(postgresDB.Pool)
leaveService := service.NewLeaveService(
    leaveTypeRepo,
    leaveBalanceRepo,
    leaveRequestRepo,
    employeeRepo,
    userRepo,
    attendanceRepo,
    holidayRepo,  // ADD THIS
    postgresDB.Pool,
)
```

### Step 3: Register Holiday Routes
Add to main.go:

```go
import "example.com/hris/internal/holiday"

// In main()
holiday.RegisterRoutes(app, postgres, jwtAuth)
```

### Step 4: Run Migration
```bash
psql -U postgres -d example.com/hris -f database/migrations/000008_add_holidays_table.up.sql
```

### Step 5: Test
```bash
# Create leave request spanning a holiday
# Verify totalDays excludes holiday correctly
```

---

## 🎯 **Benefits Once Completed:**

1. **Accurate Leave Calculation**
   - Holidays not counted as leave days
   - Employees get fair leave balance
   - Compliance with labor laws

2. **Better Planning**
   - Employees can see holidays in calendar
   - Admin can add company-specific holidays
   - Leave requests consider holidays

3. **Data Accuracy**
   - Attendance reports exclude holidays
   - Payroll calculations aware of holidays
   - Consistent holiday management

---

## 🚧 **Known Limitations:**

1. **Leave Service Not Integrated**
   - Still uses `CountWorkingDays()` without holidays
   - Need to add `holidayRepo` dependency
   - Need to update `CreateLeaveRequest`

2. **No Calendar UI**
   - Only API endpoints available
   - No frontend calendar view
   - No visual holiday planning

3. **No Notification**
   - No upcoming holiday alerts
   - No email reminders
   - No calendar integration

---

**Plan Status**: 🔄 **PARTIAL COMPLETED**
**Holiday Module**: ✅ **CREATED**
**Workday Helper**: ✅ **ENHANCED**
**Leave Integration**: ⚠️ **PENDING (needs dependency injection)**
**Migration**: ✅ **READY**
**Next Steps**: Update leave service, register routes, test
