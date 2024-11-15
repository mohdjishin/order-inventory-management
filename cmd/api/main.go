package main

import (
	"github.com/mohdjishin/order-inventory-management/config"
	_ "github.com/mohdjishin/order-inventory-management/config"
	"github.com/mohdjishin/order-inventory-management/db/migrations"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	log "github.com/mohdjishin/order-inventory-management/logger"
)

func main() {
	// Load configuration
	// config.LoadConfig()

	// Create Fiber app
	if err := migrations.Run(); err != nil {
		log.Fatal().Err(err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// User Routes
	// user := app.Group("/user")
	// handlers.UserRoutes(user)

	// Admin Routes
	// admin := app.Group("/admin")
	// admin.Use(middlewares.AdminAuthMiddleware) // Apply middleware for admin routes
	// handlers.AdminRoutes(admin)

	// Start server
	log.Info().Msg("Starting server")
	log.Info().Msg("Server started on port " + config.Get().Port)
	log.Fatal().Err(app.Listen(config.Get().Port))

}
