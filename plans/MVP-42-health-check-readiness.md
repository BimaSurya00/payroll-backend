# MVP-42: Health Check + Readiness Probe

**Estimasi**: 2 jam  
**Impact**: 🟡 SEDANG — Production Monitoring

---

## 1. Problem

Current health check (main.go lines 108–114) returns static `{"status": "ok"}` without actually checking dependencies. Load balancers and K8s probes can't detect if PostgreSQL or KeyDB is down.

## 2. Implementation

### Step 1: Create `internal/health/handler.go`

```go
package health

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

type Handler struct {
    pgPool      *pgxpool.Pool
    redisClient *redis.Client
    startTime   time.Time
}

func NewHandler(pgPool *pgxpool.Pool, redisClient *redis.Client) *Handler {
    return &Handler{
        pgPool:      pgPool,
        redisClient: redisClient,
        startTime:   time.Now(),
    }
}

// Liveness — is the process alive?
func (h *Handler) Liveness(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "ok",
        "uptime": time.Since(h.startTime).String(),
    })
}

// Readiness — can the service handle requests?
func (h *Handler) Readiness(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
    defer cancel()

    checks := fiber.Map{}
    allOk := true

    // PostgreSQL check
    if err := h.pgPool.Ping(ctx); err != nil {
        checks["postgres"] = fiber.Map{"status": "down", "error": err.Error()}
        allOk = false
    } else {
        stat := h.pgPool.Stat()
        checks["postgres"] = fiber.Map{
            "status":      "up",
            "connections": stat.TotalConns(),
            "idle":        stat.IdleConns(),
        }
    }

    // KeyDB check
    if err := h.redisClient.Ping(ctx).Err(); err != nil {
        checks["keydb"] = fiber.Map{"status": "down", "error": err.Error()}
        allOk = false
    } else {
        checks["keydb"] = fiber.Map{"status": "up"}
    }

    status := fiber.StatusOK
    statusText := "ready"
    if !allOk {
        status = fiber.StatusServiceUnavailable
        statusText = "not_ready"
    }

    return c.Status(status).JSON(fiber.Map{
        "status": statusText,
        "checks": checks,
        "uptime": time.Since(h.startTime).String(),
    })
}
```

### Step 2: Register Routes in `main.go`

Replace the current static `/health` with:

```go
healthHandler := health.NewHandler(postgres.Pool, keydb.Client)
app.Get("/health", healthHandler.Liveness)
app.Get("/health/ready", healthHandler.Readiness)
```

## 3. Files Changed

| # | File | Change |
|---|------|--------|
| 1 | `internal/health/handler.go` | [NEW] Health check with dependency verification |
| 2 | `main.go` | Replace static health check with handler |

## 4. Verification

```bash
go build ./...
curl http://localhost:3000/health           # → {"status": "ok", "uptime": "2m30s"}
curl http://localhost:3000/health/ready     # → {"status": "ready", "checks": {...}}
```
