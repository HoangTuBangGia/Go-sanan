package controllers

import (
	"gin-database-connect/initializers"
	"gin-database-connect/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=3"`
	Email    string  `json:"email" binding:"required,email"`
	Birthday *string `json:"birthday"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=3"`
	Birthday *string `json:"birthday"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

func parseBirthday(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", *s)

	if err != nil {
		return nil
	}

	return &t
}

func UserCreate(c *gin.Context) {
	var body CreateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Birthday: parseBirthday(body.Birthday),
		Phone:    body.Phone,
	}

	if body.IsActive != nil {
		user.IsActive = *body.IsActive
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func UserList(c *gin.Context) {
	var users []models.User

	if err := initializers.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func UserGet(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserUpdate(c *gin.Context) {
	id := c.Param("id")
	var body UpdateUserRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if body.Name != nil {
		user.Name = *body.Name
	}

	if body.Birthday != nil {
		user.Birthday = parseBirthday(body.Birthday)
	}

	if body.Phone != nil {
		user.Phone = body.Phone
	}

	if body.IsActive != nil {
		user.IsActive = *body.IsActive
	}

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")

	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := initializers.DB.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
