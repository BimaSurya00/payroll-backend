package main

import (
	"log"
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
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize validator
	validator.InitValidator()

	// Initialize databases
	mongoDB, err := database.NewMongoDB(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	postgres, err := database.NewPostgres(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer postgres.Close()

	keydb, err := database.NewKeyDB(cfg.KeyDB)
	if err != nil {
		log.Fatalf("Failed to connect to KeyDB: %v", err)
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
		log.Fatalf("Failed to connect to MinIO: %v", err)
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

	// Register module routes
	auth.RegisterRoutes(app, mongoDB, keydb, cfg)
	user.RegisterRoutes(app, mongoDB, minioRepo)

	// Start server
	addr := cfg.App.Host + ":" + cfg.App.Port
	log.Printf("🚀 Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}