# MVP-43: Request Tracing (Request ID + Structured Logging)

**Estimasi**: 3 jam  
**Impact**: 🟡 SEDANG — Production Debugging

---

## 1. Problem

Tidak ada cara menghubungkan log dari satu request. Ketika error terjadi di production:
- Log entry dari middleware, service, dan repository tercampur
- Tidak bisa trace request dari masuk sampai keluar
- Customer complain "request gagal" tapi sulit cari log-nya

## 2. Implementation

### Step 1: Create Request ID Middleware

**File**: `middleware/request_id.go` [NEW]

```go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Use existing request ID from header, or generate new one
        requestID := c.Get(RequestIDHeader)
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Set in response header
        c.Set(RequestIDHeader, requestID)
        
        // Store in locals for use in handlers/services
        c.Locals("requestID", requestID)
        
        return c.Next()
    }
}
```

### Step 2: Create Request Logger Middleware

**File**: `middleware/request_logger.go` [NEW]

```go
package middleware

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        requestID, _ := c.Locals("requestID").(string)
        
        // Process request
        err := c.Next()
        
        // Log after response
        duration := time.Since(start)
        statusCode := c.Response().StatusCode()
        
        fields := []zap.Field{
            zap.String("request_id", requestID),
            zap.String("method", c.Method()),
            zap.String("path", c.Path()),
            zap.Int("status", statusCode),
            zap.Duration("duration", duration),
            zap.String("ip", c.IP()),
            zap.String("user_agent", c.Get("User-Agent")),
        }
        
        // Add user_id if available (from JWT)
        if userID, ok := c.Locals("userID").(string); ok {
            fields = append(fields, zap.String("user_id", userID))
        }
        
        if statusCode >= 500 {
            zap.L().Error("Request failed", fields...)
        } else if statusCode >= 400 {
            zap.L().Warn("Client error", fields...)
        } else {
            zap.L().Info("Request completed", fields...)
        }
        
        return err
    }
}
```

### Step 3: Register in `main.go`

Add before other middleware:
```go
app.Use(middleware.RequestID())
app.Use(middleware.RequestLogger())
app.Use(middleware.GlobalRateLimiter(rateLimiterStorage))
```

### Step 4: Pass Request ID to Services via Context

Create a helper to extract request ID from Fiber context and inject into Go context:

**File**: `shared/helper/context.go` [NEW]

```go
package helper

import "context"

type contextKey string

const requestIDKey contextKey = "requestID"

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
```

## 3. Files Changed

| # | File | Change |
|---|------|--------|
| 1 | `middleware/request_id.go` | [NEW] Generate/propagate X-Request-ID |
| 2 | `middleware/request_logger.go` | [NEW] Structured access logs with zap |
| 3 | `shared/helper/context.go` | [NEW] Request ID context propagation |
| 4 | `main.go` | Register middleware |

## 4. Verification

```bash
go build ./...
curl -v http://localhost:3000/health
# Response headers should include: X-Request-ID: <uuid>
# Server logs should show structured JSON with request_id field
```
