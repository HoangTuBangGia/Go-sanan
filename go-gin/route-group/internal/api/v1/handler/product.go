package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (u *ProductHandler) GetProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "List all products (v1)"})
}

func (u *ProductHandler) GetProductsByIdV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Get product by ID (v1)"})
}

func (u *ProductHandler) PostProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "Create product (v1)"})
}

func (u *ProductHandler) PutProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Update product (v1)"})
}

func (u *ProductHandler) DeleteProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Delete product (v1)"})
}
