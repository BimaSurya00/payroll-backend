# Section H — Architectural Decisions (RESOLVED)

**Date**: 11 February 2026  
**Context**: Resolusi 4 pertanyaan kritis dari Strategic Audit

---

## Decision 1: Multi-Tenancy Strategy

### ✅ Keputusan: **Shared Database + `company_id` Column**

**Alasan:**

| Strategy | Pro | Con | Verdict |
|----------|-----|-----|---------|
| **Shared DB + company_id** | Simple, murah, 1 deployment | Harus disiplin filter | ✅ **BEST for UMKM SaaS** |
| Schema-per-tenant | Isolasi bagus | Complex migration, 1000 tenant = 1000 schema | ❌ Overkill |
| DB-per-tenant | Isolasi total | Cost tinggi, backup nightmare | ❌ Enterprise only |

**Kenapa Shared DB + `company_id` optimal untuk kasus ini:**

1. **Target market = UMKM** — data volume kecil per tenant (< 500 employees per company). Shared DB cukup performa.
2. **Operasional simple** — 1 database backup, 1 migration path, 1 deployment.
3. **Cost rendah** — 1 instance PostgreSQL bisa handle ribuan tenant UMKM.
4. **Best practice industri** — Slack, Freshdesk, Zoho HR, semuanya mulai dari shared DB saat early stage.

**Implementasi:**

```sql
-- Buat companies table (sebelum update tabel lain)
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL, -- untuk subdomain
    is_active BOOLEAN DEFAULT TRUE,
    plan VARCHAR(50) DEFAULT 'free', -- free, starter, pro
    max_employees INT DEFAULT 25,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tambah company_id ke SEMUA tabel:
ALTER TABLE employees ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE attendances ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE payrolls ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE leave_requests ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE leave_balances ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE schedules ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE departments ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE holidays ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE payroll_configs ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
ALTER TABLE audit_logs ADD COLUMN company_id UUID NOT NULL REFERENCES companies(id);
-- leave_types tetap GLOBAL (semua tenant pakai jenis cuti yang sama)
-- overtime_policies bisa per-company
-- overtime_requests ADD company_id

-- Row-Level Security (optional tapi recommended):
ALTER TABLE employees ENABLE ROW LEVEL SECURITY;
CREATE POLICY company_isolation ON employees
    USING (company_id = current_setting('app.company_id')::uuid);
```

**Middleware pattern:**

```go
// middleware/tenant.go
func TenantMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // company_id dari JWT claims (di-set saat login)
        companyID := c.Locals("company_id").(string)
        if companyID == "" {
            return fiber.ErrUnauthorized
        }
        c.Locals("company_id", companyID)
        return c.Next()
    }
}
```

**Query pattern:**

```go
// SEBELUM (single-tenant):
query := `SELECT * FROM employees WHERE id = $1`

// SESUDAH (multi-tenant):
query := `SELECT * FROM employees WHERE id = $1 AND company_id = $2`
// company_id SELALU jadi parameter terakhir di setiap query
```

---

## Decision 2: Migration Strategy

### ✅ Keputusan: **Additive Migration + Default Company**

**Masalah:**
- 11 migration files sudah ada
- Semua tabel tanpa `company_id`
- Data existing (jika ada) harus tetap valid

**Solusi:**

```sql
-- Migration 000012_add_multi_tenancy.up.sql

-- Step 1: Buat companies table
CREATE TABLE companies (...);

-- Step 2: Insert default company untuk data existing
INSERT INTO companies (id, name, slug, plan)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Company', 'default', 'pro');

-- Step 3: Tambah column dengan DEFAULT ke existing company
ALTER TABLE employees ADD COLUMN company_id UUID REFERENCES companies(id)
    DEFAULT '00000000-0000-0000-0000-000000000001';

-- Step 4: Update semua existing rows
UPDATE employees SET company_id = '00000000-0000-0000-0000-000000000001'
    WHERE company_id IS NULL;

-- Step 5: Set NOT NULL setelah data di-populate
ALTER TABLE employees ALTER COLUMN company_id SET NOT NULL;

-- Step 6: Create composite indexes untuk performance
CREATE INDEX idx_employees_company ON employees(company_id);
CREATE INDEX idx_employees_company_user ON employees(company_id, user_id);
-- ... repeat untuk semua tabel lain
```

**Kenapa ini safe:**
- Zero downtime — `ADD COLUMN` dengan `DEFAULT` tidak lock table di PostgreSQL 11+
- Data existing langsung punya `company_id` yang valid
- `NOT NULL` constraint dijamin setelah backfill
- Backward compatible — app lama masih bisa jalan (default company)

---

## Decision 3: Mixed DB (Postgres + MongoDB) — Migrasi User ke Postgres

### ✅ Keputusan: **Pindahkan Users dari MongoDB ke PostgreSQL**

**Analisis masalah saat ini:**

```go
// employee_repository.go — ini terjadi di 4 tempat:
employee.UserName = ""   // ← SELALU kosong karena user di MongoDB
employee.UserEmail = ""  // ← SELALU kosong
```

Artinya:
1. Employee list **tidak pernah punya nama/email** — harus 2 API calls dari frontend
2. Payroll slip **tanpa nama employee** — invalid secara hukum
3. Setiap fitur yang butuh join user+employee harus **cross-database fetch** — O(n) MongoDB calls
4. Tidak bisa `JOIN` users di SQL queries — selamanya harus application-level join

