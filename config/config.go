package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Initialize .env file
func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file.")
	}
}

// Getenv return the value for a given key of the loaded .env.
// You can specify a defaultValue.
func GetEnv(secret, defaultValue string) string {
	if value, exists := os.LookupEnv(secret); exists {
		return value
	}
	return defaultValue
}
