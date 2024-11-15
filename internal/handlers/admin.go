package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
)

type approveSupplierRequest struct {
	Id int `json:"id"`
}

func ApproveSupplier(c fiber.Ctx) error {
	var req approveSupplierRequest

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if req.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Supplier ID is required",
		})
	}

	var user models.User
	if err := db.GetDb().First(&user, req.Id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	if user.Role != models.SupplierRole {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User is not a supplier",
		})
	}

	if user.Approved {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Supplier is already approved",
		})
	}

	user.Approved = true
	if err := db.GetDb().Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to approve supplier",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Supplier approved successfully",
	})
}

func ListApprovedSuppliers(c fiber.Ctx) error {
	var suppliers []models.User
	if err := db.GetDb().Select("id", "first_name", "last_name", "email").
		Where("approved = ?", true).
		Where("role = ?", models.SupplierRole).
		Find(&suppliers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch approved suppliers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"suppliers": suppliers,
	})
}

func ListNonApprovedSuppliers(c fiber.Ctx) error {
	var suppliers []models.User
	if err := db.GetDb().Select("id", "first_name", "last_name", "email").
		Where("approved = ?", false).
		Where("role = ?", models.SupplierRole).
		Find(&suppliers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch non-approved suppliers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"suppliers": suppliers,
	})
}
