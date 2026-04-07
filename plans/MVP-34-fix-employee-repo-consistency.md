# MVP-34: Fix Employee Repo — Consistent Department + FullName

## Prioritas: 🟡 IMPORTANT
## Estimasi: 2 jam
## Tipe: Integration Fix (re-plan dari MVP-31 yang belum dikerjakan)

---

## Deskripsi

Employee repository queries masih inkonsisten. Hanya `FindByID` yang JOIN departments. `Employee` struct di repository tidak punya `FullName`.

## Detail Perubahan

Sama persis seperti [MVP-31](./MVP-31-fix-employee-repo-consistency.md).

- `Employee` struct: tambah `FullName string`
- `Create()`: tambah `department_id` ke INSERT query
- `Update()`: tambah `department_id` ke SET clause
- `FindByUserID()`: tambah `LEFT JOIN departments`, select `department_id`, `department_name`
- `FindAll()`: tambah `LEFT JOIN departments`, select `department_id`, `department_name`
- `FindAllWithoutPagination()`: tambah `LEFT JOIN departments`
- `FindByIDs()`: tambah `department_id`, `gender`, `full_name` ke SELECT + Scan

## Verifikasi

```bash
go build ./...
# GET /api/v1/employees → response include departmentName
# POST /api/v1/employees (with departmentId) → stored correctly
```
