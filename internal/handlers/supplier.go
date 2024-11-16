package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"github.com/mohdjishin/order-inventory-management/util"
	"go.uber.org/zap"
)

func CreateSupplier(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	return createUser(c, req, models.SupplierRole)
}

func ListOrdersForSupplier(c fiber.Ctx) error {
	supplierId, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to extract supplier ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract supplier ID from context",
		})
	}

	var orders []models.Order
	if err := db.GetDb().Where("supplier_id = ?", uint(supplierId)).Preload("Product").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	var orderList []map[string]interface{}
	for _, order := range orders {
		orderList = append(orderList, map[string]interface{}{
			"order_id":     order.ID,
			"product_id":   order.ProductID,
			"product_name": order.Product.Name,
			"quantity":     order.Quantity,
			"total_price":  order.TotalPrice,
			"status":       order.Status,
			"created_at":   order.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orders": orderList,
	})
}

const (
	approved = "APPROVED"
	rejected = "REJECTED"
)

type ApproveRejectOrderRequest struct {
	OrderID uint   `json:"orderId" validate:"required"`
	Status  string `json:"status" validate:"required,oneof=APPROVED REJECTED"`
}

func ApproveRejectOrder(c fiber.Ctx) error {
	var req ApproveRejectOrderRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	supplierID, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to extract supplier ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract supplier ID from context",
		})
	}
	validatemap, err := util.ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"status": "error",
			"fields": validatemap,
		})
	}

	var order models.Order
	if err := db.GetDb().Where("id = ? AND supplier_id = ?", req.OrderID, supplierID).First(&order).Error; err != nil {
		log.Error("Order not found or supplier not authorized", zap.Any("error", err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found or you are not authorized to approve this order",
		})
	}

	order.Status = req.Status
	if err := db.GetDb().Save(&order).Error; err != nil {
		log.Error("Failed to update order status", zap.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}
	if order.Status == approved {
		var inventory models.Inventory
		if err := db.GetDb().Where("product_id = ?", order.ProductID).First(&inventory).Error; err != nil {
			log.Error("Inventory not found for the product", zap.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Inventory not found for the product",
			})
		}
		inventory.Stock = inventory.Stock - order.Quantity

		if err := db.GetDb().Save(&inventory).Error; err != nil {
			log.Error("Failed to update inventory stock", zap.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update inventory stock",
			})
		}
		// TODO: make it dynamicly change the price of the product based on the quantity available
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order status updated successfully",
		"order":   order,
	})
}
