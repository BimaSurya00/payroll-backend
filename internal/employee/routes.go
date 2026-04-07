package employee

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	"hris/internal/employee/handler"
	"hris/internal/employee/repository"
	"hris/internal/employee/service"
	userRepository "hris/internal/user/repository"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, jwtAuth fiber.Handler) {
	// Initialize dependencies
	employeeRepo := repository.NewEmployeeRepository(postgresDB.Pool)
	userRepo := userRepository.NewUserRepository(postgresDB.Pool)
	employeeService := service.NewEmployeeService(employeeRepo, userRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Employee routes - all require authentication
	employees := app.Group("/api/v1/employees", jwtAuth)

	// Self-service routes — ALL authenticated users (must be before /:id)
	employees.Get("/me",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.GetMyProfile,
	)
	employees.Patch("/me",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.UpdateMyProfile,
	)

	// Create employee - ADMIN and SUPER_USER only
	employees.Post("/",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.CreateEmployee,
	)

	// Get all employees (paginated) - ADMIN and SUPER_USER only
	employees.Get("/",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.GetAllEmployees,
	)

	// Get employee by ID - ADMIN and SUPER_USER only
	employees.Get("/:id",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.GetEmployeeByID,
	)

	// Update employee - ADMIN and SUPER_USER only
	employees.Patch("/:id",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.UpdateEmployee,
	)

	// Delete employee - ADMIN and SUPER_USER only
	employees.Delete("/:id",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		employeeHandler.DeleteEmployee,
	)
}
