package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberstorage "github.com/gofiber/storage/redis/v3"
	"go.uber.org/zap"
	"hris/config"
	"hris/database"
	"hris/internal/attendance"
	"hris/internal/audit"
	"hris/internal/auth"
	"hris/internal/company"
	"hris/internal/dashboard"
	"hris/internal/department"
	"hris/internal/employee"
	"hris/internal/holiday"
	"hris/internal/leave"
	"hris/internal/minio"
	"hris/internal/overtime"
	"hris/internal/payroll"
	"hris/internal/schedule"
	"hris/internal/user"
	"hris/middleware"
	"hris/shared/email"
	sharedHelper "hris/shared/helper"
	"hris/shared/validator"
)

func main() {
	// Load configuration (still using plain stderr on failure)
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize Zap logger based on environment
	logger, err := newLogger(cfg.App.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // flush logs on exit

	// Make this logger the global logger used by zap.L() and zap.S()
	zap.ReplaceGlobals(logger)

	// Initialize validator
	validator.InitValidator()

	// Initialize timezone
	sharedHelper.InitTimezone(cfg.App.Timezone)

	// Initialize databases
	postgres, err := database.NewPostgres(cfg.Postgres)
	if err != nil {
		zap.L().Fatal("failed to connect to PostgreSQL", zap.Error(err))
	}
	defer postgres.Close()

	keydb, err := database.NewKeyDB(cfg.KeyDB)
	if err != nil {
		zap.L().Fatal("failed to connect to KeyDB", zap.Error(err))
	}
	defer keydb.Close()

	// Initialize distributed rate limiter storage using KeyDB
	port, _ := strconv.Atoi(cfg.KeyDB.Port)
	rateLimiterStorage := fiberstorage.New(fiberstorage.Config{
		Host:     cfg.KeyDB.Host,
		Port:     port,
		Password: cfg.KeyDB.Password,
		PoolSize: 10,
		Database: cfg.KeyDB.DB + 1, // Use different DB number to isolate from token storage
	})

	// Initialize MinIO
	minioClientInstance, err := minio.NewMinioClient(minio.MinioConfig{
		Endpoint:  cfg.MinIO.Endpoint,
		AccessKey: cfg.MinIO.AccessKey,
		SecretKey: cfg.MinIO.SecretKey,
		Bucket:    cfg.MinIO.Bucket,
		UseSSL:    cfg.MinIO.UseSSL,
	})
	if err != nil {
		zap.L().Fatal("failed to connect to MinIO", zap.Error(err))
	}

	minioRepo := minio.NewMinioRepository(minioClientInstance, cfg.MinIO.Endpoint)

	if cfg.Email.ResendAPIKey != "" {
		email.Init(cfg.Email)
	} else {
		zap.L().Warn("RESEND_API_KEY not set, email service disabled")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      cfg.App.Name,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.GlobalRateLimiter(rateLimiterStorage))
	// Configure CORS - FIX FOR DEVELOPMENT
	corsConfig := cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           86400,
	}

	// Jika environment variable kosong, gunakan default untuk development
	// IP WIFI: 192.168.10.163
	if corsConfig.AllowOrigins == "" {
		corsConfig.AllowOrigins = "http://192.168.10.163:5173,http://localhost:5173,http://127.0.0.1:5173"
	}

	app.Use(cors.New(corsConfig))

	// Build JWT middleware once using config
	jwtAuth := middleware.JWTAuth(&cfg.JWT)

	// Register module routes
	auth.RegisterRoutes(app, postgres, keydb, cfg, jwtAuth, rateLimiterStorage)
	user.RegisterRoutes(app, postgres, minioRepo, jwtAuth)
	department.RegisterRoutes(app, postgres, jwtAuth)

	// PostgreSQL-based modules
	employee.RegisterRoutes(app, postgres, jwtAuth)
	attendance.RegisterRoutes(app, postgres, jwtAuth)
	schedule.RegisterRoutes(app, postgres, jwtAuth)
	payroll.RegisterRoutes(app, postgres, jwtAuth, rateLimiterStorage)
	leave.RegisterRoutes(app, postgres, jwtAuth)
	overtime.RegisterRoutes(app, postgres, jwtAuth)
	dashboard.RegisterRoutes(app, postgres, jwtAuth)
	holiday.RegisterRoutes(app, postgres, jwtAuth)
	audit.RegisterRoutes(app, postgres, jwtAuth)
	company.RegisterRoutes(app, postgres, jwtAuth)

	// Start server in goroutine
	addr := cfg.App.Host + ":" + cfg.App.Port
	go func() {
		zap.L().Info("🚀 Server starting", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	zap.L().Info("🛑 Shutdown signal received", zap.String("signal", sig.String()))

	// Give active requests time to complete (max 30 seconds)
	shutdownTimeout := 30 * time.Second
	zap.L().Info("⏳ Shutting down server...", zap.Duration("timeout", shutdownTimeout))

	if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
		zap.L().Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connections
	zap.L().Info("🔌 Closing database connections...")
	postgres.Close()
	keydb.Close()

	zap.L().Info("✅ Server exited gracefully")
}

func newLogger(env string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if env == "local" || env == "development" {
		cfg = zap.NewDevelopmentConfig()
	}
	return cfg.Build()
}
