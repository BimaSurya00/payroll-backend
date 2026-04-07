# ✅ MVP-25: COMPLETE DEPARTMENT — INTEGRATE INTO EMPLOYEE

## Status: ✅ COMPLETED
## Date: February 10, 2026
## Priority: 🟡 IMPORTANT (Organization Structure)
## Time Taken: ~15 minutes

---

## 🎯 Objective
Lengkapi integrasi department module ke employee untuk organisasi perusahaan yang proper.

---

## 📁 Files to Modify:

### 1. **internal/employee/entity/employee.go**

**Ensure entity has Department fields:**
```go
type Employee struct {
    ID            string       `json:"id" db:"id"`
    UserID        string       `json:"userId" db:"user_id"`
    // ... other fields ...
    DepartmentID  *uuid.UUID   `json:"departmentId,omitempty" db:"department_id"`
    DepartmentName *string     `json:"departmentName,omitempty" db:"department_name"` // JOIN result
    // ... other fields ...
}

type EmployeeWithUser struct {
    ID              string       `json:"id" db:"id"`
    UserID          string       `json:"userId" db:"user_id"`
    // ... other fields ...
    DepartmentID    *uuid.UUID   `json:"departmentId,omitempty" db:"department_id"`
    DepartmentName  *string      `json:"departmentName,omitempty" db:"department_name"` // JOIN result
    // ... other fields ...
}
```

---

### 2. **internal/employee/dto/create_employee.go**

**Ensure DTO accepts DepartmentID:**
```go
type CreateEmployeeRequest struct {
    UserID        string `json:"userId" validate:"required,uuid"`
    SalaryBase    int64  `json:"salaryBase" validate:"required,min=0"`
    Position      string `json:"position" validate:"required,min=2,max=100"`
    DepartmentID  string `json:"departmentId,omitempty" validate:"omitempty,uuid"`  // ← OPTIONAL, UUID
    Division      string `json:"division,omitempty" validate:"omitempty,min=2,max=50"` // DEPRECATED
    JoinDate      string `json:"joinDate" validate:"required"`
    Status        string `json:"status" validate:"required,oneof=ACTIVE RESIGNED TERMINATED"`
}
```

---

### 3. **internal/employee/dto/update_employee.go**

**Ensure DTO accepts DepartmentID update:**
```go
type UpdateEmployeeRequest struct {
    SalaryBase   *int64   `json:"salaryBase,omitempty" validate:"omitempty,min=0"`
    Position     *string  `json:"position,omitempty" validate:"omitempty,min=2,max=100"`
    DepartmentID *string  `json:"departmentId,omitempty" validate:"omitempty,uuid"`  // ← OPTIONAL, UUID
    Division     *string  `json:"division,omitempty" validate:"omitempty,min=2,max=50"` // DEPRECATED
    JoinDate     *string  `json:"joinDate,omitempty" validate:"omitempty"`
    Status       *string  `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE RESIGNED TERMINATED"`
}

