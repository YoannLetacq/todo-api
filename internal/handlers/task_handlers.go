package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/services"
	"YoannLetacq/todo-api.git/internal/utils"

	"github.com/gin-gonic/gin"
)

var taskservices services.TaskService

// InitTaskHandlers permet d'injecter le service dans les handlers
func InitTaskHandlers(s services.TaskService) {
	taskservices = s
}

// extractuserID extrait le user_id du token JWT des headers.
func ExtractUserID(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("Authorization Token manquant")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Token mal formé")
	}

	tokenString := parts[1]
	_, claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"]
	if !ok || userID == "" {
		return "", errors.New("Token invalide: user_id manquant")
	}
	return userID, nil
}

// CreateTask crée un handler pour la création de tâches.
func CreateTask(c *gin.Context) {
	userID, err := ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorise." + err.Error()})
		return
	}

	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides."})
		return
	}

	if task.Status == "" {
		task.Status = "todo"
	}

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id invalide"})
		return
	}

	task.UserID = uint(uid)

	if err := taskservices.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echec de la création de la Task"})

		log.Println("Erreur lors de la creation de la tache:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task crée !", "task": task})
}

// GetTasks recupere toutes les tâches pour un utilisateur
func GetTasks(c *gin.Context) {
	userID, err := ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorise." + err.Error()})
		return
	}

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID invalide"})
		return
	}

	tasks, err := taskservices.GetTasksByUser(uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echech de la recuperation des tâches."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTask recupere une tâche pour un utilisateur GET /tasks/:id
func GetTask(c *gin.Context) {
	userID, err := ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorise."})
		return
	}

	uid, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID invalide"})
		return
	}

	taskID := c.Param("id")
	tid, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tâche invalide"})
		return
	}

	task, err := taskservices.GetTaskByID(uint(tid))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche introuvable."})
		return
	}

	if task.UserID != uint(uid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette tâche ne vous appartiens pas."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// UpdateTask met a jour une tâche, dont son status PUT /taks/update/:id
func UpdateTask(c *gin.Context) {
	userID, err := ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorise."})
		return
	}

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID invalide"})
		return
	}

	taskID := c.Param("id")
	tid, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tâche invalide"})
		return
	}

	task, err := taskservices.GetTaskByID(uint(tid))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche introuvable."})
		return
	}

	if task.UserID != uint(uid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette tâche ne vous appartiens pas."})
		return
	}

	var updateData models.Task
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides."})
		return
	}

	task.Title = updateData.Title
	task.Description = updateData.Description
	if updateData.Status != "" {
		task.Status = updateData.Status
	}

	if err := taskservices.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echec de la mise a jour de la Task", "detail": err.Error()})
		log.Println("Erreur lors de la mise à jour de la tâche:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task mise a jour !", "task": task})
}

// DeleteTask supprime une tâche DELETE /tasks/delete/:id
func DeleteTask(c *gin.Context) {
	userID, err := ExtractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorise."})
		return
	}

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID invalide"})
		return
	}

	taskID := c.Param("id")
	tid, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tâche invalide"})
		return
	}

	task, err := taskservices.GetTaskByID(uint(tid))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche introuvable."})
		return
	}

	if task.UserID != uint(uid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette tâche ne vous appartiens pas."})
		return
	}

	if err := taskservices.DeleteTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echec de la suppression de la Task", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task supprimée."})
}
