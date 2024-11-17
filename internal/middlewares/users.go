package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	config "github.com/mohdjishin/order-inventory-management/config"
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"go.uber.org/zap"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("{\"code\": %s, \"message\": \"%s\"}", e.Code, e.Message)
}

var (
	ErrInvalidToken                  = Error{"INVALID_TOKEN", "Invalid or expired token"}
	ErrInvalidAuthHeader             = Error{"INVALID_AUTH_HEADER", "Invalid Authorization header format"}
	ErrUnauthorizationHeaderNotFound = Error{"UNAUTHORIZATION_HEADER_NOT_FOUND", "Authorization header not found"}
)

const authorization = "Authorization"

type CtxRoleKey struct{}
type CtxUserIDKey struct{}
type CtxEmailKey struct{}

func AuthMiddleware(c fiber.Ctx) error {
	log.Debug("AuthMiddleware")

	authHeader := c.Get(authorization)
	if authHeader == "" {
		log.Error("Authorization header not found")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrUnauthorizationHeaderNotFound)
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
		log.Error("Invalid Authorization header format", zap.Any("parts", parts))
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidAuthHeader)
	}

	tokenString := strings.TrimSpace(parts[1])

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().JwtKey), nil
	})

	if err != nil || !token.Valid {
		log.Error("Failed to parse token", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		fmt.Printf("id %T\n", (*claims)["id"])
		fmt.Println("Failed to extract user ID from token", claims)
		log.Error("Failed to extract user ID from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		fmt.Printf("email%T\n", (*claims)["email"])
		log.Error("Failed to extract email from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}
	role, ok := (*claims)["role"].(string)
	if !ok {
		fmt.Printf("role %T\n", (*claims)["role"])
		log.Error("Failed to extract role from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}

	c.Locals(CtxUserIDKey{}, userID)
	c.Locals(CtxEmailKey{}, email)
	c.Locals(CtxRoleKey{}, role)
	return c.Next()
}

func OnlySuppliers(c fiber.Ctx) error {
	role, ok := c.Locals(CtxRoleKey{}).(string)
	if !ok {
		log.Error("Failed to extract role from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to extract role from context",
		})
	}

	if role != models.SupplierRole.String() {
		log.Error("User is not a supplier")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is not a supplier",
		})
	}

	userID, ok := c.Locals(CtxUserIDKey{}).(float64)
	if !ok {
		log.Error("Failed to extract user ID from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to extract user ID from context",
		})
	}

	var user models.User
	if err := db.GetDb().First(&user, "id = ? AND role = ? AND approved = ?", userID, models.SupplierRole, true).Error; err != nil {
		log.Error("Supplier not found or not approved", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Supplier not found or not approved",
		})
	}

	if user.BlackListed {
		log.Error("Supplier is blacklisted")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Supplier is blacklisted",
		})
	}

	return c.Next()
}

func OnlyCustomer(c fiber.Ctx) error {
	role, ok := c.Locals(CtxRoleKey{}).(string)
	if !ok {
		log.Error("Failed to extract role from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to extract role from context",
		})
	}
	if role != models.CustomerRole.String() {
		log.Error("User is not a customer")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is not a customer",
		})
	}
	return c.Next()
}
