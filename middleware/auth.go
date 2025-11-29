package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/helper"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
	sharedHelper "github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return sharedHelper.ErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header", nil)
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return sharedHelper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization format", nil)
		}

		tokenString := parts[1]

		// Validate token
		jwtHelper := helper.NewJWTHelper(&config.GlobalConfig.JWT)
		claims, err := jwtHelper.ValidateAccessToken(tokenString)
		if err != nil {
			return sharedHelper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
		}

		// Verify it's an access token (already checked in ValidateAccessToken)
		// Set user info in context
		c.Locals(constants.ContextKeyUserID, claims.UserID)
		c.Locals(constants.ContextKeyUserRole, claims.Role)

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals(constants.ContextKeyUserRole).(string)
		if !ok {
			return sharedHelper.ErrorResponse(c, fiber.StatusUnauthorized, "User role not found", nil)
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return sharedHelper.ErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions", nil)
	}
}
