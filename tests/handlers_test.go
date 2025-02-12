package tests

import (
	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/handlers"
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/repository"
	"YoannLetacq/todo-api.git/internal/services"
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

func initTaskServiceForTest() {
	repo := repository.NewTaskRepository()
	svc := services.NewTaskService(repo)
	handlers.InitTaskHandlers(svc)
}

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
	os.Setenv("JWT_SECRET", "my_secret_key")

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
	os.Setenv("JWT_SECRET", "my_secret_key")
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

func TestCreateTask(t *testing.T) {
	setupTestDB()
	initTaskServiceForTest()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, token := createTestUserAndToken(t)

	router := gin.Default()
	router.POST("/tasks", handlers.CreateTask)

	taskData := map[string]string{
		"title": "Test Task",
		// status: Non donné, doit être défini à "todo" par défaut
		"description": "Test Task Description",
	}

	jsonData, _ := json.Marshal(taskData)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task crée !", response["message"])
	log.Println(response)

	task, ok := response["task"].(map[string]interface{})
	assert.True(t, ok, "La tâche doit être un objet JSON")

	assert.Equal(t, "Test Task", task["title"])
	assert.Equal(t, "Test Task Description", task["description"])
	assert.Equal(t, "todo", task["status"]) // test le status par defaut
}

func TestGetTasks(t *testing.T) {
	setupTestDB()
	initTaskServiceForTest()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	user, token := createTestUserAndToken(t)

	task1 := models.Task{Title: "Task 1", Description: "Task 1 Description", Status: "todo", UserID: user.ID}
	task2 := models.Task{Title: "Task 2", Description: "Task 2 Description", Status: "done", UserID: user.ID}

	config.DB.Create(&task1)
	config.DB.Create(&task2)

	router := gin.Default()
	router.GET("/tasks", handlers.GetTasks)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	tasks, ok := response["tasks"].([]interface{})
	assert.True(t, ok, "Les tâches doivent être un tableau JSON")
	assert.Len(t, tasks, 2, "Il doit y avoir 2 tâches")
	log.Println(response)
}

// TestGetTask verifies that a single task can be retrieved.
func TestGetTask(t *testing.T) {
	setupTestDB()
	initTaskServiceForTest()
	user, token := createTestUserAndToken(t)

	// Create a task.
	task := models.Task{Title: "Task Get", Description: "Description Get", Status: "done", UserID: user.ID}
	config.DB.Create(&task)

	router := gin.Default()
	router.GET("/tasks/:id", handlers.GetTask)

	req, _ := http.NewRequest("GET", "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	taskResp, ok := response["task"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Task Get", taskResp["title"])
	assert.Equal(t, "Description Get", taskResp["description"])
	assert.Equal(t, "done", taskResp["status"])
}

// TestUpdateTask verifies that a task can be updated, including its status.
func TestUpdateTask(t *testing.T) {
	setupTestDB()
	initTaskServiceForTest()
	user, token := createTestUserAndToken(t)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Create a task.
	task := models.Task{Title: "Task Update", Description: "Old Description", Status: "todo", UserID: user.ID}
	config.DB.Create(&task)

	router := gin.Default()
	router.PUT("/tasks/:id", handlers.UpdateTask)

	updatedData := map[string]string{
		"title":       "Task Updated",
		"description": "New Description",
		"status":      "in progress",
	}
	jsonData, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(int(task.ID)), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task mise a jour !", response["message"])

	taskResp, ok := response["task"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Task Updated", taskResp["title"])
	assert.Equal(t, "New Description", taskResp["description"])
	assert.Equal(t, "in progress", taskResp["status"])
}

// TestDeleteTask verifies that a task can be deleted.
func TestDeleteTask(t *testing.T) {
	setupTestDB()
	initTaskServiceForTest()
	user, token := createTestUserAndToken(t)

	// Create a task.
	task := models.Task{Title: "Task Delete", Description: "Description Delete", Status: "done", UserID: user.ID}
	config.DB.Create(&task)

	router := gin.Default()
	router.DELETE("/tasks/:id", handlers.DeleteTask)

	req, _ := http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task supprimée.", response["message"])

	// Verify that the task was removed.
	var deletedTask models.Task
	err = config.DB.First(&deletedTask, task.ID).Error
	assert.Error(t, err) // Expected error because the task should be deleted.
}
