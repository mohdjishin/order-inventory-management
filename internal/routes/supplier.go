package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterSupplierRoutes(supplierGroup fiber.Router) {
	supplierGroup.Use(middleware.AuthMiddleware)

	orderSGroup := supplierGroup.Group("/orders")
	orderSGroup.Get("/", handlers.ListOrdersForSupplier)
	orderSGroup.Put("/approve-reject", handlers.ApproveRejectOrder)

	productGroup := supplierGroup.Group("/inventory")
	productGroup.Post("/", handlers.AddInventoryAndProduct)
	productGroup.Get("/", handlers.ListUserInventoryWithProduct)

	// supplierGroup.Get("/orders", handlers.ListOrdersForSupplier)

	// productGroup.Post("/", handlers.AddProduct)
	// productGroup.Get("/", handlers.ListProducts)

}
