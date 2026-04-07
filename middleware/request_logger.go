package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request ID and user ID
		requestID, _ := c.Locals("requestID").(string)
		userID, _ := c.Locals("userID").(string)

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
		if userID != "" {
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
