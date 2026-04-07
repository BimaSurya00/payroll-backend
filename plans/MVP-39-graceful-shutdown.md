# MVP-39: Graceful Shutdown

**Estimasi**: 2 jam  
**Impact**: MEDIUM — Production Reliability  
**Prerequisite**: None

---

## 1. Problem Statement

[`main.go`](file:///home/bima/Documents/hris/main.go) uses `app.Listen(addr)` which blocks indefinitely. On `SIGTERM` (Docker stop, k8s pod termination), the process is killed immediately:

- Active HTTP requests are dropped mid-response
- Database connections are not properly closed
- KeyDB connections are not flushed
- File uploads to MinIO may be corrupted

**Current main.go** (lines 137–141):
```go
// Start server
addr := cfg.App.Host + ":" + cfg.App.Port
zap.L().Info("🚀 Server starting", zap.String("addr", addr))
if err := app.Listen(addr); err != nil {
    zap.L().Fatal("failed to start server", zap.Error(err))
}
```

No signal handling, no graceful drain.

---

## 2. Solution

Use Go's `os/signal` to listen for `SIGTERM`/`SIGINT`, then call Fiber's `app.ShutdownWithTimeout()` to gracefully drain connections before closing database pools.

---

## 3. Implementation Steps

### Step 1: Modify `main.go`

**File**: [`main.go`](file:///home/bima/Documents/hris/main.go)

**Add to imports:**
```go
import (
    // ... existing imports ...
    "context"
    "os/signal"
    "syscall"
    "time"
)
```

**Replace the "Start server" block** (lines 136–141) with:

```go
// Start server in goroutine
addr := cfg.App.Host + ":" + cfg.App.Port
go func() {
    zap.L().Info("🚀 Server starting", zap.String("addr", addr))
    if err := app.Listen(addr); err != nil {
        zap.L().Fatal("failed to start server", zap.Error(err))
    }
}()

// Graceful shutdown
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
sig := <-quit
zap.L().Info("🛑 Shutdown signal received", zap.String("signal", sig.String()))

// Give active requests time to complete (max 30 seconds)
shutdownTimeout := 30 * time.Second
zap.L().Info("⏳ Shutting down server...", zap.Duration("timeout", shutdownTimeout))

if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
    zap.L().Error("Server forced to shutdown", zap.Error(err))
}

// Close database connections
zap.L().Info("🔌 Closing database connections...")
postgres.Close()
keydb.Close()

zap.L().Info("✅ Server exited gracefully")
```

> **Note**: `postgres.Close()` and `keydb.Close()` are already deferred in the current code. Remove the `defer` statements and call them explicitly here, after the server has stopped. This ensures the order is: stop accepting requests → drain active requests → close DB connections.

### Step 2: Remove or Keep `defer` Statements

**Current** (around lines 55–70):
```go
defer postgres.Close()
defer keydb.Close()
```

**Option A** (recommended): Remove the `defer` calls and rely on explicit close in the shutdown block.  
**Option B**: Keep `defer` as safety net — they'll fire when `main()` returns after the shutdown block runs. This is fine since `Close()` is idempotent.

**Recommended**: Keep `defer` as safety net (Option B). The explicit close in the shutdown block ensures deterministic ordering, while `defer` catches edge cases.

### Step 3: Update Health Check (Optional Enhancement)

**Current** (lines 108–114):
```go
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status":  "ok",
        "message": "Server is running",
    })
})
```

**Optionally** add a readiness flag that flips during shutdown:

```go
var isReady = true // package level or in a struct

app.Get("/health", func(c *fiber.Ctx) error {
    if !isReady {
        return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
            "status":  "shutting_down",
            "message": "Server is shutting down",
        })
    }
    return c.JSON(fiber.Map{
        "status":  "ok",
        "message": "Server is running",
    })
})
```

Then set `isReady = false` before calling `app.ShutdownWithTimeout()`.

> This is optional but useful for Kubernetes readiness probes and load balancers.

---

## 4. Files Changed Summary

| # | File | Change |
|---|------|--------|
| 1 | `main.go` | Add signal handling, graceful shutdown, explicit connection close, add `os/signal`+`syscall`+`time` imports |

---

## 5. Verification Plan

### Build
```bash
go build ./...
```

### Manual Test
```bash
# Terminal 1: Start server
./hris-server

# Terminal 2: Send SIGTERM
kill -SIGTERM <pid>
# OR
kill -SIGINT <pid>  # or Ctrl+C

# Expected output in Terminal 1:
# 🛑 Shutdown signal received signal=terminated
# ⏳ Shutting down server... timeout=30s
# 🔌 Closing database connections...
# ✅ Server exited gracefully
```

### Docker Test
```bash
docker stop <container_id>
# Container should stop within 30s (not immediate kill)
# Logs should show graceful shutdown messages
```

---

## 6. Edge Cases

| Scenario | Behavior |
|----------|----------|
| No active requests | Shutdown completes immediately |
| Long-running request (>30s) | Force-killed after timeout |
| Multiple SIGTERM signals | Second signal is ignored (channel buffered at 1) |
| Database already closed | `Close()` is idempotent — no error |
