package models

import (
	"github.com/jinzhu/gorm"
)

// Utilisateur de l'application
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Task     []Task `gorm:"foreignKey:UserID"`
}
