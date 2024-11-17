package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func RegisterUserRoutes(userGroup fiber.Router) {

	userGroup.Post("/login", handlers.Login)

	userGroup.Post("/supplier", handlers.CreateSupplier)
	userGroup.Post("/customer", handlers.CreateCustomer)
	userGroup.Use(middleware.AuthMiddleware)
	userGroup.Post("/change-password", handlers.ChangePassword)

}
