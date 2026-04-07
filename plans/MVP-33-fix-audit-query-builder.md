# MVP-33: Fix Audit Query Builder

## Prioritas: 🟡 IMPORTANT
## Estimasi: 1 jam
## Tipe: Bug Fix

---

## Deskripsi

`audit_repository_impl.go` menggunakan `string(rune('0'+argNum))` untuk generate parameter placeholder ($1, $2, dll). Ini **BROKEN untuk argNum > 9** karena `rune('0'+10)` = `:` bukan `10`.

Juga semua query pakai `SELECT *` yang fragile.

## File yang Diubah

### [MODIFY] `internal/audit/repository/audit_repository_impl.go`

**Ganti semua `string(rune('0'+argNum))` dengan `fmt.Sprintf("$%d", argNum)`:**

```diff
-  query += ` AND user_id = $` + string(rune('0'+argNum))
+  query += fmt.Sprintf(` AND user_id = $%d`, argNum)
```

Lakukan di `FindAll()` dan `Count()` — semua instance.

**Ganti `SELECT *` dengan explicit columns:**
```diff
-  query := `SELECT * FROM audit_logs WHERE 1=1`
+  query := `SELECT id, user_id, user_name, action, resource_type, resource_id, old_data, new_data, metadata, ip_address, created_at FROM audit_logs WHERE 1=1`
```

Tambah `import "fmt"` jika belum ada.

## Verifikasi

```bash
go build ./internal/audit/...
# Then test with > 9 filter params (edge case)
```
