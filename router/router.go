package router

import (
	"github.com/Hailemari/auth_task_manager/controllers"
	"github.com/Hailemari/auth_task_manager/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", controllers.CreateUser)
	r.POST("/login", controllers.LoginUser)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/tasks", controllers.GetTasks)
		auth.GET("/tasks/:id", controllers.GetTask)

		// Admin-only routes
		admin := auth.Group("/")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/tasks", controllers.AddTask)
			admin.PUT("/tasks/:id", controllers.UpdateTask)
			admin.DELETE("/tasks/:id", controllers.DeleteTask)
			admin.POST("/promote/:username", controllers.PromoteUser)
		}
	}

	return r
}
