package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
)

func RegisterSupplierRoutes(userGroup fiber.Router) {
	supplierGroup := userGroup.Group("/supplier")

	supplierGroup.Post("/", handlers.CreateSupplier)
}
