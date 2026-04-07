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
