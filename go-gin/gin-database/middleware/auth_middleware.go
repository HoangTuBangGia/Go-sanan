package middleware

import (
	"gin-database-connect/controllers"
	"gin-database-connect/initializers"
	"gin-database-connect/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		token, err := jwt.ParseWithClaims(tokenString, &controllers.Claims{}, func(token *jwt.Token) (interface{}, error) {
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				jwtSecret = "ffjfjfjfkskll123@#@#!"
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*controllers.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		var user models.User
		if err := initializers.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
			c.Abort()
			return
		}

		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userObj := user.(*models.User)
		if userObj.Role != models.AdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func UserOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userObj := user.(*models.User)

		if userObj.Role == models.AdminRole {
			c.Next()
			return
		}

		requestedUserID := c.Param("id")
		if requestedUserID != "" {
			currentUserIDStr := strconv.FormatUint(uint64(userObj.ID), 10)
			if requestedUserID != currentUserIDStr {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
