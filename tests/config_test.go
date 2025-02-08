package tests

import (
	"os"
	"testing"

	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/models"

	"github.com/joho/godotenv"
)

func TestConfigGetEnv(t *testing.T) {
	_ = godotenv.Load(".env.test")

	testKey := "TEST_KEY"
	testValue := "test_value"
	defaultValue := "default_value"

	os.Setenv(testKey, testValue)
	got := config.GetEnv(testKey, defaultValue)
	if got != testValue {
		t.Errorf("config.GetEnv(%q, %q) = %q; want %q", testKey, defaultValue, got, testValue)
	}
	os.Unsetenv(testKey)

	got = config.GetEnv(testKey, defaultValue)
	if got != defaultValue {
		t.Errorf("config.GetEnv(%q, %q) = %q; want %q", testKey, defaultValue, got, defaultValue)
	}
}

func TestInitDB(t *testing.T) {
	config.InitDB(true)

	if !config.DB.Migrator().HasTable(&models.User{}) {
		t.Fatal("Erreur : la table 'users' n'a pas été créée")
	}

	if !config.DB.Migrator().HasTable(&models.Task{}) {
		t.Fatal("Erreur : la table 'tasks' n'a pas été créée")
	}
}
