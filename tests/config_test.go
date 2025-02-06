package tests

import (
	"os"
	"testing"

	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/models"

	"github.com/joho/godotenv"
)

// Testconfig.Getenv teste la récupération des variables d'environnement
func testconfigGetenv(t *testing.T) {
	// Charger les variables d'environnement depuis un fichier .env (si présent)
	_ = godotenv.Load(".env.test") // Utiliser un fichier spécifique pour les tests

	testKey := "TEST_KEY"
	testValue := "test_value"
	defaultValue := "default_value"

	// Cas 1 : La variable d'environnement est définie dans le fichier .env.test
	os.Setenv(testKey, testValue) // On force une valeur pour le test
	got := config.GetEnv(testKey, defaultValue)
	if got != testValue {
		t.Errorf("config.Getenv(%q, %q) = %q; veut %q", testKey, defaultValue, got, testValue)
	}
	os.Unsetenv(testKey) // Nettoyer après test

	// Cas 2 : La variable d'environnement n'est pas définie, on attend la valeur par défaut
	got = config.GetEnv(testKey, defaultValue)
	if got != defaultValue {
		t.Errorf("config.Getenv(%q, %q) = %q; veut %q", testKey, defaultValue, got, defaultValue)
	}
}

// testInitDB teste que la base de donnée s'initialise correctement
func testInitDB(t *testing.T) {
	var DB = config.DB
	// Initialise la BDD
	config.InitDB(true)

	// Vérifier si la table "users" existe
	if !DB.Migrator().HasTable(&models.User{}) {
		t.Fatal("Erreur : la table 'users' n'a pas été créée")
	}

	// Vérifier si la table "tasks" existe
	if !DB.Migrator().HasTable(&models.Task{}) {
		t.Fatal("Erreur : la table 'tasks' n'a pas été créée")
	}
}
