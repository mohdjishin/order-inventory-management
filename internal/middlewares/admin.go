package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
)

// OnlyAdmin middleware checks if the user is an admin
func OnlyAdmin(c fiber.Ctx) error {
	role, ok := c.Locals(CtxRoleKey{}).(string)
	if !ok {
		log.Error("Failed to extract role from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to extract role from context",
		})
	}

	if role != models.SuperAdmin.String() {
		log.Error("User is not an admin")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is not an admin",
		})
	}

	return c.Next()
}
