package utils

import (
	"errors"
	"log"
	"time"

	"YoannLetacq/todo-api.git/config"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID, email string) (string, error) {
	secretKey := config.GetEnv("JWT_SECRET", "test_secret_key")
	log.Println("🔑 Clé utilisée pour SIGNER :", secretKey)

	if secretKey == "" {
		return "", errors.New("clé JWT manquante")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", errors.New("échec de la signature du token JWT")
	}

	log.Println("✅ Token généré :", tokenString)
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, map[string]string, error) {
	secretKey := config.GetEnv("JWT_SECRET", "test_secret_key")
	log.Println("🔑 Clé utilisée pour VERIFICATION :", secretKey)

	if secretKey == "" {
		return nil, nil, errors.New("clé JWT manquante")
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		log.Println("❌ Échec de la validation du token JWT :", err)
		return nil, nil, errors.New("échec de la validation du token JWT")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("impossible de récupérer les claims du token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return nil, nil, errors.New("user_id invalide ou vide dans le token")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, nil, errors.New("email invalide ou vide dans le token")
	}

	log.Println("✅ Token valide avec claims :", claims)
	return token, map[string]string{"user_id": userID, "email": email}, nil
}
