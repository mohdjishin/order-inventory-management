package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterCustomerRoutes(customerGroup fiber.Router) {

	//  list products
	customerGroup.Use(middleware.AuthMiddleware)
	customerGroup.Use(middleware.OnlyCustomer)
	customerGroup.Get("/products", handlers.ListProductsCustomer)
	customerGroup.Post("/order", handlers.OrderProduct)
	customerGroup.Get("/orders", handlers.ListOrdersForCustomer)
	customerGroup.Post("orders/:order_id/return-request", handlers.ReturnRequest)

}
