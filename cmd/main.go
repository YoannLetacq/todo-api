// cmd/main.go
package main

import (
	"log"
	"os"

	"YoannLetacq/todo-api.git/config"
	"YoannLetacq/todo-api.git/internal/handlers"
	"YoannLetacq/todo-api.git/internal/repository"
	"YoannLetacq/todo-api.git/internal/services"
	"YoannLetacq/todo-api.git/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env (si présent)
	if err := godotenv.Load(); err != nil {
		log.Println("Avertissement: impossible de charger le fichier .env")
	}

	// Initialiser la base de données en mode production (false signifie non-test)
	config.InitDB(false)

	// Initialiser le repository et le service pour les utilisateurs
	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)
	handlers.InitUserHandlers(userService)

	// Initialiser le repository et le service pour les tâches
	taskRepo := repository.NewTaskRepository()
	taskService := services.NewTaskService(taskRepo)
	handlers.InitTaskHandlers(taskService)

	// Configurer le routeur avec l'ensemble des routes (utilisateurs et tâches)
	router := routes.SetupRouter()

	// Récupérer le port depuis la variable d'environnement ou utiliser 8080 par défaut
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Serveur démarré sur le port " + port)

	// Démarrer le serveur
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur:", err)
	}
}
