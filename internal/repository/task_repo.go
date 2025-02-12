package repository

import (
	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/models"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetTasksByUser(userID uint) ([]models.Task, error)
	GetTaskByID(taskID uint) (*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(task *models.Task) error
}

// Implemetation par défaut de l'interface TaskRepository
type taskRepository struct{}

// Retourne une instance de TaskRepository
func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

// Créer une nouvelle tâche
func (t *taskRepository) CreateTask(task *models.Task) error {
	return config.DB.Create(task).Error
}

// Retourne toutes les tâches d'un utilisateur
func (t *taskRepository) GetTasksByUser(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := config.DB.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

// Retourne une tâche par son ID
func (t *taskRepository) GetTaskByID(taskID uint) (*models.Task, error) {
	var task models.Task
	err := config.DB.Where("id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Met à jour une tâche
func (t *taskRepository) UpdateTask(task *models.Task) error {
	return config.DB.Save(task).Error
}

// Supprime une tâche
func (t *taskRepository) DeleteTask(task *models.Task) error {
	return config.DB.Delete(task).Error
}
