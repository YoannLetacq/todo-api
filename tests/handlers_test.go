package tests

import (
	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/handlers"
	"YoannLetacq/todo-api.git/internal/models"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.InitDB(true)
	if config.DB == nil {
		log.Fatal("Erreur: Connexion à la BDD nulle.")
	}
	config.DB.Exec("DELETE from users")

	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Erreur: Echec des migrations.", err)
	}

	log.Println("Succès: Base de données initialisé avec succès.")
}

func TestRegisterHandler(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	setupTestDB()

	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	userData := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expectedResponse := `{"message":"Utilisateur enregistré avec succès !"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
