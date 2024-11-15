package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	"github.com/mohdjishin/order-inventory-management/logger"
	"github.com/mohdjishin/order-inventory-management/util"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Login(c fiber.Ctx) error {
	var req LoginRequest

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		logger.Error().Err(err).Msg("Error decoding request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad request, unable to parse input",
		})
	}

	logger.Debug().Any("request", req).Msg("Login request received")

	if validationErrors, err := util.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"fields": validationErrors,
		})
	}

	var user models.User
	if err := db.GetDb().Where("email = ?", req.Email).First(&user).Error; err != nil {
		logger.Error().Err(err).Msg("Error fetching user")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Error().Err(err).Msg("Error comparing password")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
		})
	}

	token, err := util.GenerateToken(user)
	if err != nil {
		logger.Error().Err(err).Msg("Error generating token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate token",
		})
	}

	// Return the token in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"token": token,
			"user": fiber.Map{
				"email": user.Email,
				"role":  user.Role.String(),
			},
		},
	})
}
