package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	"github.com/mohdjishin/order-inventory-management/internal/routes"
)

func New() *fiber.App {
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

	return app

}
