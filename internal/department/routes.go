package department

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/department/handler"
	departmentrepo "hris/internal/department/repository"
	"hris/internal/department/service"
	employeerepo "hris/internal/employee/repository"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize repositories
	departmentRepo := departmentrepo.NewDepartmentRepository(postgresDB.Pool)
	employeeRepo := employeerepo.NewEmployeeRepository(postgresDB.Pool)

	// Initialize service
	departmentService := service.NewDepartmentService(departmentRepo, employeeRepo)

	// Initialize handler
	departmentHandler := handler.NewDepartmentHandler(departmentService)

	// Department routes - ADMIN only
	departments := app.Group("/api/v1/departments", jwtAuth, middleware.HasRole(constants.RoleAdmin))

	departments.Get("/", departmentHandler.GetAll)
	departments.Post("/", departmentHandler.Create)
	departments.Get("/:id", departmentHandler.GetByID)
	departments.Patch("/:id", departmentHandler.Update)
	departments.Delete("/:id", departmentHandler.Delete)
}
