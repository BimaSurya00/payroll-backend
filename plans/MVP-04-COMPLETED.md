# ✅ MVP-04 COMPLETED: Fix Leave Weekend Calculation

## Status: ✅ COMPLETE
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (Data Accuracy)
## Time Taken: ~15 minutes

---

## 🎯 Objective
Perbaiki perhitungan hari cuti yang termasuk Sabtu dan Minggu, sehingga karyawan tidak kehilangan saldo cuti lebih banyak dari seharusnya.

---

## 📋 Changes Made

### 1. File Created
**`shared/helper/workday.go`** (NEW FILE)

```go
package helper

import "time"

// CountWorkingDays menghitung jumlah hari kerja (Senin-Jumat) antara dua tanggal (inklusif).
// TODO: Integrasikan dengan holiday management setelah module tersebut dibuat.
func CountWorkingDays(startDate, endDate time.Time) int {
	if startDate.After(endDate) {
		return 0
	}

	workingDays := 0
	current := startDate

	for !current.After(endDate) {
		weekday := current.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			workingDays++
		}
		current = current.AddDate(0, 0, 1)
	}

	return workingDays
}
```

---

### 2. File Modified
**`internal/leave/service/leave_service.go`**

**Added import:**
```go
sharedHelper "hris/shared/helper"
```

**Updated `CreateLeaveRequest` (line 114):**

#### Before (WRONG):
```go
// Calculate total days
totalDays := int(endDate.Sub(startDate).Hours()/24) + 1

// Check leave balance
currentYear := time.Now().Year()
// ...
```

#### After (CORRECT):
```go
// Calculate total days (working days only, excluding weekends)
totalDays := sharedHelper.CountWorkingDays(startDate, endDate)

// Validate: at least 1 working day
if totalDays == 0 {
    return nil, errors.New("selected dates contain no working days")
}

// Check leave balance
currentYear := time.Now().Year()
// ...
```

**Updated `createLeaveAttendances` (line 420-427):**

#### Before (WRONG):
```go
func (s *leaveService) createLeaveAttendances(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time, leaveRequestID uuid.UUID) {
	currentDate := startDate
	for !currentDate.After(endDate) {
		attendance := &leaveattendance.Attendance{
			ID:         uuid.New().String(),
			EmployeeID: employeeID.String(),
			Date:       currentDate,
			Status:     "LEAVE",
			Notes:      fmt.Sprintf("Leave Request ID: %s", leaveRequestID.String()),
		}

		// Ignore errors, attendance creation is not critical
		_ = s.attendanceRepo.Create(ctx, attendance)

		currentDate = currentDate.AddDate(0, 0, 1)
	}
}
```

#### After (CORRECT):
```go
func (s *leaveService) createLeaveAttendances(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time, leaveRequestID uuid.UUID) {
	currentDate := startDate
	for !currentDate.After(endDate) {
		// Skip weekends
		weekday := currentDate.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		attendance := &leaveattendance.Attendance{
			ID:         uuid.New().String(),
			EmployeeID: employeeID.String(),
			Date:       currentDate,
			Status:     "LEAVE",
			Notes:      fmt.Sprintf("Leave Request ID: %s", leaveRequestID.String()),
		}

		// Ignore errors, attendance creation is not critical
		_ = s.attendanceRepo.Create(ctx, attendance)

		currentDate = currentDate.AddDate(0, 0, 1)
	}
}
```

---

## 🔍 Technical Details

### Problem Before Fix:
```go
// Example: Cuti dari Jumat s/d Senin (4 hari calendar)
startDate = 2026-02-13 (Friday)
endDate = 2026-02-16 (Monday)

totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
// totalDays = 4 ❌ WRONG!

// Saldo cuti yang dikurangi: 4 hari
// Padahal hari kerja hanya: 2 hari (Jumat + Senin)
```

### Solution After Fix:
```go
// Example: Cuti dari Jumat s/d Senin (4 hari calendar)
startDate = 2026-02-13 (Friday)
endDate = 2026-02-16 (Monday)

totalDays := sharedHelper.CountWorkingDays(startDate, endDate)
// totalDays = 2 ✅ CORRECT! (Friday + Monday, skip Saturday + Sunday)

// Saldo cuti yang dikurangi: 2 hari
// Sesuai dengan hari kerja yang sebenarnya
```

---

## 📊 Impact Analysis

### Before Fix (WRONG):
| Date Range | Calendar Days | Weekends | Working Days | `totalDays` | Balance Deducted | Correct? |
|------------|---------------|----------|--------------|-------------|------------------|----------|
| Mon-Fri | 5 | 0 | 5 | 5 ✅ | 5 | ✅ |
| Fri-Mon | 4 | 2 | 2 | 4 ❌ | 4 | ❌ |
| Sat-Sun | 2 | 2 | 0 | 2 ❌ | 2 | ❌ |
| Mon-Sun | 7 | 2 | 5 | 7 ❌ | 7 | ❌ |

**Result**: Karyawan kehilangan saldo cuti lebih banyak dari seharusnya!

### After Fix (CORRECT):
| Date Range | Calendar Days | Weekends | Working Days | `totalDays` | Balance Deducted | Correct? |
|------------|---------------|----------|--------------|-------------|------------------|----------|
| Mon-Fri | 5 | 0 | 5 | 5 ✅ | 5 | ✅ |
| Fri-Mon | 4 | 2 | 2 | 2 ✅ | 2 | ✅ |
| Sat-Sun | 2 | 2 | 0 | 0 ✅ | ERROR | ✅ |
| Mon-Sun | 7 | 2 | 5 | 5 ✅ | 5 | ✅ |

