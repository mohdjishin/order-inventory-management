package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

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
	pending  = "PENDING"
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

type UpdateDeliveryStatusRequest struct {
	OrderID uint   `json:"orderId" validate:"required"`
	Status  string `json:"status" validate:"required,oneof=PENDING SHIPPING DELIVERED CANCELLED"`
	Remarks string `json:"remarks"`
}

func UpdateDeliveryStatus(c fiber.Ctx) error {

	var req UpdateDeliveryStatusRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})

	}

	fieldValidateError, err := util.ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"status": "error",
			"fields": fieldValidateError,
		})
	}

	req.Status = strings.ToLower(req.Status)
	fmt.Println("req.Status", req.Status)
	userId, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to ex	tract customer ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract customer ID from context",
		})
	}

	var order models.Order
	if err := db.GetDb().Where("id = ? AND supplier_id = ?", req.OrderID, uint(userId)).First(&order).Error; err != nil {
		log.Error("Order not found or customer not authorized", zap.Any("error", err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found or you are not authorized to update this order",
		})
	}

	currentStatus := order.ShippingState
	newStatus, ok := models.ShippingStates[req.Status]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	fmt.Println("currentStatus", currentStatus)
	fmt.Println("newStatus", newStatus)

	switch {
	case currentStatus == models.SDelivered:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order already delivered",
		})
	case newStatus == models.SCancelled:
	case currentStatus == models.SPending && newStatus == models.SShipping:
	case currentStatus == models.SShipping && newStatus == models.SDelivered:
	default:
		curSts := order.ShippingState.String()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid status transition from %s to %s", curSts, req.Status),
		})
	}

	tx := db.GetDb().Begin()
	shipmentStatus := &models.ShipmentStatus{
		OrderID:        order.ID,
		Status:         newStatus,
		AdditionalInfo: req.Remarks,
	}
	if err := tx.Create(shipmentStatus).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to update shipping status", zap.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update shipping status",
		})
	}

	order.ShippingState = newStatus

	if err := tx.Save(&order).Error; err != nil {
		log.Error("Failed to update shipping status", zap.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update shipping status",
		})
	}

	if newStatus == models.SCancelled {
		var inventory models.Inventory
		if err := tx.Where("product_id = ?", order.ProductID).First(&inventory).Error; err != nil {
			tx.Rollback()
			log.Error("Inventory not found for the product", zap.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Inventory not found for the product",
			})
		}
		inventory.Stock = inventory.Stock + order.Quantity

		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			log.Error("Failed to update inventory stock", zap.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update inventory stock",
			})
		}
	}
	if err := tx.Commit().Error; err != nil {
		log.Error("Failed to commit transaction", zap.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update shipping status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Order status updated to %s successfully", newStatus.String()),
		"order":   order,
		"status":  newStatus.String(),
	})
}

func ListReturnRequests(c fiber.Ctx) error {
	userId, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to extract supplier ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract supplier ID from context",
		})
	}

	var orders []models.Order

	if err := db.GetDb().Where("supplier_id = ? AND shipping_state = ? AND return_status = ?", uint(userId), models.SDelivered, "PENDING").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	var orderList []map[string]interface{}

	for _, order := range orders {
		orderList = append(orderList, map[string]interface{}{
			"orderId":    order.ID,
			"productId":  order.ProductID,
			"quantity":   order.Quantity,
			"totalPrice": order.TotalPrice,
			"createdAt":  order.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orders": orderList,
	})
}

type ApproveRejectReturnRequest struct {
	OrderID uint   `json:"orderId" validate:"required"`
	Status  string `json:"status" validate:"required,oneof=APPROVED REJECTED"`
	Remarks string `json:"remarks"`
}

// TOBE CHECKED.
func ApproveRejectReturnRequesthandler(c fiber.Ctx) error {
	fmt.Println("ApproveRejectReturnRequesthandler")
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
	tx := db.GetDb().Begin()
	var order models.Order
	if err := tx.Where("id = ? AND supplier_id = ?", req.OrderID, supplierID).First(&order).Error; err != nil {
		log.Error("Order not found or supplier not authorized", zap.String("error", err.Error())) // Log the error message as a string
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found or you are not authorized to approve this order",
		})
	}
	if order.ReturnStatus == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Return request not found",
		})
	}
	if order.ReturnStatus != pending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Return request already processed",
		})
	}

	order.ReturnStatus = req.Status
	updateFields := map[string]interface{}{}

	updateFields["return_status"] = req.Status
	if req.Status == approved {
		updateFields["status"] = req.Status
	}
	if err := tx.Model(&models.Order{}).
		Where("supplier_id = ? AND id = ?", supplierID, req.OrderID).
		Updates(updateFields).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	if order.ReturnStatus == approved {

		var inventory models.Inventory
		if err := tx.Where("product_id = ?", order.ProductID).First(&inventory).Error; err != nil {
			tx.Rollback()
			log.Error("Inventory not found for the product", zap.String("error", err.Error())) // Log the error message as a string
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Inventory not found for the product",
			})
		}
		inventory.Stock = inventory.Stock + order.Quantity

		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			log.Error("Failed to update inventory stock", zap.String("error", err.Error())) // Log the error message as a string
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update inventory stock",
			})
		}
	}
	if err := tx.Commit().Error; err != nil {
		log.Error("Failed to commit transaction", zap.String("error", err.Error())) // Log the error message as a string
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	var message string
	if order.ReturnStatus == approved {
		message = "Return request approved successfully"
	} else {
		message = "Return request rejected successfully"
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"order":   order, // some info should be hidden keep like this as of now
	})
}
