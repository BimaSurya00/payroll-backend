package auth

import (
	"github.com/gofiber/fiber/v2"
	"hris/config"
	"hris/database"
	"hris/internal/auth/handler"
	"hris/internal/auth/repository"
	"hris/internal/auth/service"
	userRepo "hris/internal/user/repository"
	"hris/middleware"
)

func RegisterRoutes(
	app *fiber.App,
	postgresDB *database.Postgres,
	keydb *database.KeyDB,
	cfg *config.Config,
	jwtAuth fiber.Handler,
	rateLimiterStorage fiber.Storage,
) {
	// Initialize dependencies — user repo now uses PostgreSQL
	userRepository := userRepo.NewUserRepository(postgresDB.Pool)
	tokenRepository := repository.NewTokenRepository(keydb)
	authService := service.NewAuthService(userRepository, tokenRepository, &cfg.JWT)
	authHandler := handler.NewAuthHandler(authService)

	// Auth routes
	auth := app.Group("/api/v1/auth")

	// Public routes with auth rate limiter
	auth.Post("/register", middleware.AuthRateLimiter(rateLimiterStorage), authHandler.Register)
	auth.Post("/login", middleware.AuthRateLimiter(rateLimiterStorage), authHandler.Login)
	auth.Post("/refresh", middleware.AuthRateLimiter(rateLimiterStorage), authHandler.RefreshToken)
	auth.Post("/forgot-password", middleware.AuthRateLimiter(rateLimiterStorage), authHandler.ForgotPassword)
	auth.Post("/reset-password", middleware.AuthRateLimiter(rateLimiterStorage), authHandler.ResetPassword)

	// Protected routes
	auth.Post("/logout", jwtAuth, authHandler.Logout)
	auth.Post("/logout-all", jwtAuth, authHandler.LogoutAll)
	auth.Put("/change-password", jwtAuth, authHandler.ChangePassword)
}
