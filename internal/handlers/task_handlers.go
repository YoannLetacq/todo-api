package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/utils"

	"github.com/gin-gonic/gin"
)

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

	if err := config.DB.Create(&task).Error; err != nil {
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
	var tasks []models.Task
	if err := config.DB.Where("user_id= ?", uint(uid)).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echec pour recuperer la Task."})
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
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche introuvable."})
		return
	}

	if task.UserID != uint(uid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette tâche ne vous apartiens pas."})
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
	var task models.Task

	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche introuvable."})
		return
	}

	if task.UserID != uint(uid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette tâche ne vous apartiens pas."})
		return
	}

	var updateTask models.Task
	if err := c.ShouldBindJSON(&updateTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides."})
		return
	}

	task.Title = updateTask.Title
	task.Description = updateTask.Description
	if updateTask.Status != "" {
		task.Status = updateTask.Status
	}

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echec de la mise a jour de la Task."})

		log.Println("Erreur lors de la mise a jour de la tache:", err)
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
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tâche introuvable."})
		return
	}

	if task.UserID != uint(uid) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette tâche ne vous apartiens pas."})
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Echec de la suppression de la Task."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task supprimée."})
}
