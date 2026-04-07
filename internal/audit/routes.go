package audit

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/audit/handler"
	auditrepo "hris/internal/audit/repository"
	"hris/internal/audit/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize repository
	auditRepo := auditrepo.NewAuditRepository(postgresDB.Pool)

	// Initialize service
	auditService := service.NewAuditService(auditRepo)

	// Initialize handler
	auditHandler := handler.NewAuditHandler(auditService)

	// Audit routes - SUPER_USER only (security-sensitive audit logs)
	audit := app.Group("/api/v1/audit", jwtAuth, middleware.HasRole(constants.RoleSuperUser))

	audit.Get("/logs", auditHandler.GetAll)
	audit.Get("/logs/:resourceType/:resourceId", auditHandler.GetByResource)
}
