package middleware

import (
	"github.com/gofiber/fiber/v2"
	"hris/shared/constants"
	sharedHelper "hris/shared/helper"
)

// TenantGuard extracts company_id from JWT context (set by JWTAuth)
// and ensures the request is scoped to a valid tenant.
// Must be used AFTER JWTAuth middleware.
func TenantGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID, ok := c.Locals(constants.ContextKeyCompanyID).(string)
		if !ok || companyID == "" {
			return sharedHelper.ErrorResponse(c, fiber.StatusForbidden, "Company context not found", "Tenant isolation requires company_id in token")
		}

		return c.Next()
	}
}
