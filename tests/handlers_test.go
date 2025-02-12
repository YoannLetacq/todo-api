package tests

import (
	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/handlers"
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.InitDB(true)
	if config.DB == nil {
		log.Fatal("Erreur: Connexion à la BDD nulle.")
	}
	// Nettoyer les tables
	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM tasks")
	if err := config.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("Erreur : Echec des migrations.", err)
	}

	log.Println("Succès: Base de données initialisée avec succès.")
}

func TestRegisterHandler(t *testing.T) {
	setupTestDB()

	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	userData := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "securepassword",
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Utilisateur enregistré avec succès !", response["message"])
}

func TestLoginHandler(t *testing.T) {
	setupTestDB()
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Créer un utilisateur avec un mot de passe hashé
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := models.User{
		Username: "testUser",
		Email:    "test@example.com",
		Password: string(hashedPass),
	}
	config.DB.Create(&user)

	router := gin.Default()
	router.POST("/login", handlers.LoginHandler)

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	token, exists := response["token"]
	assert.True(t, exists, "Le token JWT est absent de la réponse")

	parsedToken, claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	assert.NotEmpty(t, claims["user_id"], "Le user_id du token ne doit pas être vide")
	assert.NotEmpty(t, claims["email"], "L'email du token ne doit pas être vide")

	assert.Equal(t, strconv.Itoa(int(user.ID)), claims["user_id"])
	assert.Equal(t, user.Email, claims["email"])
}

func createTestUserAndToken(t *testing.T) (models.User, string) {
	os.Setenv("JWT_SECRET", "test_secret_key")
	password := "password"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Username: "testUser",
		Email:    "testUser@example.com",
		Password: string(hashedPass),
	}

	config.DB.Create(&user)

	token, err := utils.GenerateJWT(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		t.Fatal("Erreur: Impossible de générer le token JWT")
	}
	return user, token
}
