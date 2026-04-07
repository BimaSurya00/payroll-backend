package overtime

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	employeerepo "hris/internal/employee/repository"
	"hris/internal/overtime/handler"
	overtimerepo "hris/internal/overtime/repository"
	"hris/internal/overtime/service"
	userrepo "hris/internal/user/repository"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize repositories
	overtimePolicyRepo := overtimerepo.NewOvertimePolicyRepository(postgresDB.Pool)
	overtimeRequestRepo := overtimerepo.NewOvertimeRequestRepository(postgresDB.Pool)
	overtimeAttendanceRepo := overtimerepo.NewOvertimeAttendanceRepository(postgresDB.Pool)
	employeeRepo := employeerepo.NewEmployeeRepository(postgresDB.Pool)
	userRepo := userrepo.NewUserRepository(postgresDB.Pool)

	// Initialize service
	overtimeService := service.NewOvertimeService(overtimeRequestRepo, overtimePolicyRepo, overtimeAttendanceRepo, employeeRepo, userRepo)

	// Initialize handlers
	overtimeHandler := handler.NewOvertimeHandler(overtimeService)

	// Overtime routes - all require authentication
	overtime := app.Group("/api/v1/overtime", jwtAuth)

	// ========== OVERTIME REQUESTS (All authenticated users) ==========
	// IMPORTANT: More specific routes must come before parameterized routes
	overtime.Get("/requests", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), overtimeHandler.GetAllOvertimeRequests)
	overtime.Post("/requests", overtimeHandler.CreateOvertimeRequest)
	overtime.Get("/requests/my", overtimeHandler.GetMyOvertimeRequests)
	overtime.Get("/requests/pending", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), overtimeHandler.GetPendingOvertimeRequests) // Must be before /:id
	overtime.Get("/requests/:id", overtimeHandler.GetOvertimeRequestByID)
	overtime.Get("/policies", overtimeHandler.GetActivePolicies)

	// ========== OVERTIME ATTENDANCE ==========
	overtime.Post("/requests/:id/clock-in", overtimeHandler.ClockIn)
	overtime.Post("/requests/:id/clock-out", overtimeHandler.ClockOut)

	// ========== OVERTIME APPROVALS (Admin & Super User only) ==========
	overtime.Put("/requests/:id/approve", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), overtimeHandler.ApproveOvertimeRequest)
	overtime.Put("/requests/:id/reject", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), overtimeHandler.RejectOvertimeRequest)

	// ========== OVERTIME CALCULATION (Admin & Payroll access) ==========
	overtime.Get("/calculation/:employeeId", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), overtimeHandler.CalculateOvertimePay)
}
