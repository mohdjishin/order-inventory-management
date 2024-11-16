package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"github.com/mohdjishin/order-inventory-management/util"
	"go.uber.org/zap"
)

func ListUserInventoryWithProduct(c fiber.Ctx) error {
	val, ok := c.Locals("userId").(float64)
	if !ok {
		log.Error("Failed to extract user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract user ID from context",
		})
	}

	userID := uint(val)

	var inventories []models.Inventory
	if err := db.GetDb().Preload("Product").Where("added_by = ?", userID).Find(&inventories).Error; err != nil {
		log.Error("Failed to fetch inventories and products for user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch inventories and products",
		})
	}

	if len(inventories) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "No inventories found",
			"data":    []models.Inventory{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Inventories and products retrieved successfully",
		"data":    inventories,
	})
}

type AddInventoryRequest struct {
	Stock       int     `json:"stock" validate:"required,gte=0"`
	ProductName string  `json:"productName" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Category    string  `json:"category" validate:"required"`
}

func AddInventoryAndProduct(c fiber.Ctx) error {
	var input AddInventoryRequest

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		log.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	userId, ok := c.Locals("userId").(float64)
	if !ok {
		log.Error("Failed to extract user ID from context", zap.Any("userId", userId))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to extract user ID from context",
		})
	}

	if validationErrors, err := util.ValidateStruct(input); err != nil {
		log.Error("Validation failed for AddInventoryRequest", zap.Any("input", input), zap.Any("errors", validationErrors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or missing fields",
			"fields":  validationErrors,
		})
	}

	product := models.Product{
		Name:        input.ProductName,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		AddedBy:     uint(userId),
	}

	tx := db.GetDb().Begin()
	if tx.Error != nil {
		log.Error("Failed to start database transaction", zap.Error(tx.Error))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process transaction",
		})
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to create product", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create product",
		})
	}

	inventory := models.Inventory{
		Stock:     input.Stock,
		AddedBy:   uint(userId),
		ProductID: product.ID,
		BasePrice: input.Price,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to create inventory", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create inventory",
		})
	}

	product.InventoryID = inventory.ID
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		log.Error("Failed to update product with inventory ID", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update product with inventory ID",
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Transaction failed to complete",
		})
	}

	log.Info("Product and Inventory added successfully", zap.Any("product", product), zap.Any("inventory", inventory))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Inventory and product added successfully",
		"data": map[string]interface{}{
			"inventory_id": inventory.ID,
			"product_id":   product.ID,
			"product_name": product.Name,
			"category":     product.Category,
			"price":        product.Price,
			"stock":        inventory.Stock,
		},
	})
}

type UpdateInventoryRequest struct {
	ProductID uint    `json:"productId" validate:"required"`
	NewStock  int     `json:"newStock" validate:"required,gt=0"`
	NewPrice  float64 `json:"newPrice" validate:"required,gt=0"`
}

func UpdateInventories(c fiber.Ctx) error {
	var req UpdateInventoryRequest
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
	validationFields, err := util.ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"status": "error",
			"fields": validationFields,
		})
	}
	if err := UpdateInventory(&req, uint(userId)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Inventory updated successfully",
	})
}