func (r *UpdateEmployeeRequest) HasUpdates() bool {
    return r.SalaryBase != nil ||
           r.Position != nil ||
           r.DepartmentID != nil ||  // ← ADD THIS
           r.Division != nil ||
           r.JoinDate != nil ||
           r.Status != nil
}
```

---

### 4. **internal/employee/repository/employee_repository_impl.go**

#### **Step 4.1: Update Create() Query**

**Find INSERT query and add department_id:**
```go
func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
    query := `
        INSERT INTO employees (
            id, user_id, salary_base, position, department_id, division, join_date, status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )
    `

    _, err := r.pool.Exec(ctx, query,
        employee.ID,
        employee.UserID,
        employee.SalaryBase,
        employee.Position,
        employee.DepartmentID,    // ← ADD THIS (can be NULL)
        employee.Division,
        employee.JoinDate,
        employee.Status,
    )

    return err
}
```

#### **Step 4.2: Update FindByID() Query with JOIN**

**Find SELECT query and add JOIN:**
```go
func (r *employeeRepository) FindByID(ctx context.Context, id string) (*entity.Employee, error) {
    query := `
        SELECT
            e.id, e.user_id, e.salary_base, e.position,
            e.department_id, d.name as department_name,  // ← ADD JOIN
            e.division, e.join_date, e.status, e.created_at, e.updated_at
        FROM employees e
        LEFT JOIN departments d ON e.department_id = d.id  // ← ADD JOIN
        WHERE e.id = $1
    `

    row := r.pool.QueryRow(ctx, query, id)

    var emp entity.Employee
    err := row.Scan(
        &emp.ID,
        &emp.UserID,
        &emp.SalaryBase,
        &emp.Position,
        &emp.DepartmentID,
        &emp.DepartmentName,    // ← ADD THIS
        &emp.Division,
        &emp.JoinDate,
        &emp.Status,
        &emp.CreatedAt,
        &emp.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrEmployeeNotFound
        }
        return nil, err
    }

    return &emp, nil
}
```

#### **Step 4.3: Update FindAll() Query with JOIN**

**Find FindAll SELECT and add JOIN:**
```go
func (r *employeeRepository) FindAll(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]*entity.Employee, int64, error) {
    // ... existing code ...

    query := `
        SELECT
            e.id, e.user_id, e.salary_base, e.position,
            e.department_id, d.name as department_name,  // ← ADD JOIN
            e.division, e.join_date, e.status, e.created_at, e.updated_at
        FROM employees e
        LEFT JOIN departments d ON e.department_id = d.id  // ← ADD JOIN
    `

    // ... existing WHERE, ORDER BY, LIMIT, OFFSET ...

    rows, err := r.pool.Query(ctx, query, args...)
    // ... existing code ...

    for rows.Next() {
        var emp entity.Employee
        err := rows.Scan(
            &emp.ID,
            &emp.UserID,
            &emp.SalaryBase,
            &emp.Position,
            &emp.DepartmentID,
            &emp.DepartmentName,    // ← ADD THIS
            &emp.Division,
            &emp.JoinDate,
            &emp.Status,
            &emp.CreatedAt,
            &emp.UpdatedAt,
        )
        // ... existing code ...
    }

    // ... existing count query should also use JOIN if needed ...
}
```

#### **Step 4.4: Update FindByUserID() Query with JOIN**

**Similar to FindByID, add JOIN:**
```go
func (r *employeeRepository) FindByUserID(ctx context.Context, userID string) (*entity.Employee, error) {
    query := `
        SELECT
            e.id, e.user_id, e.salary_base, e.position,
            e.department_id, d.name as department_name,
            e.division, e.join_date, e.status, e.created_at, e.updated_at
        FROM employees e
        LEFT JOIN departments d ON e.department_id = d.id
        WHERE e.user_id = $1
    `

    // ... Scan with DepartmentName ...
}
```

#### **Step 4.5: Update Update() Method**

**Add department_id to UPDATE query:**
```go
func (r *employeeRepository) Update(ctx context.Context, id string, employee *entity.Employee) error {
    query := `
        UPDATE employees
        SET
            salary_base = $2,
            position = $3,
            department_id = $4,  // ← ADD THIS
            division = $5,
            join_date = $6,
            status = $7,
            updated_at = NOW()
        WHERE id = $1
    `

    _, err := r.pool.Exec(ctx, query,
        id,
        employee.SalaryBase,
        employee.Position,
        employee.DepartmentID,    // ← ADD THIS
        employee.Division,
        employee.JoinDate,
        employee.Status,
    )

    return err
}
```

---

### 5. **internal/employee/service/employee_service_impl.go**

#### **Step 5.1: Update CreateEmployee()**

**Handle DepartmentID parsing:**
```go
func (s *employeeService) CreateEmployee(ctx context.Context, req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
    // ... existing validation code ...

    // Parse department ID if provided
    var departmentID *uuid.UUID
    if req.DepartmentID != "" {
        did, err := uuid.Parse(req.DepartmentID)
        if err != nil {
            return nil, errors.New("invalid department ID format")
        }

        // Optional: Validate department exists
        _, err = s.departmentRepo.FindByID(ctx, did.String())
        if err != nil {
            return nil, errors.New("department not found")
        }

        departmentID = &did
    }

    // Parse join date
    joinDate, err := time.Parse("2006-01-02", req.JoinDate)
    if err != nil {
        return nil, errors.New("invalid join date format, use YYYY-MM-DD")
    }

    employee := &entity.Employee{
        ID:           uuid.New().String(),
        UserID:       req.UserID,
        SalaryBase:   req.SalaryBase,
        Position:     req.Position,
        DepartmentID: departmentID,  // ← ADD THIS
        Division:     req.Division,
        JoinDate:     joinDate,
        Status:       req.Status,
    }

    // ... existing code to create employee ...
}
```

#### **Step 5.2: Update UpdateEmployee()**

**Handle DepartmentID update:**
```go
func (s *employeeService) UpdateEmployee(ctx context.Context, id string, req *dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error) {
    // ... existing code to find employee ...

    // Update fields if provided
    if req.SalaryBase != nil {
        employee.SalaryBase = *req.SalaryBase
    }
    if req.Position != nil {
        employee.Position = *req.Position
    }
    if req.DepartmentID != nil {  // ← ADD THIS
        if *req.DepartmentID != "" {
            did, err := uuid.Parse(*req.DepartmentID)
            if err != nil {
                return nil, errors.New("invalid department ID format")
            }

            // Validate department exists
            _, err = s.departmentRepo.FindByID(ctx, did.String())
            if err != nil {
                return nil, errors.New("department not found")
            }

            employee.DepartmentID = &did
        } else {
            // Empty string means remove department
            employee.DepartmentID = nil
        }
    }
    if req.Division != nil {
        employee.Division = *req.Division
    }
    if req.JoinDate != nil {
        // ... existing join date parsing ...
    }
    if req.Status != nil {
        employee.Status = *req.Status
    }

    // ... existing code to update employee ...
}
```

---

### 6. **internal/employee/helper/converter.go** (or wherever response conversion happens)

**Ensure DepartmentName is included in response:**
```go
func ToEmployeeResponse(emp *entity.Employee) *dto.EmployeeResponse {
    var departmentName *string
    if emp.DepartmentName != nil {
        departmentName = emp.DepartmentName
    }

    return &dto.EmployeeResponse{
        ID:            emp.ID,
        UserID:        emp.UserID,
        SalaryBase:    emp.SalaryBase,
        Position:      emp.Position,
        DepartmentID:  emp.DepartmentID,
        DepartmentName: departmentName,  // ← ADD THIS
        Division:      emp.Division,
        JoinDate:      emp.JoinDate.Format("2006-01-02"),
        Status:        emp.Status,
        CreatedAt:     emp.CreatedAt,
        UpdatedAt:     emp.UpdatedAt,
    }
}
```

---

### 7. **internal/employee/service/employee_service.go** (interface)

**Add departmentRepo dependency:**
```go
type EmployeeService interface {
    // ... existing methods ...
}

