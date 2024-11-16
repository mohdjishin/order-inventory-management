package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
)

func ListProductsCustomer(c fiber.Ctx) error {
	type ProductWithStock struct {
		Id          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	var productsWithStock []ProductWithStock

	if err := db.GetDb().Table("products").
		Select("products.id, products.name, products.description, products.price, inventories.stock").
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
