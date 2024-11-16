package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	middleware "github.com/mohdjishin/order-inventory-management/internal/middlewares"
)

func ListProductsCustomer(c fiber.Ctx) error {
	type Products struct {
		Id          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Availablity string  `json:"availablity,omitempty"`
	}

	var productsWithStock []Products

	if err := db.GetDb().Table("products").
		Select(`
			products.id, 
			products.name, 
			products.description, 
			products.price, 
			inventories.stock, 
			CASE 
				WHEN inventories.stock = 0 THEN 'Out of Stock'
				ELSE 'In Stock'
			END AS availablity`).
		Joins("left join inventories on inventories.product_id = products.id").
		Scan(&productsWithStock).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve products with stock",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": productsWithStock,
	})
}

// TODO: custom marshaler for pricing history
type ProductWithHistory struct {
	ProductID      uint            `json:"product_id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Category       string          `json:"category"`
	CurrentPrice   float64         `json:"current_price"`
	PricingHistory json.RawMessage `json:"pricing_history"`
}

func GetAllProductsWithPricingHistory(c fiber.Ctx) error {
	var products []ProductWithHistory
	userID, ok := c.Locals(middleware.CtxUserIDKey{}).(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to extract user ID from context",
		})
	}
	query := `
		SELECT 
			p.id AS product_id,
			p.name,
			p.description,
			p.category,
			p.price AS current_price,
			COALESCE(
				(SELECT JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', ph.id,
						'oldPrice', ph.old_price,
						'newPrice', ph.new_price,
						'pricingDate', ph.created_at
					)
				) FROM pricing_histories ph WHERE ph.product_id = p.id), '[]') AS pricing_history
		FROM 
			products p
		WHERE 
			p.added_by = ?`

	if err := db.GetDb().Raw(query, userID).Scan(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products with pricing history",
		})
	}

	return c.JSON(products)
}
