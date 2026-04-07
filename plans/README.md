# 📋 MVP Implementation Plans — HRIS SaaS

> Master index semua MVP plans. Round 1-3 sudah dieksekusi, Round 4 siap.

---

## Round 1 — ✅ COMPLETED

| # | Plan | Status |
|---|------|--------|
| 01 | Fix Payroll Routes Security | ✅ Done |
| 02 | Fix Timezone Handling | ✅ Done |
| 03 | Fix Payroll-Attendance Integration | ✅ Done |
| 04 | Fix Leave Weekend Calculation | ✅ Done |
| 05 | Fix Leave Pagination Count | ✅ Done |
| 06 | Dashboard Summary API | ✅ Done |
| 07 | Employee Self-Service Profile | ✅ Done |
| 08 | Payroll Slip per Employee | ✅ Done |
| 09 | Change Password API | ✅ Done |
| 10 | DB Transaction for Critical Ops | ✅ Done |

---

## Round 2 — ✅ COMPLETED (with partials fixed in R3)

| # | Plan | Status |
|---|------|--------|
| 11 | Fix Dashboard Scan Bug | ✅ Done |
| 12 | Fix N+1 Query Payroll GetAll | ✅ Done |
| 13 | Fix N+1 Query Leave Pending | ✅ Done |
| 14 | Rate Limiting Middleware | ✅ Done |
| 15 | Attendance Report API | ✅ Done |
| 16 | Attendance Correction Flow | ✅ Done |
| 17 | Dynamic Allowance Config | ⚠️ Config exists but not used in GenerateBulk |
| 18 | Holiday/Calendar | ⚠️ Module exists, leave struct broken |
| 19 | Audit Trail Module | ⚠️ Module exists, not integrated |
| 20 | Department Master Data | ⚠️ Module exists, queries inconsistent |

---

## Round 3 — ✅ COMPLETED (partially)

| # | Plan | Status | Issue Found |
|---|------|--------|-------------|
| 21 | Fix main.go — Register All Modules | ✅ Done | — |
| 22 | Complete Payroll Config | ⚠️ Partial | GenerateBulk still uses old calculator |
| 23 | Complete Holiday-Leave Integration | 🔴 Bug | `holidayRepo` missing from struct |
| 24 | Complete Audit Trail Integration | ❌ Not Done | Zero references in services |
| 25 | Complete Department-Employee | ⚠️ Partial | Only `FindByID` has JOIN |
| 26 | Fix Leave Error Handling | ✅ Done | — |
| 27 | Fix Employee Name in Reports | ⚠️ Partial | Entity has FullName, repo doesn't |

📄 **[Verification Report (27 MVPs)](./VERIFICATION-REPORT-ALL-27-MVPS.md)**

---

## Round 4 — 🆕 READY (Fix Remaining Integration Issues)

| # | Plan | Priority | Estimasi | File |
|---|------|----------|----------|------|
| 28 | **Fix Leave holidayRepo Struct** | 🔴 CRITICAL | 15 min | [MVP-28](./MVP-28-fix-leave-holidayrepo-struct.md) |
| 29 | **Fix Payroll Use Config Calculator** | 🔴 CRITICAL | 30 min | [MVP-29](./MVP-29-fix-payroll-use-config-calculator.md) |
| 30 | **Integrate Audit Into Services** | 🟡 IMPORTANT | 1.5 jam | [MVP-30](./MVP-30-integrate-audit-into-services.md) |
| 31 | **Fix Employee Repo Consistency** | 🟡 IMPORTANT | 2 jam | [MVP-31](./MVP-31-fix-employee-repo-consistency.md) |

### Urutan Eksekusi Round 4
```
1. MVP-28 (15 min) ← Compile error fix
2. MVP-29 (30 min) ← Payroll config integration
3. MVP-31 (2 jam)  ← Employee repo consistency
4. MVP-30 (1.5 jam) ← Audit trail integration
```

**Total: ~4 jam 15 menit**
