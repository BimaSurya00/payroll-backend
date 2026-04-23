package dashboard

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/dashboard/handler"
	"hris/internal/dashboard/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	dashboardService := service.NewDashboardService(postgresDB.Pool)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// Superuser dashboard - MUST be registered BEFORE admin dashboard
	// because both share the /api/v1/dashboard prefix
	superDash := app.Group("/api/v1/dashboard/super", jwtAuth, middleware.HasRole(constants.RoleSuperUser))
	superDash.Get("/summary", dashboardHandler.GetSuperUserSummary)

	// Admin dashboard routes
	dash := app.Group("/api/v1/dashboard", jwtAuth, middleware.HasRole(constants.RoleAdmin))
	dash.Get("/summary", dashboardHandler.GetSummary)
	dash.Get("/attendance-stats", dashboardHandler.GetAttendanceStats)
	dash.Get("/payroll-stats", dashboardHandler.GetPayrollStats)
	dash.Get("/employee-stats", dashboardHandler.GetEmployeeStats)
	dash.Get("/recent-activities", dashboardHandler.GetRecentActivities)
}
