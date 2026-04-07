package leave

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	attendanceRepo "hris/internal/attendance/repository"
	auditrepository "hris/internal/audit/repository"
	auditservice "hris/internal/audit/service"
	employeerepo "hris/internal/employee/repository"
	holidayrepo "hris/internal/holiday/repository"
	"hris/internal/leave/handler"
	"hris/internal/leave/repository"
	"hris/internal/leave/service"
	userRepo "hris/internal/user/repository"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize repositories
	leaveTypeRepo := repository.NewLeaveTypeRepository(postgresDB.Pool)
	leaveBalanceRepo := repository.NewLeaveBalanceRepository(postgresDB.Pool)
	leaveRequestRepo := repository.NewLeaveRequestRepository(postgresDB.Pool)
	employeeRepo := employeerepo.NewEmployeeRepository(postgresDB.Pool)
	userRepo := userRepo.NewUserRepository(postgresDB.Pool)
	attendanceRepo := attendanceRepo.NewAttendanceRepository(postgresDB.Pool)
	holidayRepo := holidayrepo.NewHolidayRepository(postgresDB.Pool)
	auditRepo := auditrepository.NewAuditRepository(postgresDB.Pool)

	// Initialize services
	auditService := auditservice.NewAuditService(auditRepo)
	leaveService := service.NewLeaveService(leaveTypeRepo, leaveBalanceRepo, leaveRequestRepo, employeeRepo, userRepo, attendanceRepo, holidayRepo, auditService, postgresDB.Pool)
	leaveTypeService := service.NewLeaveTypeService(leaveTypeRepo)

	// Initialize handlers
	leaveRequestHandler := handler.NewLeaveRequestHandler(leaveService)
	leaveTypeHandler := handler.NewLeaveTypeHandler(leaveTypeService)

	// Leave routes - all require authentication
	leave := app.Group("/api/v1/leave", jwtAuth)

	// ========== LEAVE TYPES (All authenticated users) ==========
	leave.Get("/types", leaveTypeHandler.GetActiveLeaveTypes)
	leave.Get("/types/all", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveTypeHandler.GetAllLeaveTypes)
	leave.Get("/types/:id", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveTypeHandler.GetLeaveTypeByID)
	leave.Post("/types", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveTypeHandler.CreateLeaveType)
	leave.Put("/types/:id", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveTypeHandler.UpdateLeaveType)
	leave.Delete("/types/:id", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveTypeHandler.DeleteLeaveType)

	// ========== LEAVE REQUESTS (All authenticated users) ==========
	leave.Post("/requests", leaveRequestHandler.CreateLeaveRequest)
	leave.Get("/requests/my", leaveRequestHandler.GetMyLeaveRequests)
	leave.Get("/balances/my", leaveRequestHandler.GetMyLeaveBalances)

	// ========== LEAVE APPROVALS (Admin & Super User only) ==========
	// IMPORTANT: Specific routes must be defined BEFORE parameterized routes
	leave.Get("/requests/pending", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveRequestHandler.GetPendingLeaveRequests)
	leave.Get("/requests", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveRequestHandler.GetAllLeaveRequests)
	leave.Get("/requests/:id", leaveRequestHandler.GetLeaveRequestByID)
	leave.Put("/requests/:id/approve", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveRequestHandler.ApproveLeaveRequest)
	leave.Put("/requests/:id/reject", middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser), leaveRequestHandler.RejectLeaveRequest)
}
