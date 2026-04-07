# MVP-27: Fix Employee Name in Reports

## Prioritas: 🟢 IMPROVEMENT
## Estimasi: 2 jam
## Tipe: Data Quality Fix

---

## Deskripsi Masalah

Beberapa endpoint menampilkan data yang salah sebagai "nama karyawan":
- Payroll GetAll → menampilkan `emp.Position` (jabatan) bukan nama
- Attendance Report → menampilkan `employee_id` bukan nama

Root cause: Employee data di PostgreSQL, User data (termasuk nama) di MongoDB. Tidak ada denormalized `full_name` di employee table.

## Solusi

Tambah field `full_name` di employee table (denormalized dari User di MongoDB). Update saat create/update employee.

## File yang Diubah

### 1. [NEW] Migration: `000011_add_employee_fullname.up.sql`

```sql
ALTER TABLE employees ADD COLUMN full_name VARCHAR(255);

-- Populate from existing position (temporary, harus diupdate dari MongoDB)
UPDATE employees SET full_name = position WHERE full_name IS NULL;
```

### 2. [MODIFY] Employee entity + create/update flow

Terima `fullName` di create employee, simpan ke DB.

### 3. [MODIFY] Payroll `GetAll` response

Gunakan `emp.FullName` alih-alih `emp.Position`.

### 4. [MODIFY] Attendance report response

Gunakan `emp.FullName` alih-alih employee_id.

## Verifikasi

1. `go build ./...`
2. `GET /api/v1/payrolls` → employeeName = nama asli (bukan jabatan)
3. `GET /api/v1/attendance/report/monthly` → employeeName = nama asli
