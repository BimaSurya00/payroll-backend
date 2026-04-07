# MVP-25: Complete Department — Integrate into Employee

## Prioritas: 🟡 IMPORTANT
## Estimasi: 2 jam
## Tipe: Complete Partial (MVP-20)
## Dependency: MVP-21 (main.go fix)

---

## Deskripsi Masalah

Department module sudah lengkap, employee entity sudah punya `DepartmentID`, tapi:
- Employee repository queries belum JOIN departments
- Employee service belum accept/validate departmentID
- Employee response belum include departmentName

## File yang Diubah

### 1. [MODIFY] `internal/employee/repository/employee_repository.go` (interface)

Pastikan methods sudah support department_id.

### 2. [MODIFY] `internal/employee/repository/employee_repository_impl.go`

**Update `Create()`:**
```sql
INSERT INTO employees (..., department_id) VALUES (..., $N)
```

**Update `FindByID()` dan `FindAll()`:**
```sql
SELECT e.*, d.name as department_name
FROM employees e
LEFT JOIN departments d ON e.department_id = d.id
WHERE ...
```

**Update Scan:** Tambah `&emp.DepartmentID`, `&emp.DepartmentName`

### 3. [MODIFY] `internal/employee/service/employee_service.go`

**Update CreateEmployee:**
```go
if req.DepartmentID != "" {
    did, err := uuid.Parse(req.DepartmentID)
    if err != nil { return nil, errors.New("invalid department ID") }
    employee.DepartmentID = &did
}
```

**Update UpdateEmployee:** Handle department_id update.

### 4. [MODIFY] Employee response DTOs

Pastikan `DepartmentName` di-return di response.

## Verifikasi

1. `go build ./...`
2. Create employee with departmentId → stored correctly
3. `GET /api/v1/employees/:id` → response include `departmentName`
4. `GET /api/v1/employees` → semua include `departmentName`
