package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"github.com/mohdjishin/order-inventory-management/util"
)

type addProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	Category    string  `json:"category" validate:"required"`
}

func AddProduct(c fiber.Ctx) error {
	var input addProductRequest

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON format",
		})
	}

	if validationErrors, err := util.ValidateStruct(input); err != nil {
		log.Warn().Err(err).Msg("Validation failed for product input")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or missing fields",
			"fields":  validationErrors,
		})
	}
	val, ok := c.Locals("userId").(float64)
	if !ok {
		fmt.Println("Failed to extract user userId from context", c.Locals("userId"))
		log.Error().Msg("Failed to extract user ID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add product",
		})
	}
	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Category:    input.Category,
		AddedBy:     uint(val),
	}

	tx := db.GetDb().Begin()
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
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to create product inventory")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add product inventory",
		})
	}

	tx.Commit()

	log.Info().Msgf("Product added successfully: ID=%d, Name=%s", product.ID, product.Name)
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
