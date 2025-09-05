package controllers

import (
	"gin-database-connect/initializers"
	"gin-database-connect/models"
	"gin-database-connect/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string  `json:"name" binding:"required,min=3"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Birthday *string `json:"birthday"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Claims struct {
	UserID uint            `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func generateRefreshToken(user *models.User) (string, error) {
	claims := &RefreshClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtSecret == "" {
		jwtSecret = "refresh_secret_hfdnjdjdkkk12&*^$acahihihaha"
	}

	return token.SignedString([]byte(jwtSecret))
}

func validateRefreshTokenJWT(tokenString string) (*RefreshClaims, error) {
	jwtSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtSecret == "" {
		jwtSecret = "refresh_secret_hfdnjdjdkkk12&*^$acahihihaha"
	}

	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

func validateRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims, err := validateRefreshTokenJWT(tokenString)
	if err != nil {
		return nil, err
	}

	var refreshToken models.RefreshToken
	result := initializers.DB.Where("token = ? AND is_revoked = false AND expires_at > ?",
		tokenString, time.Now()).First(&refreshToken)

	if result.Error != nil {
		return nil, err
	}

	return claims, nil
}

func saveRefreshToken(userID uint, tokenString string) error {
	initializers.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{})

	refreshToken := models.RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	return initializers.DB.Create(&refreshToken).Error
}

func revokeRefreshToken(tokenString string) error {
	return initializers.DB.Model(&models.RefreshToken{}).
		Where("token = ?", tokenString).
		Update("is_revoked", true).Error
}

func revokeAllUserTokens(userID uint) error {
	return initializers.DB.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func generateAccessToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "hfdnjdjdkkk12&*^$acahihihaha"
	}

	return token.SignedString([]byte(jwtSecret))
}

func createRefreshToken(user *models.User) (string, error) {
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		return "", err
	}

	if err := saveRefreshToken(user.ID, refreshToken); err != nil {
		return "", err
	}

	return refreshToken, nil
}

func Register(c *gin.Context) {
	var body RegisterRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := initializers.DB.Where("email = ?", body.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
		Role:     models.RegularUserRole,
		Birthday: utils.ParseBirthday(body.Birthday),
		Phone:    body.Phone,
	}

	if body.Role != nil {
		if currentUser, exists := c.Get("user"); exists {
			if currentUser.(*models.User).Role == models.AdminRole {
				if *body.Role == "admin" {
					user.Role = models.AdminRole
				}
			}
		}
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    user,
	})
}

func Login(c *gin.Context) {
	var body LoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
		return
	}

	accessToken, err := generateAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := createRefreshToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"tokens": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    300,
		},
		"user": user,
	})
}

func GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Refresh(c *gin.Context) {
	var body RefreshRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := validateRefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
		return
	}

	if err := revokeRefreshToken(body.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old token"})
		return
	}

	accessToken, err := generateAccessToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := createRefreshToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"tokens": TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			ExpiresIn:    300,
		},
	})
}

func Logout(c *gin.Context) {
	var body RefreshRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := validateRefreshTokenJWT(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	if err := revokeRefreshToken(body.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func LogoutAllDevices(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userObj := user.(*models.User)

	if err := revokeAllUserTokens(userObj.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout from all devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out from all devices successfully",
	})
}