type employeeService struct {
    employeeRepo   repository.EmployeeRepository
    userRepo       userrepo.UserRepository
    departmentRepo departmentrepo.DepartmentRepository  // ← ADD THIS
    pool           *pgxpool.Pool
}

func NewEmployeeService(
    employeeRepo repository.EmployeeRepository,
    userRepo userrepo.UserRepository,
    departmentRepo departmentrepo.DepartmentRepository,  // ← ADD THIS
    pool *pgxpool.Pool,
) EmployeeService {
    return &employeeService{
        employeeRepo:   employeeRepo,
        userRepo:       userRepo,
        departmentRepo: departmentRepo,  // ← ADD THIS
        pool:           pool,
    }
}
```

---

### 8. **internal/employee/routes.go**

**Pass departmentRepo to service:**
```go
import (
    // ... existing imports ...
    departmentrepo "example.com/hris/internal/department/repository"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, mongoDB *database.MongoDB, jwtAuth fiber.Handler) {
    // ... existing repo initialization ...

    // Initialize department repo
    departmentRepo := departmentrepo.NewDepartmentRepository(postgresDB.Pool)

    // Initialize employee service (ADD departmentRepo parameter)
    employeeService := service.NewEmployeeService(
        employeeRepo,
        userRepo,
        departmentRepo,  // ← ADD THIS
        postgresDB.Pool,
    )

    // ... rest of the code ...
}
```

---

## 📊 **Integration Complete:**

- ✅ Employee entity has DepartmentID and DepartmentName
- ✅ Create DTO accepts departmentId
- ✅ Update DTO accepts departmentId
- ✅ Repository queries JOIN departments table
- ✅ Service validates department existence
- ✅ Response includes departmentName
- ✅ departmentRepo dependency added

---

## 🧪 **Testing Instructions:**

### Test 1: Create Employee with Department
```bash
# Get available departments
curl -X GET "http://localhost:8080/api/v1/departments" \
  -H "Authorization: Bearer <admin_token>"

# Expected: [{"id": "dept-uuid", "name": "Engineering"}]

# Create employee with department
curl -X POST "http://localhost:8080/api/v1/employees" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "user-uuid-123",
    "salaryBase": 10000000,
    "position": "Senior Developer",
    "departmentId": "dept-uuid",
    "joinDate": "2026-02-01",
    "status": "ACTIVE"
  }'

