// tests/handler_test.go
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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// setHandlerTestDB initialise la BDD pour les tests de handlers.
func setHandlerTestDB() {
	config.InitDB(true)
	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM tasks")
	config.DB.AutoMigrate(&models.User{}, &models.Task{})
}

// initHandlerTestServices initialise et injecte les services dans les handlers.
func initHandlerTestServices() {
	// Service User
	userRepo := repository.NewUserRepository()
	userSvc := services.NewUserService(userRepo)
	handlers.InitUserHanlers(userSvc)

	// Service Task
	taskRepo := repository.NewTaskRepository()
	taskSvc := services.NewTaskService(taskRepo)
	handlers.InitTaskHandlers(taskSvc)
}

// createTestUser crée un utilisateur en BDD.
func createTestUser(t *testing.T) models.User {
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
	return user
}

// generateTokenForUser génère un token JWT pour un utilisateur.
func generateTokenForUser(user models.User, t *testing.T) string {
	token, err := utils.GenerateJWT(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		t.Fatal("Erreur lors de la génération du token JWT:", err)
	}
	return token
}

// TestRegisterUserHandler teste directement le handler RegisterUser.
func TestRegisterUserHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	registerData := map[string]string{
		"username": "handlerUser",
		"email":    "handlerUser@example.com",
		"password": "handlerPassword",
	}
	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Créer un contexte de test Gin
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.RegisterUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	assert.Equal(t, "Utilisateur enregistré avec succès !", resp["message"])
}

// TestLoginUserHandler teste directement le handler LoginHandler.
func TestLoginUserHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	// Créer un utilisateur
	user := createTestUser(t)

	// Préparer la requête de login
	loginData := map[string]string{
		"email":    user.Email,
		"password": "password",
	}
	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.LoginHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	token, exists := resp["token"]
	assert.True(t, exists, "Le token JWT est absent")
	parsedToken, claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	assert.Equal(t, strconv.Itoa(int(user.ID)), claims["user_id"])
}

// TestCreateTaskHandler teste directement le handler CreateTask.
func TestCreateTaskHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	user := createTestUser(t)
	token := generateTokenForUser(user, t)

	taskData := map[string]string{
		"title":       "Handler Task",
		"description": "Task created via handler",
	}
	jsonData, _ := json.Marshal(taskData)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.CreateTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	assert.Equal(t, "Task crée !", resp["message"])
	task, ok := resp["task"].(map[string]interface{})
	assert.True(t, ok, "La tâche doit être un objet JSON")
	assert.Equal(t, "Handler Task", task["title"])
	assert.Equal(t, "Task created via handler", task["description"])
	assert.Equal(t, "todo", task["status"])
}

// TestGetTasksHandler teste directement le handler GetTasks.
func TestGetTasksHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	user, token := createTestUserAndToken(t)
	// Créer deux tâches
	task1 := models.Task{Title: "Task 1", Description: "Desc 1", Status: "todo", UserID: user.ID}
	task2 := models.Task{Title: "Task 2", Description: "Desc 2", Status: "done", UserID: user.ID}
	config.DB.Create(&task1)
	config.DB.Create(&task2)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.GetTasks(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	tasks, ok := resp["tasks"].([]interface{})
	assert.True(t, ok, "Les tâches doivent être un tableau JSON")
	assert.Equal(t, 2, len(tasks))
}

// TestGetTaskHandler teste directement le handler GetTask.
func TestGetTaskHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	user, token := createTestUserAndToken(t)
	task := models.Task{Title: "Task Get", Description: "Desc Get", Status: "done", UserID: user.ID}
	config.DB.Create(&task)

	req, _ := http.NewRequest("GET", "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.GetTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	taskResp, ok := resp["task"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Task Get", taskResp["title"])
	assert.Equal(t, "Desc Get", taskResp["description"])
	assert.Equal(t, "done", taskResp["status"])
}

// TestUpdateTaskHandler teste directement le handler UpdateTask.
func TestUpdateTaskHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	user, token := createTestUserAndToken(t)
	task := models.Task{Title: "Task Update", Description: "Old Desc", Status: "todo", UserID: user.ID}
	config.DB.Create(&task)

	updatedData := map[string]string{
		"title":       "Task Updated",
		"description": "New Desc",
		"status":      "in progress",
	}
	jsonData, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(int(task.ID)), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.UpdateTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	assert.Equal(t, "Task mise a jour !", resp["message"])
	updatedTask, ok := resp["task"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Task Updated", updatedTask["title"])
	assert.Equal(t, "New Desc", updatedTask["description"])
	assert.Equal(t, "in progress", updatedTask["status"])
}

// TestDeleteTaskHandler teste directement le handler DeleteTask.
func TestDeleteTaskHandler(t *testing.T) {
	setHandlerTestDB()
	initHandlerTestServices()

	user, token := createTestUserAndToken(t)
	task := models.Task{Title: "Task Delete", Description: "Desc Delete", Status: "done", UserID: user.ID}
	config.DB.Create(&task)

	req, _ := http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.DeleteTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Erreur de parsing:", err)
	}
	assert.Equal(t, "Task supprimée.", resp["message"])

	// Vérifier que la tâche n'existe plus en base
	var deletedTask models.Task
	err := config.DB.First(&deletedTask, task.ID).Error
	assert.Error(t, err)
}
