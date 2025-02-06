package utils

import (
	"time"

	"YoannLetacq/todo-api.git/config"

	"github.com/golang-jwt/jwt"
)

// Genere un token jwt valide 24h
func GenerateJWT(userID, email string) (string, error) {

	secretKey := config.GetEnv("JWT_SECRET", "my_secret_key")

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