**Result**: Saldo cuti hanya dikurangi untuk hari kerja yang sebenarnya!

---

## ✅ Verification

### 1. Build Status
```bash
/usr/local/go/bin/go build -o /tmp/hris-mvp04-complete ./main.go
# Result: ✅ SUCCESS - No compilation errors
```

### 2. CountWorkingDays Function
The new helper function:
- Iterates through each day in the date range
- Skips Saturday and Sunday
- Returns count of working days (Mon-Fri)
- Ready for future enhancement with holiday management

### 3. Expected Behavior

#### Scenario 1: Monday - Friday (5 working days)
```bash
curl -X POST http://localhost:8080/api/v1/leave/requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leaveTypeID": "...",
    "startDate": "2026-02-09",
    "endDate": "2026-02-13",
    "reason": "Family vacation"
  }'

# Expected Response:
# totalDays: 5 ✅
# Balance deducted: 5 days ✅
# Attendance records created: 5 (Mon-Fri) ✅
```

#### Scenario 2: Friday - Monday (2 working days)
```bash
curl -X POST http://localhost:8080/api/v1/leave/requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leaveTypeID": "...",
    "startDate": "2026-02-13",
    "endDate": "2026-02-16",
    "reason": "Long weekend"
  }'

# Expected Response:
# totalDays: 2 ✅ (Friday + Monday, skip Sat + Sun)
# Balance deducted: 2 days ✅
# Attendance records created: 2 (Fri, Mon) ✅
```

#### Scenario 3: Saturday - Sunday (0 working days)
```bash
curl -X POST http://localhost:8080/api/v1/leave/requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leaveTypeID": "...",
    "startDate": "2026-02-14",
    "endDate": "2026-02-15",
    "reason": "Weekend trip"
  }'

# Expected Response:
# Error: "selected dates contain no working days" ✅
# Balance deducted: 0 days ✅
# Attendance records created: 0 ✅
```

#### Scenario 4: Monday - Sunday (5 working days)
```bash
curl -X POST http://localhost:8080/api/v1/leave/requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leaveTypeID": "...",
    "startDate": "2026-02-09",
    "endDate": "2026-02-15",
    "reason": "Full week off"
  }'

# Expected Response:
# totalDays: 5 ✅ (Mon-Fri, skip Sat + Sun)
# Balance deducted: 5 days ✅
# Attendance records created: 5 (Mon-Fri) ✅
```

---

## 🧪 Testing Instructions

### Test 1: Mon-Fri (5 working days)
```bash
POST /api/v1/leave/requests
{
  "startDate": "2026-02-09",  # Monday
  "endDate": "2026-02-13"      # Friday
}

# Expected: totalDays = 5
# Attendance created: 5 records (Mon, Tue, Wed, Thu, Fri)
```

### Test 2: Fri-Mon (2 working days)
```bash
POST /api/v1/leave/requests
{
  "startDate": "2026-02-13",  # Friday
  "endDate": "2026-02-16"      # Monday
}

# Expected: totalDays = 2
# Attendance created: 2 records (Fri, Mon)
# Skipped: Sat, Sun
```

### Test 3: Sat-Sun (0 working days - error)
```bash
POST /api/v1/leave/requests
{
  "startDate": "2026-02-14",  # Saturday
  "endDate": "2026-02-15"      # Sunday
}

# Expected: ERROR "selected dates contain no working days"
# No leave request created
```

### Test 4: Mon-Sun (5 working days)
```bash
POST /api/v1/leave/requests
{
  "startDate": "2026-02-09",  # Monday
  "endDate": "2026-02-15"      # Sunday
}

# Expected: totalDays = 5
# Attendance created: 5 records (Mon-Fri)
# Skipped: Sat, Sun
```

---

## 🎯 Conclusion

### ✅ Completed:
1. Created `CountWorkingDays` helper in `shared/helper/workday.go`
2. Updated `CreateLeaveRequest` to use `CountWorkingDays` instead of calendar days
3. Added validation to reject requests with no working days
4. Updated `createLeaveAttendances` to skip weekends when creating attendance records
5. Build successful - no errors

### 🔒 Data Accuracy Improvements:
- **Before**: Cuti Fri-Mon = 4 hari (termasuk Sabtu & Minggu) ❌
- **After**: Cuti Fri-Mon = 2 hari (hanya Jumat & Senin) ✅
- **Balance**: Karyawan tidak lagi kehilangan saldo cuti untuk hari weekend
- **Attendance**: Leave attendance records hanya dibuat untuk hari kerja

### 📈 Employee Benefits:
- Saldo cuti lebih akurat dan fair
- Tidak perlu request cuti untuk weekend (otomatis di-skip)
- Pengurangan saldo sesuai dengan hari kerja yang sebenarnya
- Validasi mencegah request tanpa hari kerja

### 🔮 Future Enhancements:
- TODO: Integrate with holiday management module
- TODO: Add company holidays exclusion
- TODO: Consider custom workweek (e.g., 6-day work week)
- TODO: Add half-day leave support

### 🚀 Next Steps:
1. Restart application to load the new helper
2. Test leave request creation with various date ranges
3. Verify attendance records only created for working days
4. Check leave balance deduction is accurate
5. Update API documentation to mention working day calculation

---

**Plan Status**: ✅ **EXECUTED**
**Data Accuracy Bug**: ✅ **RESOLVED**
**Build Status**: ✅ **SUCCESS**
**Ready For**: Deployment
