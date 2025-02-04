package config

import (
	"fmt"
	"log"

	"YoannLetacq/todo-api.git/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// connexion à la base de donnée
var DB *gorm.DB

// initialise la connexion à la DBB
func InitDB(testing bool) {
	var err error

	if testing {
		DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	}

	dbType := GetEnv("DB_TYPE", "DB_TYPE")

	if dbType == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			GetEnv("DB_HOST", "DB_HOST"),
			GetEnv("DB_USER", "DB_USER"),
			GetEnv("DB_PASSWORD", "DB_PASSWORD"),
			GetEnv("DB_NAME", "DB_NAME"),
			GetEnv("DB_PORT", "DB_PORT"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Base de connecté avec succès !")

	// Applicaiton des migrations
	DB.AutoMigrate(&models.User{}, &models.Task{})
}
