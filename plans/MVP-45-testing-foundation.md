# MVP-45: Testing Foundation

**Estimasi**: 5 hari  
**Impact**: 🔴 TINGGI — Safety Net sebelum Production

---

## 1. Problem

Saat ini hampir tidak ada test. Setiap perubahan di satu modul bisa break modul lain tanpa terdeteksi. Untuk production launch, kita **wajib** punya test coverage pada critical paths.

## 2. Scope: Critical Paths Only (MVP)

Tidak perlu 100% coverage. Fokus pada 3 flow paling critical:

| # | Area | Priority | Kenapa |
|---|------|----------|--------|
| 1 | **Auth** (login, JWT, refresh) | 🔴 | Kalau rusak, semua user locked out |
| 2 | **Payroll** (generate, approve, pay) | 🔴 | Kalau rusak, gaji salah |
| 3 | **Leave** (request, approve, balance deduction) | 🟡 | Kalau rusak, saldo cuti salah |

## 3. Test Strategy

### 3a. Unit Tests (Service Layer)

Mock repositories using interfaces. Test business logic in isolation.

**Package**: `internal/<module>/service/*_test.go`

### Example — Payroll State Machine Tests:

```go
func TestIsValidStatusTransition(t *testing.T) {
    tests := []struct {
        name     string
        from     string
        to       string
        expected bool
    }{
        {"DRAFT to APPROVED", "DRAFT", "APPROVED", true},
        {"DRAFT to PAID", "DRAFT", "PAID", false},
        {"APPROVED to PAID", "APPROVED", "PAID", true},
        {"PAID to anything", "PAID", "DRAFT", false},
        {"PAID is immutable", "PAID", "APPROVED", false},
        {"CANCELLED is immutable", "CANCELLED", "DRAFT", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := isValidStatusTransition(tt.from, tt.to)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 3b. Integration Tests (Repository Layer)

Test SQL queries against a real PostgreSQL test database.

**Setup**: Use Docker container or test database with migrations applied.

**Package**: `internal/<module>/repository/*_test.go`

### 3c. Test File Structure

```
internal/
├── auth/
│   └── service/
│       └── auth_service_test.go         # login, JWT validation, refresh
├── payroll/
│   └── service/
│       ├── payroll_service_test.go       # GenerateBulk, UpdateStatus
│       └── payroll_transition_test.go    # state machine tests
├── leave/
│   └── service/
│       └── leave_service_test.go        # CreateRequest, Approve, Balance deduction
└── shared/
    └── helper/
        └── helper_test.go              # utility function tests
```

## 4. Implementation Steps

### Step 1: Setup Test Infrastructure

**File**: `internal/testutil/setup.go` [NEW]

```go
package testutil

import (
    "context"
    "testing"
    
    "github.com/jackc/pgx/v5/pgxpool"
)

// NewTestDB creates a connection to test database
func NewTestDB(t *testing.T) *pgxpool.Pool {
    t.Helper()
    pool, err := pgxpool.New(context.Background(), 
        "postgres://test:test@localhost:5432/hris_test?sslmode=disable")
    if err != nil {
        t.Fatalf("failed to connect to test db: %v", err)
    }
    t.Cleanup(func() { pool.Close() })
    return pool
}

// CleanupTables truncates specified tables for test isolation
func CleanupTables(t *testing.T, pool *pgxpool.Pool, tables ...string) {
    t.Helper()
    for _, table := range tables {
        _, err := pool.Exec(context.Background(), "TRUNCATE TABLE "+table+" CASCADE")
        if err != nil {
            t.Fatalf("failed to truncate %s: %v", table, err)
        }
    }
}
```

### Step 2: Create Mock Repositories

**File**: `internal/testutil/mocks.go` [NEW]

Create mock implementations of repository interfaces for unit tests.

### Step 3: Write Auth Tests

**File**: `internal/auth/service/auth_service_test.go` [NEW]

Test cases:
- ✅ Login with valid credentials
- ❌ Login with wrong password
- ❌ Login with non-existent email
- ✅ JWT token generation and validation
- ✅ Refresh token rotation
- ❌ Expired refresh token

### Step 4: Write Payroll Tests

**File**: `internal/payroll/service/payroll_service_test.go` [NEW]

Test cases:
- ✅ GenerateBulk creates payrolls for all employees
- ❌ GenerateBulk with no employees
- ✅ Status transition DRAFT → APPROVED
- ❌ Status transition DRAFT → PAID (invalid)
- ✅ Status transition APPROVED → PAID
- ❌ Edit APPROVED payroll (should fail)
- ✅ Salary calculation accuracy

### Step 5: Write Leave Tests

**File**: `internal/leave/service/leave_service_test.go` [NEW]

Test cases:
- ✅ Create leave request deducts from balance
- ❌ Create leave request with insufficient balance
- ✅ Approve request changes status
- ✅ Reject request restores pending balance
- ❌ Request overlapping dates
- ✅ Weekend/holiday exclusion

### Step 6: Write Helper Tests

**File**: `shared/helper/helper_test.go` [NEW]

Test utility functions (date calculations, timezone handling, etc.)

## 5. Files Changed

| # | File | Change |
|---|------|--------|
| 1 | `internal/testutil/setup.go` | [NEW] Test DB setup + cleanup |
| 2 | `internal/testutil/mocks.go` | [NEW] Mock repository implementations |
| 3 | `internal/auth/service/auth_service_test.go` | [NEW] Auth flow tests |
| 4 | `internal/payroll/service/payroll_service_test.go` | [NEW] Payroll generation + state machine tests |
| 5 | `internal/leave/service/leave_service_test.go` | [NEW] Leave request + balance tests |
| 6 | `shared/helper/helper_test.go` | [NEW] Utility tests |

## 6. Verification

```bash
# Run all tests
go test ./... -v -count=1

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Target: ≥60% coverage on critical paths (auth, payroll, leave)
```

## 7. Minimum Test Count for Launch

| Module | Min Test Cases | Focus |
|--------|---------------|-------|
| Auth | 6 | Login, JWT, refresh |
| Payroll | 7 | Generate, transitions, calculations |
| Leave | 6 | Request, approve/reject, balance |
| Helpers | 4 | Date/timezone utils |
| **Total** | **23** | Critical paths covered |
