package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	"github.com/mohdjishin/order-inventory-management/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=100"`
	LastName  string `json:"lastName" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required,len=10,numeric"`
}

func createUser(c fiber.Ctx, req CreateUserRequest, role models.Role) error {
	if validationErrors, err := util.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Validation failed",
			"fields": validationErrors,
		})
	}

	var existingUser models.User
	if err := db.GetDb().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	if err := db.GetDb().Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Phone number already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		Role:      role,
	}

	if err := db.GetDb().Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":        user.Id,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role.String(),
		},
	})
}

func UpdateInventory(req *UpdateInventoryRequest, userID uint) error {
	tx := db.GetDb().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var inventory models.Inventory
	if err := tx.Where("product_id = ? AND added_by = ?", req.ProductID, userID).First(&inventory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("inventory not found for the given product and user")
		}
		tx.Rollback()
		return err
	}

	inventory.Stock += req.NewStock
	if req.NewPrice > 0 {
		inventory.BasePrice = req.NewPrice
	}

	if err := tx.Save(&inventory).Error; err != nil {
		tx.Rollback()
		return err
	}

	var product models.Product
	if err := tx.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("product not found: %w", err)
	}

	if req.NewPrice > 0 && product.Price != req.NewPrice {
		pricingHistory := models.PricingHistory{
			ProductID: product.ID,
			OldPrice:  product.Price,
			NewPrice:  req.NewPrice,
		}
		if err := tx.Create(&pricingHistory).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create pricing history: %w", err)
		}

		product.Price = req.NewPrice
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update product price: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
