package payroll

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	attendancerepository "hris/internal/attendance/repository"
	auditrepository "hris/internal/audit/repository"
	auditservice "hris/internal/audit/service"
	employeerepository "hris/internal/employee/repository"
	"hris/internal/payroll/handler"
	payrollrepository "hris/internal/payroll/repository"
	payrollService "hris/internal/payroll/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler, rateLimiterStorage fiber.Storage) {
	// Initialize repositories
	payrollRepo := payrollrepository.NewPayrollRepository(postgresDB.Pool)
	payrollConfigRepo := payrollrepository.NewPayrollConfigRepository(postgresDB.Pool)
	employeeRepo := employeerepository.NewEmployeeRepository(postgresDB.Pool)
	attendanceRepo := attendancerepository.NewAttendanceRepository(postgresDB.Pool)
	auditRepo := auditrepository.NewAuditRepository(postgresDB.Pool)

	// Initialize services
	auditService := auditservice.NewAuditService(auditRepo)
	payrollService := payrollService.NewPayrollService(payrollRepo, employeeRepo, attendanceRepo, payrollConfigRepo, auditService, postgresDB.Pool)

	// Initialize handler
	payrollHandler := handler.NewPayrollHandler(payrollService)

	// Public routes (with JWT)
	api := app.Group("/api/v1/payrolls")
	api.Use(jwtAuth)

	// Self-service routes — all authenticated users (must be before /:id)
	api.Get("/my", payrollHandler.GetMyPayrolls)
	api.Get("/my/:id", payrollHandler.GetMyPayrollByID)

	// Admin only routes - ADMIN only
	admin := api.Group("", middleware.HasRole(constants.RoleAdmin))

	admin.Post("/generate", middleware.PayrollRateLimiter(rateLimiterStorage), payrollHandler.GenerateBulk)
	admin.Get("/", payrollHandler.GetAllPayrolls)
	admin.Get("/export/csv", payrollHandler.ExportCSV)
	admin.Patch("/:id/status", payrollHandler.UpdateStatus)

	// All authenticated users can view their own payroll details
	api.Get("/:id", payrollHandler.GetPayrollByID)
}
