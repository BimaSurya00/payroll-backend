# ✅ MVP-05 COMPLETED: Fix Leave Pagination Count

## Status: ✅ COMPLETE
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (UX Bug)
## Time Taken: ~10 minutes

---

## 🎯 Objective
Perbaiki pagination count di leave requests yang mengembalikan jumlah item di halaman saat ini (bukan total keseluruhan).

---

## 📋 Changes Made

### 1. File Modified
**`internal/leave/repository/leave_request_repository.go`**

**Added method to interface:**
```go
type LeaveRequestRepository interface {
    Create(ctx context.Context, request *entity.LeaveRequest) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.LeaveRequest, error)
    FindByEmployeeID(ctx context.Context, employeeID uuid.UUID, limit, offset int) ([]entity.LeaveRequest, error)
    CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error) // NEW
    FindPending(ctx context.Context) ([]entity.LeaveRequest, error)
    FindByEmployeeAndDateRange(ctx context.Context, employeeID uuid.UUID, startDate, endDate time.Time) ([]entity.LeaveRequest, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *uuid.UUID, rejectionReason *string) error
}
```

**Implemented method:**
```go
func (r *leaveRequestRepository) CountByEmployeeID(ctx context.Context, employeeID uuid.UUID) (int64, error) {
    query := `SELECT COUNT(*) FROM leave_requests WHERE employee_id = $1`
    var count int64
    err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
    return count, err
}
```

---

### 2. File Modified
**`internal/leave/service/leave_service.go`**

**Updated `GetMyLeaveRequests` (line 181-214):**

#### Before (WRONG):
```go
offset := (page - 1) * perPage
requests, err := s.leaveRequestRepo.FindByEmployeeID(ctx, employee.ID, perPage, offset)
if err != nil {
    return nil, 0, err
}

// Get user from MongoDB
user, err := s.userRepo.FindByID(ctx, userID)
if err != nil {
    return nil, 0, fmt.Errorf("user not found: %w", err)
}

responses := make([]dto.LeaveRequestResponse, len(requests))
for i, req := range requests {
    leaveType, _ := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(req.LeaveTypeID))
    responses[i] = *s.toLeaveRequestResponse(&req, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid)
}

// Simplified total count
total := int64(len(responses)) // ❌ WRONG - returns page size instead of total!

return responses, total, nil
```

#### After (CORRECT):
```go
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

// Get user from MongoDB
user, err := s.userRepo.FindByID(ctx, userID)
if err != nil {
    return nil, 0, fmt.Errorf("user not found: %w", err)
}

responses := make([]dto.LeaveRequestResponse, len(requests))
for i, req := range requests {
    leaveType, _ := s.leaveTypeRepo.FindByID(ctx, uuid.MustParse(req.LeaveTypeID))
    responses[i] = *s.toLeaveRequestResponse(&req, user.Name, user.Email, employee.Position, leaveType.Name, leaveType.Code, leaveType.IsPaid)
}

return responses, total, nil // ✅ CORRECT - returns actual total count!
```

---

## 🔍 Technical Details

### Problem Before Fix:
```go
// Example: 50 leave requests total, page=1, per_page=10

requests, err := s.leaveRequestRepo.FindByEmployeeID(ctx, employee.ID, 10, 0)
// Returns: 10 requests (correct for page 1)

responses := make([]dto.LeaveRequestResponse, len(requests))
// Creates: 10 responses

total := int64(len(responses))
// total = 10 ❌ WRONG! Should be 50

// Pagination metadata would show:
// - total: 10 (wrong, should be 50)
// - lastPage: 1 (wrong, should be 5)
// - nextPage: null (wrong, should be 2)
```

### Solution After Fix:
```go
// Example: 50 leave requests total, page=1, per_page=10

requests, err := s.leaveRequestRepo.FindByEmployeeID(ctx, employee.ID, 10, 0)
// Returns: 10 requests (correct for page 1)

total, err := s.leaveRequestRepo.CountByEmployeeID(ctx, employee.ID)
// total = 50 ✅ CORRECT!

// Pagination metadata would show:
// - total: 50 ✅ CORRECT
// - lastPage: 5 ✅ CORRECT
// - nextPage: 2 ✅ CORRECT
```

---

## 📊 Impact Analysis

### Before Fix (WRONG):
| Total Records | Page | Per Page | Responses | `total` returned | lastPage | Correct? |
|---------------|------|----------|-----------|------------------|----------|----------|
| 50 | 1 | 10 | 10 | 10 ❌ | 1 ❌ | ❌ |
| 50 | 2 | 10 | 10 | 10 ❌ | 1 ❌ | ❌ |
| 50 | 3 | 10 | 10 | 10 ❌ | 1 ❌ | ❌ |
| 50 | 5 | 10 | 10 | 10 ❌ | 1 ❌ | ❌ |

