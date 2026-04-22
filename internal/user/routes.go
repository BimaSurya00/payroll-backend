package user

import (
	"github.com/gofiber/fiber/v2"
	"hris/database"
	auditrepository "hris/internal/audit/repository"
	auditservice "hris/internal/audit/service"
	employeeRepo "hris/internal/employee/repository"
	minioClient "hris/internal/minio"
	scheduleRepo "hris/internal/schedule/repository"
	"hris/internal/user/handler"
	"hris/internal/user/repository"
	"hris/internal/user/service"
	"hris/middleware"
	"hris/shared/constants"
)

func RegisterRoutes(app *fiber.App, postgresDB *database.Postgres, minioRepo minioClient.MinioRepository, jwtAuth fiber.Handler) {
	// Initialize user repo from PostgreSQL
	userRepo := repository.NewUserRepository(postgresDB.Pool)

	// Initialize employee and schedule repos for auto-create employee feature
	employeeRepository := employeeRepo.NewEmployeeRepository(postgresDB.Pool)
	scheduleRepository := scheduleRepo.NewScheduleRepository(postgresDB.Pool)

	// Initialize audit service
	auditRepo := auditrepository.NewAuditRepository(postgresDB.Pool)
	auditService := auditservice.NewAuditService(auditRepo)

	// Create service with both employee auto-creation and audit support
	userService := service.NewUserServiceWithEmployeeAndAudit(userRepo, employeeRepository, scheduleRepository, auditService)
	minioService := minioClient.NewMinioService(minioRepo, userRepo)

	userHandler := handler.NewUserHandler(userService, minioService)
	registerUserRoutes(app, userHandler, jwtAuth)
}

func registerUserRoutes(app *fiber.App, userHandler *handler.UserHandler, jwtAuth fiber.Handler) {
	// User routes - all require authentication
	users := app.Group("/api/v1/users", jwtAuth)

	// Get own profile - accessible by USER and ADMIN
	users.Get("/me",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		userHandler.GetOwnProfile,
	)

	// Create user - SUPER_USER only (platform-level user management)
	users.Post("/",
		middleware.HasRole(constants.RoleSuperUser),
		userHandler.CreateUser,
	)

	// Get all users (paginated) - SUPER_USER only
	users.Get("/",
		middleware.HasRole(constants.RoleSuperUser),
		userHandler.GetUsers,
	)

	// Get user by ID - SUPER_USER only
	users.Get("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		userHandler.GetUserByID,
	)

	// Update user - SUPER_USER only
	users.Patch("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		userHandler.UpdateUser,
	)

	// Delete user - SUPER_USER only
	users.Delete("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		userHandler.DeleteUser,
	)

	// Profile Image Routes
	// Upload profile image - USER and ADMIN
	users.Post("/:id/profile-image",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		userHandler.UploadProfileImage,
	)

	// Update profile image - USER and ADMIN
	users.Put("/:id/profile-image",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		userHandler.UpdateProfileImage,
	)

	// Delete profile image - USER and ADMIN
	users.Delete("/:id/profile-image",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin),
		userHandler.DeleteProfileImage,
	)
}
