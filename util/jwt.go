package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohdjishin/order-inventory-management/config"
	"github.com/mohdjishin/order-inventory-management/internal/models"
)

// var jwtSecretKey = []byte("your-jwt-secret-key")

func GenerateToken(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"role":  user.Role.String(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Get().JwtKey))
}
