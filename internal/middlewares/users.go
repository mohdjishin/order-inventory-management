package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	config "github.com/mohdjishin/order-inventory-management/config"
	log "github.com/mohdjishin/order-inventory-management/logger"
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

func AuthMiddleware(c fiber.Ctx) error {
	log.Debug().Str("request-url", c.OriginalURL()).Msg("Auth middleware")

	authHeader := c.Get(authorization)
	if authHeader == "" {
		log.Error().Msg("Authorization header not found")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrUnauthorizationHeaderNotFound)
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
		log.Error().Str("authHeader", authHeader).Msg("Invalid Authorization header format")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidAuthHeader)
	}

	tokenString := strings.TrimSpace(parts[1])

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().JwtKey), nil
	})

	if err != nil || !token.Valid {
		log.Error().Err(err).Msg("Failed to parse token")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}
	fmt.Println("claims", claims)
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		fmt.Printf("id %T\n", (*claims)["id"])
		fmt.Println("Failed to extract user ID from token", claims)
		log.Error().Msg("Failed to extract user ID from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}
	fmt.Println("userID", userID)

	email, ok := (*claims)["email"].(string)
	if !ok {
		fmt.Printf("email%T\n", (*claims)["email"])
		log.Error().Msg("Failed to extract email from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}
	role, ok := (*claims)["role"].(string)
	if !ok {
		fmt.Printf("role %T\n", (*claims)["role"])
		log.Error().Msg("Failed to extract role from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(ErrInvalidToken)
	}

	c.Locals("userId", userID)
	c.Locals("email", email)
	c.Locals("role", role)
	return c.Next()
}
