package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/helper"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
	sharedHelper "github.com/itsahyarr/go-fiber-boilerplate/shared/helper"
)

func JWTAuth(jwtCfg *config.JWTConfig) fiber.Handler {
	jwtHelper := helper.NewJWTHelper(jwtCfg)

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

		// Validate token (now using the pre‑built helper)
		claims, err := jwtHelper.ValidateAccessToken(tokenString)
		if err != nil {
			return sharedHelper.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", err.Error())
		}

		// Set user info in context
		c.Locals(constants.ContextKeyUserID, claims.UserID)
		c.Locals(constants.ContextKeyUserRole, claims.Role)

		return c.Next()
	}
}

// HasRole checks if the user has one of the required roles
// SUPER_USER automatically has access to everything
func HasRole(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user role from context (set by JWTAuth middleware)
		userRole, ok := c.Locals(constants.ContextKeyUserRole).(string)
		if !ok {
			return sharedHelper.ErrorResponse(c, fiber.StatusUnauthorized, "User role not found in context", nil)
		}

		// SUPER_USER has access to everything (best practice)
		if userRole == constants.RoleSuperUser {
			return c.Next()
		}

		// Check if user's role is in the required roles
		for _, role := range requiredRoles {
			if userRole == role {
				return c.Next()
			}
		}

		// User doesn't have required role
		return sharedHelper.ErrorResponse(
			c,
			fiber.StatusForbidden,
			"Insufficient permissions",
			map[string]any{
				"required": requiredRoles,
				"actual":   userRole,
			},
		)
	}
}

// RequireRole is an alias for HasRole for backward compatibility
func RequireRole(roles ...string) fiber.Handler {
	return HasRole(roles...)
}
