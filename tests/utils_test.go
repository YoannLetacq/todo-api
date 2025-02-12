package tests

import (
	"YoannLetacq/todo-api.git/internal/utils"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "my_secret_key")

	userID := "1" // user_id est stocké en string
	email := "test@example.com"

	token, err := utils.GenerateJWT(userID, email)
	assert.NoError(t, err, "Erreur lors de la génération du token")
	assert.NotEmpty(t, token, "Le token généré est vide")

	parsedToken, claims, err := utils.ParseToken(token)
	assert.NoError(t, err, "Erreur lors du parsing du token")
	assert.True(t, parsedToken.Valid, "Le token généré n'est pas valide")

	// Vérifier les claims
	assert.Equal(t, userID, claims["user_id"], "L'ID utilisateur du token est invalide")
	assert.Equal(t, email, claims["email"], "L'email du token est invalide")
}

func TestParseToken(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Setenv("JWT_SECRET", "my_secret_key")

	userID := "1"
	email := "test@example.com"

	tokenString, err := utils.GenerateJWT(userID, email)
	assert.NoError(t, err, "Erreur lors de la génération du token.")

	parsedToken, claims, err := utils.ParseToken(tokenString)
	assert.NoError(t, err, "Erreur lors du parsing du token JWT")
	assert.True(t, parsedToken.Valid, "Le token JWT n'est pas valide")

	// Vérifier les claims
	assert.Equal(t, userID, claims["user_id"], "L'ID utilisateur du token est invalide")
	assert.Equal(t, email, claims["email"], "L'email du token est invalide")

	// Test avec un token invalide
	invalidToken := "invalid.token.string"
	_, _, err = utils.ParseToken(invalidToken)
	assert.Error(t, err, "Le parsing d'un token invalide aurait dû échouer")

	// Test avec un token signé avec une autre clé secrète
	os.Setenv("JWT_SECRET", "wrong_secret_key")
	_, _, err = utils.ParseToken(tokenString)
	assert.Error(t, err, "Le parsing aurait dû échouer avec une clé invalide")
}
