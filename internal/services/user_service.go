package services

import (
	"YoannLetacq/todo-api.git/internal/models"
	"YoannLetacq/todo-api.git/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(user *models.User) error
	LoginUser(email, password string) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService cree une nouvelle instance de UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// RegisterUser permet d'enregistrer un utilisateur
func (s *userService) RegisterUser(user *models.User) error {
	return s.repo.CreateUser(user)
}

// LoginUser permet de connecter un utilisateur
func (s *userService) LoginUser(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
