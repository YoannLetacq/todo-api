package models

import "github.com/jinzhu/gorm"

// Tâche de l'Utilisateur
type Task struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	Status      string `gorm:"default:'todo'" json:"status"`
	UserID      uint   `json:"user_id"`
}
