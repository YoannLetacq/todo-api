// tests/router_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/handlers"
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/repository"
	"YoannLetacq/todo-api.git/internal/services"
	"YoannLetacq/todo-api.git/internal/utils"
	"YoannLetacq/todo-api.git/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// setRouterTestDB initialise la base de données pour les tests.
func setRouterTestDB() {
	config.InitDB(true)
	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM tasks")
	config.DB.AutoMigrate(&models.User{}, &models.Task{})
}

// initRouterTest initialise les services (User et Task), les injecte dans les handlers,
// et retourne l'instance du router.
func initRouterTest() *gin.Engine {
	// Service User
	userRepo := repository.NewUserRepository()
	userSvc := services.NewUserService(userRepo)
	handlers.InitUserHanlers(userSvc)

	// Service Task
	taskRepo := repository.NewTaskRepository()
	taskSvc := services.NewTaskService(taskRepo)
	handlers.InitTaskHandlers(taskSvc)

	return routes.SetupRouter()
}

// createTestUserAndToken crée un utilisateur dans la base de test et retourne l'utilisateur ainsi qu'un token JWT valide.
func createTestUserAndToken(t *testing.T) (models.User, string) {
	os.Setenv("JWT_SECRET", "my_secret_key")
	password := "password"
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal("Erreur lors du hachage du mot de passe:", err)
	}

	user := models.User{
		Username: "testUser",
		Email:    "testUser@example.com",
		Password: string(hashedPass),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		t.Fatal("Erreur lors de la création de l'utilisateur:", err)
	}

	token, err := utils.GenerateJWT(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		t.Fatal("Erreur lors de la génération du token JWT:", err)
	}
	return user, token
}

// TestRouterRegisterAndLogin teste les endpoints d'inscription et de connexion.
func TestRouterRegisterAndLogin(t *testing.T) {
	os.Setenv("JWT_SECRET", "my_secret_key")
	setRouterTestDB()
	router := initRouterTest()

	// --- Inscription ---
	registerData := map[string]string{
		"username": "routerUser",
		"email":    "routerUser@example.com",
		"password": "routerPassword",
	}
	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// --- Connexion ---
	loginData := map[string]string{
		"email":    "routerUser@example.com",
		"password": "routerPassword",
	}
	jsonData, _ = json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatal("Erreur de parsing de la réponse de login:", err)
	}
	_, exists := loginResp["token"]
	assert.True(t, exists, "Le token JWT est absent de la réponse")
}

// TestRouterTasksEndpoints teste les endpoints liés aux tâches (CRUD).
func TestRouterTasksEndpoints(t *testing.T) {
	os.Setenv("JWT_SECRET", "my_secret_key")
	setRouterTestDB()
	router := initRouterTest()

	// Création d'un utilisateur et génération de token.
	_, token := createTestUserAndToken(t)

	// --- Création de tâche ---
	taskData := map[string]string{
		"title":       "Router Task",
		"description": "Task created via router",
	}
	jsonData, _ := json.Marshal(taskData)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &createResp); err != nil {
		t.Fatal("Erreur de parsing de la réponse de création de tâche:", err)
	}
	assert.Equal(t, "Task crée !", createResp["message"])

	// --- Récupération de toutes les tâches ---
	req, _ = http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var getResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &getResp); err != nil {
		t.Fatal("Erreur de parsing de la réponse GET /tasks:", err)
	}
	tasks, ok := getResp["tasks"].([]interface{})
	assert.True(t, ok, "Les tâches doivent être un tableau JSON")
	assert.GreaterOrEqual(t, len(tasks), 1, "Il doit y avoir au moins 1 tâche")

	// --- Récupération d'une tâche par son ID ---
	createdTask := createResp["task"].(map[string]interface{})
	taskID := strconv.Itoa(int(createdTask["ID"].(float64)))
	req, _ = http.NewRequest("GET", "/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var getOneResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &getOneResp); err != nil {
		t.Fatal("Erreur de parsing de la réponse GET /tasks/:id:", err)
	}
	taskResp, ok := getOneResp["task"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Router Task", taskResp["title"])

	// --- Mise à jour de la tâche ---
	updatedData := map[string]string{
		"title":       "Router Task Updated",
		"description": "Task updated via router",
		"status":      "in progress",
	}
	jsonData, _ = json.Marshal(updatedData)
	req, _ = http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var updateResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &updateResp); err != nil {
		t.Fatal("Erreur de parsing de la réponse PUT /tasks/:id:", err)
	}
	assert.Equal(t, "Task mise a jour !", updateResp["message"])
	updatedTask := updateResp["task"].(map[string]interface{})
	assert.Equal(t, "Router Task Updated", updatedTask["title"])

	// --- Suppression de la tâche ---
	req, _ = http.NewRequest("DELETE", "/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var deleteResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &deleteResp); err != nil {
		t.Fatal("Erreur de parsing de la réponse DELETE /tasks/:id:", err)
	}
	assert.Equal(t, "Task supprimée.", deleteResp["message"])
}