**Evaluasi opsi:**

| Opsi | Pro | Con | Verdict |
|------|-----|-----|---------|
| **Keep MongoDB for users** | Sudah jalan | Cross-DB join impossible, N+1 selamanya, complexity tinggi | ❌ Tidak scalable |
| **Pindah users ke Postgres** | 1 database, proper JOINs, simple architecture | Migration effort, ubah auth repo | ✅ **BEST PRACTICE** |
| **Duplicate data (CQRS-lite)** | Fast reads | Data sync nightmare, consistency issues | ❌ Overengineered |

**Kenapa pindah ke Postgres:**

1. **User data SUDAH relational** — `user_id` adalah FK di `employees`, `audit_logs`, dll. MongoDB tidak memberikan value apapun untuk data relational.
2. **Eliminasi cross-DB problem** — Semua query bisa `JOIN users ON employees.user_id = users.id`
3. **Consistency** — 1 transaction bisa cover user creation + employee creation
4. **Simpler ops** — 1 database less to manage, backup, monitor
5. **Industry standard** — HRIS data is fundamentally relational

**MongoDB tetap dipakai untuk:**
- Token blacklist → **Pindah ke KeyDB** (sudah ada, better fit)
- Tidak ada use case MongoDB yang tersisa → **Remove MongoDB dari stack**

**Resulting stack:**
```
SEBELUM: PostgreSQL + MongoDB + KeyDB + MinIO
SESUDAH: PostgreSQL + KeyDB + MinIO
         ↑ simpler, cheaper, faster
```

**Migration plan:**

```sql
-- Migration 000012_add_users_table_postgres.up.sql

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'employee',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    profile_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, email) -- email unique PER company
);

CREATE INDEX idx_users_company ON users(company_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_company_email ON users(company_id, email);

-- Add proper FK to employees
ALTER TABLE employees ADD CONSTRAINT fk_employees_user
    FOREIGN KEY (user_id) REFERENCES users(id);
```

**Code changes:**
- `internal/user/repository/` → Rewrite dari `mongo.Collection` ke `pgxpool.Pool`
- `internal/auth/repository/` → Token storage pindah ke KeyDB (sudah ada)
- `internal/employee/repository/` → Bisa JOIN `users` table langsung — hapus `employee.UserName = ""`
- `main.go` → Remove MongoDB dependency
- `docker-compose.yml` → Remove MongoDB service

> [!IMPORTANT]
> Ini perubahan besar tapi **HARUS dilakukan sebelum launch**. Semakin banyak data di MongoDB, semakin susah migrasi.

---

## Decision 4: Definition of "Done"

### ✅ Keputusan: **DoD Checklist — 5 Level**

MVP tidak boleh dianggap "selesai" hanya karena file dibuat. Harus memenuhi 5 kriteria:

```markdown
## Definition of Done (DoD)

### Level 1: Code Complete ✍️
- [ ] Semua file yang di-plan sudah dibuat/dimodifikasi
- [ ] Tidak ada TODO/FIXME di code baru

### Level 2: Build Pass ⚙️
- [ ] `go build ./...` → zero errors
- [ ] `go vet ./...` → zero warnings
- [ ] Tidak ada unused imports

### Level 3: Integration Complete 🔗
- [ ] Module baru TERDAFTAR di main.go
- [ ] Module baru DIPANGGIL oleh consumers (bukan cuma exist)
- [ ] Routes accessible via HTTP test
- [ ] Database migration applied successfully

### Level 4: Test Pass ✅
- [ ] Unit tests pass
- [ ] Manual API test (curl/Postman) untuk happy path
- [ ] Manual API test untuk error cases (400, 401, 404, 409)

### Level 5: Verified End-to-End 🎯
- [ ] Feature berfungsi dari API → Service → Repository → DB → Response
- [ ] Data yang disimpan bisa dibaca kembali dengan benar
- [ ] Cross-module integration (jika ada) terverifikasi
```

**Enforcement:**
- Setiap MVP completion report HARUS include:
  1. Build output (`go build ./...`)
  2. Screenshot/curl log dari API test
  3. Checklist DoD yang dicentang

---

## Ringkasan Decisions

| # | Pertanyaan | Keputusan | Effort |
|---|-----------|-----------|--------|
| 1 | Multi-tenancy | Shared DB + `company_id` | 3–5 hari |
| 2 | Migration | Additive migration + default company | Include di #1 |
| 3 | Mixed DB | Pindah users ke Postgres, hapus MongoDB | 2–3 hari |
| 4 | Definition of Done | 5-level DoD checklist | Immediate |

**Total additional effort: ~1–1.5 minggu**

**Execution order yang recommended:**
```
Day 1:    Enforce DoD (immediately)
Day 1-2:  Fix compile errors (MVP-32)
Day 3-5:  Create companies table + add company_id (multi-tenancy)
Day 6-8:  Migrate users to Postgres + remove MongoDB
Day 9-10: Update all queries for multi-tenancy + user JOINs
```

Setelah ini selesai, arsitektur sudah solid untuk scale ke ribuan tenant UMKM.
