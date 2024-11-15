package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterSupplierRoutes(supplierGroup fiber.Router) {

	productGroup := supplierGroup.Group("/product")
	productGroup.Use(middleware.AuthMiddleware)
	productGroup.Post("/", handlers.AddProduct)

}
