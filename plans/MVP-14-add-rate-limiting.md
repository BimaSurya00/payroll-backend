# MVP-14: Add Rate Limiting Middleware

## Prioritas: 🟡 IMPORTANT — Security
## Estimasi: 2 jam
## Tipe: New Feature (Middleware)

---

## Deskripsi Masalah

Tidak ada rate limiting di seluruh API. Risiko:
1. **Brute force login** — attacker bisa mencoba password tanpa batas
2. **API abuse** — bisa flood server dengan bulk requests
3. **DoS** — tidak ada proteksi resource exhaustion

## Solusi

Gunakan library `gofiber/limiter` yang sudah tersedia sebagai middleware Fiber. Implementasikan:
1. **Global rate limit** — 100 req/s per IP
2. **Auth rate limit** — 5 req/min pada `/login`, `/register`, `/refresh`
3. Store di KeyDB untuk distributed rate limiting

## File yang Diubah

### 1. Tambah dependency

```bash
go get github.com/gofiber/fiber/v2/middleware/limiter
```

### 2. [NEW] `middleware/rate_limiter.go`

```go
package middleware

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
    "hris/shared/helper"
)

// GlobalRateLimiter — 100 requests per menit per IP
func GlobalRateLimiter() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:               100,
        Expiration:        1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return helper.ErrorResponse(c, fiber.StatusTooManyRequests,
                "Too many requests. Please try again later.", nil)
        },
    })
}

// AuthRateLimiter — 5 requests per menit per IP (untuk login, register)
func AuthRateLimiter() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:               5,
        Expiration:        1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP() + ":" + c.Path()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return helper.ErrorResponse(c, fiber.StatusTooManyRequests,
                "Too many attempts. Please wait before trying again.", nil)
        },
    })
}

// PayrollRateLimiter — 3 requests per menit (heavy operation)
func PayrollRateLimiter() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:               3,
        Expiration:        1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP() + ":payroll"
        },
        LimitReached: func(c *fiber.Ctx) error {
            return helper.ErrorResponse(c, fiber.StatusTooManyRequests,
                "Too many payroll generation requests. Please wait.", nil)
        },
    })
}
```

### 3. [MODIFY] `main.go`

**Tambah global rate limiter sebelum routing:**
```go
// Setelah app := fiber.New(...)
app.Use(middleware.GlobalRateLimiter())
```

### 4. [MODIFY] `internal/auth/routes.go`

**Tambah AuthRateLimiter ke public routes:**
```go
import "hris/middleware"

// Public routes — dengan rate limiter
auth.Post("/register", middleware.AuthRateLimiter(), authHandler.Register)
auth.Post("/login", middleware.AuthRateLimiter(), authHandler.Login)
auth.Post("/refresh", middleware.AuthRateLimiter(), authHandler.RefreshToken)
```

### 5. [MODIFY] `internal/payroll/routes.go`

**Tambah PayrollRateLimiter ke generate endpoint:**
```go
admin.Post("/generate", middleware.PayrollRateLimiter(), payrollHandler.GenerateBulk)
```

## Verifikasi

1. `go build ./...` — compile sukses
2. Test login rate limit:
   - Hit `POST /api/v1/auth/login` 6 kali dalam 1 menit → request ke-6 harus return `429 Too Many Requests`
3. Test global rate limit:
   - Hit endpoint mana saja 101 kali dalam 1 menit → return `429`
4. Pastikan rate limit reset setelah 1 menit
