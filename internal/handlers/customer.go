package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
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

	userId, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to extract user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract user ID from context",
		})
	}
	tx := db.GetDb().Begin()
	var product models.Product
	if err := tx.First(&product, req.ProductID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	var inventory models.Inventory
	if err := tx.Where("product_id = ?", req.ProductID).First(&inventory).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Inventory not found for the product",
		})
	}

	if inventory.Stock < req.Quantity {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient stock available",
		})
	}

	var supplier models.User
	if err := tx.First(&supplier, product.AddedBy).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Supplier not found for the product",
		})
	}

	totalPrice := float64(req.Quantity) * product.Price

	order := models.Order{
		UserID:        uint(userId),
		ProductID:     req.ProductID,
		Quantity:      req.Quantity,
		TotalPrice:    totalPrice,
		Status:        pending,
		SupplierID:    supplier.Id,
		ShippingState: models.SPending,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create the order",
		})
	}
	shipmentStatus := models.ShipmentStatus{}

	shipmentStatus.OrderID = order.ID
	shipmentStatus.Status = models.SPending
	shipmentStatus.AdditionalInfo = "Order placed successfully. Awaiting supplier approval."
	if err := tx.Create(&shipmentStatus).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create shipment status",
		})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order placed successfully. Awaiting supplier approval.",
		"order":   order,
	})
}

func ListOrdersForCustomer(c fiber.Ctx) error {
	userId, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
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

func ReturnRequest(c fiber.Ctx) error {
	order_id := c.Params("order_id")

	userId, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to extract user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract user ID from context",
		})
	}
	db := db.GetDb()
	var order models.Order
	if err := db.Where("id = ? AND user_id = ?", order_id, userId).First(&order).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found or you are not authorized to update this order",
		})
	}

	if order.ShippingState != models.SDelivered {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order is not delivered yet",
		})
	}

	order.ReturnStatus = pending
	if err := db.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Return requested successfully",
	})

}
