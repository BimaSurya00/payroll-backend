package schedule

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/schedule/handler"
	"hris/internal/schedule/repository"
	"hris/internal/schedule/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize dependencies
	scheduleRepo := repository.NewScheduleRepository(postgresDB.Pool)
	scheduleService := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	// Schedule routes - all require authentication
	schedules := app.Group("/api/v1/schedules", jwtAuth)

	// Get all schedules (paginated) - SUPER_USER only (platform-level config)
	schedules.Get("/",
		middleware.HasRole(constants.RoleSuperUser),
		scheduleHandler.GetAllSchedules,
	)

	// Get schedule by ID - SUPER_USER only
	schedules.Get("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		scheduleHandler.GetScheduleByID,
	)

	// Create schedule - SUPER_USER only
	schedules.Post("/",
		middleware.HasRole(constants.RoleSuperUser),
		scheduleHandler.CreateSchedule,
	)

	// Update schedule - SUPER_USER only
	schedules.Patch("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		scheduleHandler.UpdateSchedule,
	)

	// Delete schedule - SUPER_USER only
	schedules.Delete("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		scheduleHandler.DeleteSchedule,
	)
}