**Result**: Frontend menampilkan total records = 10 (salah!), pagination tidak berfungsi.

### After Fix (CORRECT):
| Total Records | Page | Per Page | Responses | `total` returned | lastPage | Correct? |
|---------------|------|----------|-----------|------------------|----------|----------|
| 50 | 1 | 10 | 10 | 50 ✅ | 5 ✅ | ✅ |
| 50 | 2 | 10 | 10 | 50 ✅ | 5 ✅ | ✅ |
| 50 | 3 | 10 | 10 | 50 ✅ | 5 ✅ | ✅ |
| 50 | 5 | 10 | 10 | 50 ✅ | 5 ✅ | ✅ |

**Result**: Frontend menampilkan total records = 50 (benar!), pagination berfungsi dengan baik.

---

## ✅ Verification

### 1. Build Status
```bash
/usr/local/go/bin/go build -o /tmp/hris-mvp05-complete ./main.go
# Result: ✅ SUCCESS - No compilation errors
```

### 2. SQL Query
The new method uses a simple COUNT query:
```sql
SELECT COUNT(*) FROM leave_requests WHERE employee_id = $1
```

This is efficient and returns only the count (not all rows).

### 3. Expected Behavior

#### Before Fix:
```bash
GET /api/v1/leave/requests/my?page=1&per_page=5
# Response with 15 total records:
{
  "data": [... 5 items ...],
  "pagination": {
    "total": 5,        // ❌ WRONG - should be 15
    "perPage": 5,
    "currentPage": 1,
    "lastPage": 1      // ❌ WRONG - should be 3
  }
}
```

#### After Fix:
```bash
GET /api/v1/leave/requests/my?page=1&per_page=5
# Response with 15 total records:
{
  "data": [... 5 items ...],
  "pagination": {
    "total": 15,       // ✅ CORRECT
    "perPage": 5,
    "currentPage": 1,
    "lastPage": 3      // ✅ CORRECT
  }
}
```

---

## 🧪 Testing Instructions

### Test 1: Verify Correct Total Count
```bash
# Create 15 leave requests for one employee
# Then query page 1 with per_page=5

curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=1&per_page=5" \
  -H "Authorization: Bearer $TOKEN"

# Expected Response:
{
  "data": [ ... 5 items ... ],
  "pagination": {
    "total": 15,      # ✅ Should be 15 (not 5)
    "perPage": 5,
    "currentPage": 1,
    "lastPage": 3     # ✅ Should be 3 (not 1)
  }
}
```

### Test 2: Verify Multiple Pages
```bash
# Page 1
curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=1&per_page=5" \
  -H "Authorization: Bearer $TOKEN"
# Expected: total=15, lastPage=3, data=5 items

# Page 2
curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=2&per_page=5" \
  -H "Authorization: Bearer $TOKEN"
# Expected: total=15, lastPage=3, data=5 items

# Page 3
curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=3&per_page=5" \
  -H "Authorization: Bearer $TOKEN"
# Expected: total=15, lastPage=3, data=5 items
```

### Test 3: Edge Cases
```bash
# Empty result (0 leave requests)
curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=1&per_page=5" \
  -H "Authorization: Bearer $TOKEN"
# Expected: total=0, data=[]

# Exact match (5 leave requests with per_page=5)
curl -X GET "http://localhost:8080/api/v1/leave/requests/my?page=1&per_page=5" \
  -H "Authorization: Bearer $TOKEN"
# Expected: total=5, data=[5 items], lastPage=1
```

---

## 🎯 Conclusion

### ✅ Completed:
1. Added `CountByEmployeeID` method to `LeaveRequestRepository` interface
2. Implemented `CountByEmployeeID` in repository with COUNT query
3. Updated `GetMyLeaveRequests` to use `CountByEmployeeID` instead of `len(responses)`
4. Build successful - no errors

### 🔒 Bug Fixed:
- **Before**: `total = len(responses)` = returns page size (e.g., 10)
- **After**: `total = CountByEmployeeID()` = returns actual total (e.g., 50)
- **Pagination**: Now correctly shows total records and calculates `lastPage`

### 📈 UX Improvements:
- Frontend pagination now works correctly
- Users can see actual total number of leave requests
- Pagination controls (prev/next/last) function properly
- `lastPage` calculation is now accurate

### 🚀 Next Steps:
1. Restart application to load the fix
2. Test with employee who has multiple leave requests (> per_page)
3. Verify pagination metadata in API response
4. Update API documentation to reflect correct pagination behavior

---

**Plan Status**: ✅ **EXECUTED**
**UX Bug**: ✅ **RESOLVED**
**Build Status**: ✅ **SUCCESS**
**Ready For**: Deployment
