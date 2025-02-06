package tests

import (
	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {

	userID := "1"
	email := "test@example.com"

	token, err := utils.GenerateJWT(userID, email)

	assert.NoError(t, err, "Erreur lors de la generation du token.")
	assert.NotEmpty(t, token, "Le token genere est vide.")

	tokenString, err := jwt.ParseWithClaims(token, &jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetEnv("JWT_SECRET", "my_secret_key")), nil
		})
	assert.NoError(t, err, "Erreur lors du parsing do token.")
	assert.True(t, tokenString.Valid, "Le token genere est valide.")

	claims, ok := tokenString.Claims.(*jwt.MapClaims)

	assert.True(t, ok, "Impossible  de recuperer les claims du token.")
	assert.Equal(t, userID, (*claims)["user_id"],
		"L'ID utilisateur du tokent est invalide.",
	)
	assert.Equal(t, email, (*claims)["email"], "L'email de token est invalide.")

	expiration, ok := (*claims)["exp"].(float64)
	assert.True(t, ok, "Le claims d'expiration est incorrecte ou absent.")
	assert.Greater(t, expiration, float64(time.Now().Unix()),
		"Le token est expire.",
	)

}
