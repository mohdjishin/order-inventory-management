package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterAdminRoutes(app fiber.Router) {
	adminGroup := app.Group("/")
	// TODO MIDDLEWARES tobe implemented.
	adminGroup.Use(middleware.AuthMiddleware)
	adminGroup.Use(middleware.OnlyAdmin)

	adminGroup.Post("/approve-supplier/", handlers.ApproveSupplier)

	adminGroup.Get("/approved-suppliers", handlers.ListApprovedSuppliers)
	// blacklist supplier
	adminGroup.Post("/blacklist-supplier", handlers.BlacklistSupplier)

	adminGroup.Get("/non-approved-suppliers", handlers.ListNonApprovedSuppliers)
}
