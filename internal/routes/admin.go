package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
)

func RegisterAdminRoutes(app fiber.Router) {
	adminGroup := app.Group("/")

	adminGroup.Post("/approve-supplier/", handlers.ApproveSupplier)

	adminGroup.Get("/approved-suppliers", handlers.ListApprovedSuppliers)

	adminGroup.Get("/non-approved-suppliers", handlers.ListNonApprovedSuppliers)
}
