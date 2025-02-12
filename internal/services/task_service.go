package services

import (
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/repository"
)

type TaskService interface {
	CreateTask(task *models.Task) error
	GetTasksByUser(userID uint) ([]models.Task, error)
	GetTaskByID(taskID uint) (*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(task *models.Task) error
}

// retourne une instance de TaskService
type taskService struct {
	repo repository.TaskRepository
}

// Retourne une instance de TaskService
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

// Créer une nouvelle tâche
func (s *taskService) CreateTask(task *models.Task) error {
	return s.repo.CreateTask(task)
}

// Retourne toutes les tâches d'un utilisateur
func (s *taskService) GetTasksByUser(userID uint) ([]models.Task, error) {
	return s.repo.GetTasksByUser(userID)
}

// Retourne une tâche par son ID
func (s *taskService) GetTaskByID(taskID uint) (*models.Task, error) {
	return s.repo.GetTaskByID(taskID)
}

// Met à jour une tâche
func (s *taskService) UpdateTask(task *models.Task) error {
	return s.repo.UpdateTask(task)
}

// Supprime une tâche
func (s *taskService) DeleteTask(task *models.Task) error {
	return s.repo.DeleteTask(task)
}
