package holiday

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/holiday/handler"
	holidayrepo "hris/internal/holiday/repository"
	"hris/internal/holiday/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize repository
	holidayRepo := holidayrepo.NewHolidayRepository(postgresDB.Pool)

	// Initialize service
	holidayService := service.NewHolidayService(holidayRepo)

	// Initialize handler
	holidayHandler := handler.NewHolidayHandler(holidayService)

	// Holiday routes - ADMIN and SUPER_USER only
	holidays := app.Group("/api/v1/holidays", jwtAuth, middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser))

	holidays.Get("/", holidayHandler.GetAllByYear)
	holidays.Post("/", holidayHandler.Create)
	holidays.Get("/:id", holidayHandler.GetByID)
	holidays.Patch("/:id", holidayHandler.Update)
	holidays.Delete("/:id", holidayHandler.Delete)
}
