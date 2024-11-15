package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/handlers"
)

func RegisterUserRoutes(userGroup fiber.Router) {

	userGroup.Post("/login", handlers.Login)

	userGroup.Post("/supplier", handlers.CreateSupplier)

}
