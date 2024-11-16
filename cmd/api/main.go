package main

import (
	"github.com/mohdjishin/order-inventory-management/config"
	_ "github.com/mohdjishin/order-inventory-management/config"
	"github.com/mohdjishin/order-inventory-management/db/migrations"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	"github.com/mohdjishin/order-inventory-management/internal/routes"
	"go.uber.org/zap"

	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	log "github.com/mohdjishin/order-inventory-management/logger"
)

func main() {

	if err := migrations.Run(); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Get("/info", handlers.GetVersion)
	app.Get("/health", handlers.HealthCheck)
	userGroup := app.Group("/user")
	routes.RegisterUserRoutes(userGroup)

	customer := app.Group("/customer")
	routes.RegisterCustomerRoutes(customer)
	supplier := app.Group("/supplier")
	routes.RegisterSupplierRoutes(supplier)

	adminGroup := app.Group("/admin")
	routes.RegisterAdminRoutes(adminGroup)

	// Start the server
	log.Info("Starting server")
	log.Info("Server started on port " + config.Get().Port)
	log.Fatal("Listen ", zap.Error(app.Listen(config.Get().Port)))
}
