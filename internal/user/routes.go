package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/database"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/handler"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/repository"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/go-fiber-boilerplate/middleware"
)

func RegisterRoutes(app *fiber.App, db *database.MongoDB) {
	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// User routes
	users := app.Group("/api/v1/users")
	users.Use(middleware.JWTAuth()) // Protected routes

	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Patch("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
