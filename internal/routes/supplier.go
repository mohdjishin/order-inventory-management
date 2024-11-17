package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterSupplierRoutes(supplierGroup fiber.Router) {
	supplierGroup.Use(middleware.AuthMiddleware)
	supplierGroup.Use(middleware.OnlySuppliers)
	orderSGroup := supplierGroup.Group("/orders")
	orderSGroup.Get("/", handlers.ListOrdersForSupplier)
	orderSGroup.Put("/approve-reject", handlers.ApproveRejectOrder)

	inventoryGroup := supplierGroup.Group("/inventory")

	inventoryGroup.Post("/", handlers.AddInventoryAndProduct)
	inventoryGroup.Get("/", handlers.ListUserInventoryWithProduct)
	inventoryGroup.Put("/", handlers.UpdateInventories)
	inventoryGroup.Delete("/:id", handlers.DeleteInventories)

	productGroups := supplierGroup.Group("/product")
	productGroups.Get("/with-pricing-history", handlers.GetAllProductsWithPricingHistory)

	// supplierGroup.Get("/orders", handlers.ListOrdersForSupplier)

	// productGroup.Post("/", handlers.AddProduct)
	// productGroup.Get("/", handlers.ListProducts)

}
