package routes

import (
	"YoannLetacq/todo-api.git/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter ... Configure les routes
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginHandler)

	taskGroup := router.Group("/tasks")
	{
		taskGroup.POST("", handlers.CreateTask)
		taskGroup.GET("", handlers.GetTasks)
		taskGroup.GET("/:id", handlers.GetTask)
		taskGroup.PUT("/:id", handlers.UpdateTask)
		taskGroup.DELETE("/:id", handlers.DeleteTask)
	}

	return router
}
