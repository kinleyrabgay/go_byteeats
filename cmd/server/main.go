package main

import (
	"log"

	"go_byteeats/config"
	"go_byteeats/internal/database"
	"go_byteeats/internal/handlers"
	"go_byteeats/internal/repository"
	"go_byteeats/internal/services"
	"go_byteeats/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.ConnectDB(cfg)
	if err != nil {
		logger.Fatal(err, "Failed to connect to database")
	}

	// Get the underlying *sql.DB to handle connection cleanup
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal(err, "Failed to get underlying *sql.DB")
	}
	defer sqlDB.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		logger.Fatal(err, "Failed to run database migrations")
	}

	// Seed data in development environment
	if cfg.AppEnv == "development" {
		if err := database.SeedUsers(db); err != nil {
			logger.Fatal(err, "Failed to seed database")
		}
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	tokenService := services.NewTokenService()
	emailConfig := services.EmailConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		FromName: cfg.SMTPFromName,
		FromAddr: cfg.SMTPFromAddr,
	}
	emailService := services.NewEmailService(emailConfig)
	userService := services.NewUserService(userRepo, emailService, tokenService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.CustomErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// API Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes
	users := v1.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
	users.Post("/verify-email", userHandler.VerifyEmail)
	users.Post("/reset-password/initiate", userHandler.InitiatePasswordReset)
	users.Post("/reset-password/complete", userHandler.ResetPassword)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		logger.Fatal(err, "Failed to start server")
	}
}
