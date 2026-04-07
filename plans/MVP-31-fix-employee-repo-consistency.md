# MVP-31: Fix Employee Repository — Consistent Department + FullName

## Prioritas: 🟡 IMPORTANT
## Estimasi: 2 jam
## Tipe: Integration Fix (Completes MVP-25 + MVP-27)

---

## Deskripsi Masalah

Employee repository memiliki 3 inkonsistensi:

1. **Department JOIN inkonsisten**: Hanya `FindByID` yang JOIN departments. `FindByUserID`, `FindAll`, `FindAllWithoutPagination`, `FindByIDs` semua **tidak JOIN departments** → tidak return `departmentName`.

2. **Create/Update tidak handle department_id**: `Create()` query tidak include `department_id` parameter ($17). `Update()` juga tidak include `department_id`.

3. **FullName field mismatch**: Entity (`entity/employee.go`) punya `FullName` field. Repository struct (`repository/employee_repository.go`) **TIDAK PUNYA** `FullName` field → `full_name` column di DB tidak pernah di-SELECT atau INSERT.

## File yang Diubah

### 1. [MODIFY] `internal/employee/repository/employee_repository.go`

**Tambah FullName ke repository structs:**
```diff
 type Employee struct {
+    FullName           string     `db:"full_name"`
     // ... existing fields ...
 }

 type EmployeeWithUser struct {
+    FullName           string     `db:"full_name"`
     // ... existing fields ...
 }
```

**Update Create() query (line 86-88):**
```diff
- INSERT INTO employees (id, user_id, position, phone_number, ..., division, created_at, updated_at)
- VALUES ($1, $2, ..., $15, $16, $17)
+ INSERT INTO employees (id, user_id, full_name, position, phone_number, ..., division, department_id, created_at, updated_at)
+ VALUES ($1, $2, $3, ..., $16, $17, $18, $19)
```

**Update Update() query (line 364-370):**
```diff
- SET position = $2, phone_number = $3, ...
+ SET full_name = $2, position = $3, phone_number = $4, ..., department_id = $15, ...
```

**Update FindByUserID() (line 176-188):** Tambah department JOIN dan department_id/department_name ke SELECT + Scan

**Update FindAll() (line 245-258):** Tambah department JOIN dan department_id/department_name ke SELECT + Scan

**Update FindAllWithoutPagination() (line 304-317):** Tambah department JOIN

**Update FindByIDs() (line 413-416):** Tambah `full_name`, `department_id`, `gender` ke SELECT + Scan

## Verifikasi

```bash
go build ./...
# Harus: ✅ SUCCESS

# Then:
# 1. GET /api/v1/employees → response include departmentName + fullName
# 2. GET /api/v1/employees/:id → departmentName + fullName
# 3. POST /api/v1/employees (with departmentId) → verify stored
# 4. GET /api/v1/payrolls → employeeName = fullName (not position)
```
