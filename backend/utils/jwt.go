package utils

import (
	"time"

	"companyInternalManagement/config"
	"companyInternalManagement/middlewares"
	"companyInternalManagement/models"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user *models.User) (string, error) {
	claims := middlewares.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "company-management",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecret)
}
