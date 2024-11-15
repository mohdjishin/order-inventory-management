package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	"github.com/mohdjishin/order-inventory-management/logger"
)

var jwtSecretKey = []byte("your-jwt-secret-key")

func GenerateToken(user models.User) (string, error) {

	claims := jwt.MapClaims{
		"userID": user.Id,
		"email":  user.Email,
		"role":   user.Role.String(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		logger.Error().Err(err).Msg("Error signing the token")
		return "", err
	}

	return signedToken, nil
}
