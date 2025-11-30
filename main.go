package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/itsahyarr/go-fiber-boilerplate/config"
	"github.com/itsahyarr/go-fiber-boilerplate/database"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/auth"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/minio"
	"github.com/itsahyarr/go-fiber-boilerplate/internal/user"
	"github.com/itsahyarr/go-fiber-boilerplate/middleware"
	"github.com/itsahyarr/go-fiber-boilerplate/shared/validator"
	"go.uber.org/zap"
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

	// Initialize databases
	mongoDB, err := database.NewMongoDB(cfg.MongoDB)
	if err != nil {
		zap.L().Fatal("failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoDB.Close()

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

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      cfg.App.Name,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPatch,
			fiber.MethodPut,
			fiber.MethodDelete,
		}, ","),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Build JWT middleware once using config
	jwtAuth := middleware.JWTAuth(&cfg.JWT)

	// Register module routes
	auth.RegisterRoutes(app, mongoDB, keydb, cfg, jwtAuth)
	user.RegisterRoutes(app, mongoDB, minioRepo, jwtAuth)

	// Start server
	addr := cfg.App.Host + ":" + cfg.App.Port
	zap.L().Info("🚀 Server starting", zap.String("addr", addr))
	if err := app.Listen(addr); err != nil {
		zap.L().Fatal("failed to start server", zap.Error(err))
	}
}

func newLogger(env string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if env == "local" || env == "development" {
		cfg = zap.NewDevelopmentConfig()
	}
	return cfg.Build()
}
