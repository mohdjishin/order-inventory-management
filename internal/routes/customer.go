package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterCustomerRoutes(customerGroup fiber.Router) {

	//  list products
	customerGroup.Use(middleware.AuthMiddleware)
	customerGroup.Get("/products", handlers.ListProductsCustomer)
	customerGroup.Post("/buy", handlers.BuyProduct)

}
