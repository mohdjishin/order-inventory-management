package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
)

func CreateCustomer(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	return createUser(c, req, models.CustomerRole)
}

type OrderProductRequest struct {
	ProductID uint `json:"productId" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

func OrderProduct(c fiber.Ctx) error {
	var req OrderProductRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	userId, ok := c.Locals("userId").(float64)
	if !ok {
		log.Error("Failed to extract user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract user ID from context",
		})
	}

	// Fetch the product details
	var product models.Product
	if err := db.GetDb().First(&product, req.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	var inventory models.Inventory
	if err := db.GetDb().Where("product_id = ?", req.ProductID).First(&inventory).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Inventory not found for the product",
		})
	}

	if inventory.Stock < req.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient stock available",
		})
	}

	var supplier models.User
	if err := db.GetDb().First(&supplier, product.AddedBy).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Supplier not found for the product",
		})
	}

	totalPrice := float64(req.Quantity) * product.Price

	order := models.Order{
		UserID:     uint(userId),
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
		Status:     "PENDING",
		SupplierID: supplier.Id,
	}

	if err := db.GetDb().Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create the order",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order placed successfully. Awaiting supplier approval.",
		"order":   order,
	})
}

func ListOrdersForCustomer(c fiber.Ctx) error {
	userId, ok := c.Locals("userId").(float64)
	if !ok {
		log.Error("Failed to extract user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract user ID from context",
		})
	}

	var orders []models.Order
	if err := db.GetDb().Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orders": orders,
	})
}
