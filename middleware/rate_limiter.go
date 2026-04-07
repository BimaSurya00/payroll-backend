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
