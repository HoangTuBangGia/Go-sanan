package main

import (
	v1handler "hung/route-group/internal/api/v1/handler"
	v2handler "hung/route-group/internal/api/v2/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		user := v1.Group("users")
		{
			userHandlerV1 := v1handler.NewUserHandler()

			user.GET("", userHandlerV1.GetUsersV1)
			user.GET(":id", userHandlerV1.GetUsersByIdV1)
			user.POST("", userHandlerV1.PostUsersV1)
			user.PUT(":id", userHandlerV1.PutUsersV1)
			user.DELETE(":id", userHandlerV1.DeleteUsersV1)
		}

		product := v1.Group("/products")
		{
			productHandlerV1 := v1handler.NewProductHandler()

			product.GET("", productHandlerV1.GetProductsV1)
			product.GET("/:id", productHandlerV1.GetProductsByIdV1)
			product.POST("", productHandlerV1.PostProductsV1)
			product.PUT("/:id", productHandlerV1.PutProductsV1)
			product.DELETE("/:id", productHandlerV1.DeleteProductsV1)
		}
	}

	v2 := r.Group("/api/v2")
	{
		user := v2.Group("/users")
		{
			userHandlerV2 := v2handler.NewUserHandler()

			user.GET("", userHandlerV2.GetUsersV2)
			user.GET(":id", userHandlerV2.GetUsersByIdV2)
			user.POST("", userHandlerV2.PostUsersV2)
			user.PUT(":id", userHandlerV2.PutUsersV2)
			user.DELETE(":id", userHandlerV2.DeleteUsersV2)
		}
	}

	r.Run(":8080")
}
