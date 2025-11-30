package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/itsahyarr/go-fiber-boilerplate/database"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/handler"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/repository"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth/service"
	userRepo "github.com/itsahyarr/go-fiber-boilerplate/internal/user/repository"
)

func RegisterRoutes(
	app *fiber.App,
	db *database.MongoDB,
	keydb *database.KeyDB,
	cfg *config.Config,
	jwtAuth fiber.Handler,
) {
	// Initialize dependencies
	userRepository := userRepo.NewUserRepository(db)
	tokenRepository := repository.NewTokenRepository(keydb)
	authService := service.NewAuthService(userRepository, tokenRepository, &cfg.JWT)
	authHandler := handler.NewAuthHandler(authService)

	// Auth routes
	auth := app.Group("/api/v1/auth")

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	auth.Post("/logout", jwtAuth, authHandler.Logout)
	auth.Post("/logout-all", jwtAuth, authHandler.LogoutAll)
}
