package main

import (
	"gin-database-connect/controllers"
	"gin-database-connect/initializers"
	"gin-database-connect/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/refresh", controllers.Refresh)
			auth.POST("/logout", controllers.Logout)
		}

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.POST("/logout-all-devices", controllers.LogoutAllDevices)

			users := protected.Group("/users")
			{
				users.POST("", middleware.AdminMiddleware(), controllers.UserCreate)
				users.GET("", middleware.AdminMiddleware(), controllers.UserList)
				users.GET("/:id", middleware.UserOwnershipMiddleware(), controllers.UserGet)
				users.PUT("/:id", middleware.UserOwnershipMiddleware(), controllers.UserUpdate)
				users.DELETE("/:id", middleware.AdminMiddleware(), controllers.UserDelete)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
