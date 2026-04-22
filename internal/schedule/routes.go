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
	scheduleRepo := repository.NewScheduleRepository(postgresDB.Pool)
	scheduleService := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	schedules := app.Group("/api/v1/schedules", jwtAuth, middleware.HasRole(constants.RoleAdmin))

	schedules.Get("/", scheduleHandler.GetAllSchedules)
	schedules.Get("/:id", scheduleHandler.GetScheduleByID)
	schedules.Post("/", scheduleHandler.CreateSchedule)
	schedules.Patch("/:id", scheduleHandler.UpdateSchedule)
	schedules.Delete("/:id", scheduleHandler.DeleteSchedule)
}
