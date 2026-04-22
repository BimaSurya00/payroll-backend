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

	// Clock in - USER and ADMIN only
	attendance.Post("/clock-in",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		attendanceHandler.ClockIn,
	)

	// Clock out - USER and ADMIN only
	attendance.Post("/clock-out",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		attendanceHandler.ClockOut,
	)

	// Get own history - USER and ADMIN only
	attendance.Get("/history",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		attendanceHandler.GetHistory,
	)

	// Get all attendances - ADMIN only
	attendance.Get("/all",
		middleware.HasRole(constants.RoleAdmin),
		attendanceHandler.GetAllAttendances,
	)

	// Report endpoints
	attendance.Get("/report/monthly",
		middleware.HasRole(constants.RoleAdmin),
		attendanceHandler.GetMonthlyReport,
	)

	attendance.Get("/report/my",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		attendanceHandler.GetMyMonthlySummary,
	)

	// Correction endpoints - ADMIN only
	attendance.Post("/correction",
		middleware.HasRole(constants.RoleAdmin),
		attendanceHandler.CreateCorrection,
	)

	attendance.Patch("/:id/correction",
		middleware.HasRole(constants.RoleAdmin),
		attendanceHandler.UpdateCorrection,
	)
}
