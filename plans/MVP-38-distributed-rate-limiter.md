# MVP-38: Distributed Rate Limiter (KeyDB-Based)

**Estimasi**: 2 jam  
**Impact**: MEDIUM — Security in Production  
**Prerequisite**: KeyDB already configured ✅

---

## 1. Problem Statement

All 3 rate limiters in [`middleware/rate_limiter.go`](file:///home/bima/Documents/hris/middleware/rate_limiter.go) use **in-memory storage** via Fiber's built-in `limiter.New()`. This does NOT sync between instances:

- Deploy 3 instances → attacker sends 300 req/min (100 per instance)
- DDoS protection is ineffective at scale
- On restart, all rate limit counters are reset

KeyDB (Redis-compatible) is already running and configured in [`database/keydb.go`](file:///home/bima/Documents/hris/database/keydb.go) with `redis.Client`.

### Current Code (3 limiters, all in-memory)

```go
// GlobalRateLimiter — 100 req/min per IP
return limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
    ...
})

// AuthRateLimiter — 5 req/min per IP+path
return limiter.New(limiter.Config{
    Max:        5,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string { return c.IP() + ":" + c.Path() },
    ...
})

// PayrollRateLimiter — 3 req/min per IP
return limiter.New(limiter.Config{
    Max:        3,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string { return c.IP() + ":payroll" },
    ...
})
```

---

## 2. Solution: Fiber Redis Storage Adapter

Fiber has a built-in Redis storage adapter: `github.com/gofiber/storage/redis/v3`.

The `limiter.Config` has a `Storage` field that accepts `fiber.Storage`. When set, counters are stored in Redis/KeyDB instead of in-memory.

### Approach

1. Create a `fiberstorage.Storage` instance backed by KeyDB
2. Pass it to all 3 limiters via `Storage` field
3. Modify the functions to accept the storage parameter
4. Update `main.go` to create storage and pass it to middleware

---

## 3. Implementation Steps

### Step 1: Install Fiber Redis Storage

```bash
go get github.com/gofiber/storage/redis/v3
```

### Step 2: Rewrite `middleware/rate_limiter.go`

**File**: [`middleware/rate_limiter.go`](file:///home/bima/Documents/hris/middleware/rate_limiter.go)

Replace the entire file with:

```go
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"hris/shared/helper"
)

// GlobalRateLimiter — 100 requests per minute per IP
// Uses distributed storage (KeyDB) if provided, falls back to in-memory
func GlobalRateLimiter(storage ...fiber.Storage) fiber.Handler {
	cfg := limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "rl:global:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return helper.ErrorResponse(c, fiber.StatusTooManyRequests,
				"Too many requests. Please try again later.", nil)
		},
	}
	if len(storage) > 0 && storage[0] != nil {
		cfg.Storage = storage[0]
	}
	return limiter.New(cfg)
}

// AuthRateLimiter — 5 requests per minute per IP+path (login, register, refresh)
func AuthRateLimiter(storage ...fiber.Storage) fiber.Handler {
	cfg := limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "rl:auth:" + c.IP() + ":" + c.Path()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return helper.ErrorResponse(c, fiber.StatusTooManyRequests,
				"Too many attempts. Please wait before trying again.", nil)
		},
	}
	if len(storage) > 0 && storage[0] != nil {
		cfg.Storage = storage[0]
	}
	return limiter.New(cfg)
}

// PayrollRateLimiter — 3 requests per minute (heavy operation)
func PayrollRateLimiter(storage ...fiber.Storage) fiber.Handler {
	cfg := limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "rl:payroll:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return helper.ErrorResponse(c, fiber.StatusTooManyRequests,
				"Too many payroll generation requests. Please wait.", nil)
		},
	}
	if len(storage) > 0 && storage[0] != nil {
		cfg.Storage = storage[0]
	}
	return limiter.New(cfg)
}
```

**Key Changes:**
- Each function now accepts optional `fiber.Storage` variadic parameter
- KeyDB storage is used when provided, falls back to in-memory (backward compatible)
- Key prefix added (`rl:global:`, `rl:auth:`, `rl:payroll:`) to namespace in KeyDB

### Step 3: Create Storage Initialization in `main.go`

**File**: [`main.go`](file:///home/bima/Documents/hris/main.go)

**Add import:**
```go
fiberedis "github.com/gofiber/storage/redis/v3"
```

**Add after KeyDB initialization** (around line ~68, after `keydb` is created):

```go
// Initialize distributed rate limiter storage using KeyDB
rateLimiterStorage := fiberedis.New(fiberedis.Config{
    Host:     cfg.KeyDB.Host,
    Port:     cfg.KeyDB.Port,
    Password: cfg.KeyDB.Password,
    Database: cfg.KeyDB.DB + 1, // Use different DB number to isolate from token storage
    Reset:    false,
})
```

> **Note**: Use `cfg.KeyDB.DB + 1` to isolate rate limiter keys from token/session keys stored in DB 0. Alternatively, rely on key prefixes.

**Update global rate limiter call.** Find where `GlobalRateLimiter()` is called in `main.go`:

```diff
- app.Use(middleware.GlobalRateLimiter())
+ app.Use(middleware.GlobalRateLimiter(rateLimiterStorage))
```

### Step 4: Update Routes That Call Rate Limiters

Search for all calls to `AuthRateLimiter()` and `PayrollRateLimiter()` and update them:

**Auth routes** (`internal/auth/routes.go`):
```diff
- api.Post("/login", middleware.AuthRateLimiter(), authHandler.Login)
+ api.Post("/login", middleware.AuthRateLimiter(rateLimiterStorage), authHandler.Login)
```

> **Option A**: Pass `rateLimiterStorage` through `RegisterRoutes()` parameters  
> **Option B**: Make `rateLimiterStorage` a package-level variable in middleware  
> **Recommended**: Option A for explicitness

If using Option A, update `auth.RegisterRoutes` and `payroll.RegisterRoutes` to accept `fiber.Storage`:

```go
// auth/routes.go
func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, keydb *database.KeyDB, cfg *config.Config, jwtAuth fiber.Handler, rateLimiterStorage fiber.Storage) {
```

```go
// payroll/routes.go  
func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler, rateLimiterStorage fiber.Storage) {
```

Then in `main.go`:
```diff
- auth.RegisterRoutes(app, postgres, keydb, cfg, jwtAuth)
+ auth.RegisterRoutes(app, postgres, keydb, cfg, jwtAuth, rateLimiterStorage)

- payroll.RegisterRoutes(app, postgres, jwtAuth)
+ payroll.RegisterRoutes(app, postgres, jwtAuth, rateLimiterStorage)
```

### Step 5: Run `go mod tidy`

```bash
go mod tidy
```

---

## 4. Files Changed Summary

| # | File | Change |
|---|------|--------|
| 1 | `middleware/rate_limiter.go` | Accept optional `fiber.Storage`, add key prefixes |
| 2 | `main.go` | Create `fiberedis.New()` storage, pass to GlobalRateLimiter and RegisterRoutes |
| 3 | `internal/auth/routes.go` | Accept `fiber.Storage` param, pass to `AuthRateLimiter(storage)` |
| 4 | `internal/payroll/routes.go` | Accept `fiber.Storage` param, pass to `PayrollRateLimiter(storage)` |
| 5 | `go.mod` / `go.sum` | New dep: `github.com/gofiber/storage/redis/v3` |

---

## 5. Verification Plan

```bash
# Build
go build ./...

# Check no in-memory-only limiters remain
grep -n "limiter.New(" middleware/rate_limiter.go
# Should see cfg.Storage being conditionally set

# Verify KeyDB connection with rate limiter keys
# After starting the server and hitting an endpoint:
redis-cli -p 6379 KEYS "rl:*"
# Should show keys like rl:global:127.0.0.1, rl:auth:127.0.0.1:/api/v1/auth/login
```

---

## 6. Rollback

If `fiberedis.New()` fails to connect, the variadic approach gracefully falls back to in-memory (since Storage will be nil). No downtime risk.
