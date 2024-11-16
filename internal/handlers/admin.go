package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	"github.com/mohdjishin/order-inventory-management/util"
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

type BlacklistSupplierRequest struct {
	Id int `json:"id" validate:"required"`
	// WithInventory bool `json:"withInventory"`  // anyway if blacklisted, all inventory will be blacklisted
}

func BlacklistSupplier(c fiber.Ctx) error {
	//  blacklist supplier with inventory,products and orders
	var req BlacklistSupplierRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	validateFields, err := util.ValidateStruct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"status": "error",
			"fields": validateFields,
		})
	}

	//  load supplier
	tx := db.GetDb().Begin()
	var supplier models.User
	if err := tx.First(&supplier).Where("role = ?", models.SupplierRole).Where("id = ?", req.Id).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}
	if supplier.Blacklisted {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Supplier is already blacklisted",
		})
	}
	supplier.Blacklisted = true
	if err := tx.Save(&supplier).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to blacklist supplier",
		})
	}
	var inventory []models.Inventory
	if err := tx.Where("added_by = ?", supplier.Id).Find(&inventory).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find supplier's inventory",
		})
	}
	if len(inventory) == 0 {
		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to blacklist supplier",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Supplier blacklisted successfully",
		})
	}

	if err := tx.Model(&models.Inventory{}).
		Where("added_by = ?", supplier.Id).
		Update("blacklisted", true).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to blacklist inventory",
		})
	}

	var prods []models.Product
	if err := tx.Where("added_by = ?", supplier.Id).Find(&prods).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find supplier's products",
		})
	}

	if len(prods) == 0 {
		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to blacklist supplier",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Supplier blacklisted successfully",
		})
	}

	if err := tx.Model(&models.Product{}).
		Where("added_by = ?", supplier.Id).
		Update("blacklisted", true).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to blacklist inventory",
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to blacklist supplier",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Supplier blacklisted successfully",
		"supplier": supplier,
	})

}
