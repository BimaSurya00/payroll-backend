package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/database"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/handler"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/repository"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/go-fiber-boilerplate/middleware"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/constants"
)

func RegisterRoutes(app *fiber.App, db *database.MongoDB) {
	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// User routes - all require authentication
	users := app.Group("/api/v1/users", middleware.JWTAuth())

	// Get own profile - accessible by all authenticated users (USER, ADMIN, SUPER_USER)
	users.Get("/me",
		middleware.HasRole(constants.RoleUser, constants.RoleAdmin, constants.RoleSuperUser),
		userHandler.GetOwnProfile,
	)

	// Create user - ADMIN and SUPER_USER only
	users.Post("/",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		userHandler.CreateUser,
	)

	// Get all users (paginated) - ADMIN and SUPER_USER only
	users.Get("/",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		userHandler.GetUsers,
	)

	// Get user by ID - ADMIN and SUPER_USER only
	users.Get("/:id",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		userHandler.GetUserByID,
	)

	// Update user - ADMIN and SUPER_USER only
	users.Patch("/:id",
		middleware.HasRole(constants.RoleAdmin, constants.RoleSuperUser),
		userHandler.UpdateUser,
	)

	// Delete user - SUPER_USER only
	users.Delete("/:id",
		middleware.HasRole(constants.RoleSuperUser),
		userHandler.DeleteUser,
	)
}
