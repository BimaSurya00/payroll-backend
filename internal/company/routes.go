package company

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/company/handler"
	"hris/internal/company/repository"
	"hris/internal/company/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize repository
	companyRepo := repository.NewCompanyRepository(postgresDB.Pool)

	// Initialize service
	companySvc := service.NewCompanyService(companyRepo)

	// Initialize handler
	companyHandler := handler.NewCompanyHandler(companySvc)

	// Public company routes - accessible by all authenticated users
	// For getting own company info
	publicApi := app.Group("/api/v1/companies")
	publicApi.Use(jwtAuth)
	publicApi.Get("/current", companyHandler.GetCurrent)
	publicApi.Get("/:id", companyHandler.GetByID)

	// Admin/Super User routes - company management
	adminApi := publicApi.Group("/")
	adminApi.Use(middleware.HasRole(constants.RoleSuperUser))
	adminApi.Post("/", companyHandler.Create)
	adminApi.Get("/", companyHandler.GetAll)
	adminApi.Put("/:id", companyHandler.Update)
	adminApi.Delete("/:id", companyHandler.Delete)
}
