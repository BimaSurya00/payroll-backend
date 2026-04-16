package attendance

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/attendance/handler"
	"hris/internal/attendance/repository"
	"hris/internal/attendance/service"
	companyRepo "hris/internal/company/repository"
	employeeRepo "hris/internal/employee/repository"
	scheduleRepo "hris/internal/schedule/repository"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(
	app *fiber.App,
	postgresDB *database.Postgres,
	jwtAuth fiber.Handler,
) {
	attendanceRepo := repository.NewAttendanceRepository(postgresDB.Pool)
	employeeRepo := employeeRepo.NewEmployeeRepository(postgresDB.Pool)
	scheduleRepo := scheduleRepo.NewScheduleRepository(postgresDB.Pool)
	companyRepo := companyRepo.NewCompanyRepository(postgresDB.Pool)

	attendanceService := service.NewAttendanceService(
		attendanceRepo,
		employeeRepo,
		scheduleRepo,
		companyRepo,
	)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// Attendance routes - all require authentication
	attendance := app.Group("/api/v1/attendance", jwtAuth)

	// Clock in - All authenticated users (USER, ADMIN, SUPER_USER)
	attendance.Post("/clock-in",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.ClockIn,
	)

	// Clock out - All authenticated users (USER, ADMIN, SUPER_USER)
	attendance.Post("/clock-out",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.ClockOut,
	)

	// Get own history - All authenticated users (USER, ADMIN, SUPER_USER)
	attendance.Get("/history",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.GetHistory,
	)

	// Get all attendances - Admin and Super User only
	attendance.Get("/all",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.GetAllAttendances,
	)

	// Report endpoints
	attendance.Get("/report/monthly",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.GetMonthlyReport,
	)

	attendance.Get("/report/my",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.GetMyMonthlySummary,
	)

	// Correction endpoints - Admin and Super User only
	attendance.Post("/correction",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.CreateCorrection,
	)

	attendance.Patch("/:id/correction",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		attendanceHandler.UpdateCorrection,
	)
}