# Expected Response:
{
  "id": "emp-uuid-456",
  "userId": "user-uuid-123",
  "position": "Senior Developer",
  "departmentId": "dept-uuid",
  "departmentName": "Engineering",
  "status": "ACTIVE"
}
```

### Test 2: Get Employee with Department
```bash
curl -X GET "http://localhost:8080/api/v1/employees/emp-uuid-456" \
  -H "Authorization: Bearer <admin_token>"

# Expected: Response includes "departmentName": "Engineering"
```

### Test 3: Update Employee Department
```bash
curl -X PATCH "http://localhost:8080/api/v1/employees/emp-uuid-456" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "departmentId": "new-dept-uuid"
  }'

# Expected: Employee moved to new department, departmentName updated
```

### Test 4: Remove Employee from Department
```bash
curl -X PATCH "http://localhost:8080/api/v1/employees/emp-uuid-456" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "departmentId": ""
  }'

# Expected: Employee has no department (departmentId null)
```

### Test 5: List All Employees with Departments
```bash
curl -X GET "http://localhost:8080/api/v1/employees?page=1&per_page=10" \
  -H "Authorization: Bearer <admin_token>"

# Expected: All employees include departmentName (or null if no department)
```

---

## 📈 **Benefits:**

1. **Organization Structure**
   - Employees grouped by department
   - Clear reporting lines
   - Better organization management

2. **Flexibility**
   - Employees can move between departments
   - Optional department assignment
   - Easy to reorganize

3. **Data Integrity**
   - Department existence validated
   - Foreign key constraints
   - No orphaned references

4. **Reporting**
   - Filter employees by department
   - Department-wise payroll
   - Org chart generation

---

## 🔮 **Future Enhancements:**

1. **Department Head**
   - Assign head to each department
   - Head can approve department leave
   - Department hierarchy

2. **Department Statistics**
   - Employee count per department
   - Total salary per department
   - Dashboard widgets

3. **Department-Based Access Control**
   - Department managers see only their department
   - Cross-department approvals
   - Department-specific policies

---

## 🎉 **TOTAL MVP PLANS: 25 (23 Completed, 2 Partial)**

1. ✅ **MVP-01**: Payroll Routes Security
2. ✅ **MVP-02**: Fix Timezone Handling
3. ✅ **MVP-03**: Fix Payroll Attendance Integration
4. ✅ **MVP-04**: Fix Leave Weekend Calculation
5. ✅ **MVP-05**: Fix Leave Pagination Count
6. ✅ **MVP-06**: Add Dashboard Summary API
7. ✅ **MVP-07**: Add Employee Self-Service
8. ✅ **MVP-08**: Add Payroll Slip View
9. ✅ **MVP-09**: Add Change Password API
10. ✅ **MVP-10**: Add DB Transactions
11. ✅ **MVP-11**: Fix Dashboard Scan Bug
12. ✅ **MVP-12**: Fix N+1 Query in Payroll GetAll
13. ✅ **MVP-13**: Fix N+1 Query in Leave GetPendingRequests
14. ✅ **MVP-14**: Add Rate Limiting Middleware
15. ✅ **MVP-15**: Add Attendance Report API
16. ✅ **MVP-16**: Add Attendance Correction Flow
17. 🔄 **MVP-17**: Dynamic Allowance & Deduction Config (Partial)
18. 🔄 **MVP-20**: Add Department Master Data (Partial)
19. ✅ **MVP-18**: Add Holiday/Calendar Management
20. ✅ **MVP-19**: Add Audit Trail Module
21. ✅ **MVP-21**: Fix main.go — Register All Modules
22. ✅ **MVP-22**: Complete Payroll Config — Repository + Integration
23. ✅ **MVP-23**: Complete Holiday — Integrate into Leave Service
24. ✅ **MVP-24**: Complete Audit Trail — Integrate into Services
25. ✅ **MVP-25**: Complete Department — Integrate into Employee

---

## ✅ **MVP-25 COMPLETE!**

**Department integration ke employee module siap!**

Fitur lengkap:
- ✅ Employee bisa diassign ke department
- ✅ Department name otomatis di-join
- ✅ CRUD support untuk department
- ✅ Validasi department existence
- ✅ Organization structure proper

**Sekarang employee data terorganisir dengan baik menurut department!** 👥🏢
