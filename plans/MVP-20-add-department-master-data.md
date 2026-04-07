# MVP-20: Add Department Master Data

## Prioritas: 🟢 NICE TO HAVE → 🟡 IMPORTANT untuk Reporting
## Estimasi: 2 jam
## Tipe: New Feature + Refactor

---

## Deskripsi Masalah

Saat ini `division` di employee hanya free-text string. Masalah:
- Inkonsistensi data: "IT", "Information Technology", "it" = 3 department berbeda
- Tidak bisa report per department
- Tidak bisa filter attendance/payroll per department
- Tidak ada manajemen structure organisasi

## Solusi

Buat master data `departments` + ubah employee `division` dari free-text ke foreign key.

## File yang Diubah

### 1. [NEW] Database migration: `000010_add_departments.up.sql`

```sql
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    head_employee_id UUID REFERENCES employees(id), -- Kepala departemen
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_departments_code ON departments(code);

-- Migrate existing divisions to departments
INSERT INTO departments (name, code)
SELECT DISTINCT division, UPPER(REPLACE(division, ' ', '_'))
FROM employees
WHERE division IS NOT NULL AND division != ''
ON CONFLICT (code) DO NOTHING;

-- Add department_id column to employees
ALTER TABLE employees ADD COLUMN department_id UUID REFERENCES departments(id);

-- Populate department_id from existing division data
UPDATE employees e
SET department_id = d.id
FROM departments d
WHERE d.name = e.division;
```

### 2. [NEW] `000010_add_departments.down.sql`

```sql
ALTER TABLE employees DROP COLUMN IF EXISTS department_id;
DROP TABLE IF EXISTS departments;
```

### 3. [NEW] `internal/department/` — Module baru

Structure:
```
internal/department/
├── entity/department.go
├── repository/department_repository.go
├── repository/department_repository_impl.go
├── handler/department_handler.go
├── dto/department_dto.go
└── routes.go
```

**Entity:**
```go
type Department struct {
    ID             uuid.UUID `json:"id" db:"id"`
    Name           string    `json:"name" db:"name"`
    Code           string    `json:"code" db:"code"`
    Description    string    `json:"description" db:"description"`
    HeadEmployeeID *string   `json:"headEmployeeId,omitempty" db:"head_employee_id"`
    IsActive       bool      `json:"isActive" db:"is_active"`
    CreatedAt      time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}
```

**Routes (Admin/SuperUser only):**
```go
departments.Get("/", departmentHandler.GetAll)
departments.Post("/", departmentHandler.Create)
departments.Get("/:id", departmentHandler.GetByID)
departments.Patch("/:id", departmentHandler.Update)
departments.Delete("/:id", departmentHandler.Delete)
```

### 4. [MODIFY] `internal/employee/entity/employee.go`

**Tambah field:**
```go
type Employee struct {
    // ... existing fields ...
    DepartmentID *uuid.UUID `json:"departmentId,omitempty" db:"department_id"` // NEW
}
```

### 5. [MODIFY] `internal/employee/dto/create_employee.go` dan `update_employee.go`

**Tambah field:**
```go
DepartmentID *string `json:"departmentId" validate:"omitempty,uuid"`
```

### 6. [MODIFY] Employee service dan repository

- Update `CreateEmployee` — accept dan store `department_id`
- Update `FindAll` — JOIN dengan departments untuk get department name
- Tetap keep `division` field untuk backward compatibility (deprecated)

### 7. [MODIFY] Employee response DTO

**Tambah department info:**
```go
type EmployeeResponse struct {
    // ... existing fields ...
    DepartmentID   *string `json:"departmentId,omitempty"`
    DepartmentName *string `json:"departmentName,omitempty"`
}
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Run migration — existing divisions auto-migrated ke departments
3. `GET /api/v1/departments` → return list departments
4. Create employee with `departmentId` → stored correctly
5. Employee response include `departmentName`
6. Filter employee by department → works
