package main

import (
	"gin-database-connect/controllers"
	"gin-database-connect/initializers"
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
		users := v1.Group("/users")

		users.POST("", controllers.UserCreate)
		users.GET("", controllers.UserList)
		users.GET("/:id", controllers.UserGet)
		users.PUT("/:id", controllers.UserUpdate)
		users.DELETE("/:id", controllers.UserDelete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
