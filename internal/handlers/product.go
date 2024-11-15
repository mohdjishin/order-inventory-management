package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"github.com/mohdjishin/order-inventory-management/util"
)

type AddProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	Category    string  `json:"category" validate:"required"`
}

func AddProduct(c fiber.Ctx) error {
	var input AddProductRequest

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if validationErrors, err := util.ValidateStruct(input); err != nil {
		log.Warn().Err(err).Msg("Validation failed for AddProductRequest")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or missing fields",
			"fields":  validationErrors,
		})
	}

	userID, ok := c.Locals("userId").(float64)
	if !ok {
		log.Error().Msg("Failed to retrieve user ID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process user ID",
		})
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		AddedBy:     uint(userID),
	}

	tx := db.GetDb().Begin()
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to start database transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to process transaction",
		})
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to create product")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add product",
		})
	}

	inventory := models.Inventory{
		ProductID: product.ID,
		Stock:     input.Stock,
		AddedBy:   uint(userID),
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to create product inventory")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add inventory for the product",
		})
	}

	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Transaction failed to complete",
		})
	}

	log.Info().Msgf("Product added successfully: ID=%d, Name=%s, AddedBy=%d", product.ID, product.Name, product.AddedBy)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Product added successfully",
		"data": map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"stock":       inventory.Stock,
		},
	})
}
