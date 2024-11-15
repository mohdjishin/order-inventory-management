package main

import (
	"github.com/mohdjishin/order-inventory-management/config"
	_ "github.com/mohdjishin/order-inventory-management/config"
	"github.com/mohdjishin/order-inventory-management/db/migrations"
	"github.com/mohdjishin/order-inventory-management/internal/routes"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	log "github.com/mohdjishin/order-inventory-management/logger"
)

func main() {

	// Run database migrations
	if err := migrations.Run(); err != nil {
		log.Fatal().Err(err)
	}

	// Initialize a new Fiber app
	app := fiber.New()

	// Apply logging middleware
	app.Use(logger.New())

	// User group for supplier-related routes
	userGroup := app.Group("/user")
	routes.RegisterSupplierRoutes(userGroup)

	// Admin group for admin-related routes
	adminGroup := app.Group("/admin")
	routes.RegisterAdminRoutes(adminGroup)

	// Start the server
	log.Info().Msg("Starting server")
	log.Info().Msg("Server started on port " + config.Get().Port)
	log.Fatal().Err(app.Listen(config.Get().Port))
}
