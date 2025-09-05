package main

import (
	"gin-database-connect/initializers"
	"gin-database-connect/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
}

func main() {
	var adminCount int64
	initializers.DB.Model(&models.User{}).Where("role = ?", models.AdminRole).Count(&adminCount)

	if adminCount > 0 {
		log.Println("Admin user already exists")
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
		log.Printf("Warning: Using default admin password. Set ADMIN_PASSWORD in .env for security!")
	}

	adminName := os.Getenv("ADMIN_NAME")
	if adminName == "" {
		adminName = "Administrator"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := models.User{
		Name:     adminName,
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     models.AdminRole,
		IsActive: true,
	}

	if err := initializers.DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("Admin user created successfully!")
	log.Printf("Name: %s", admin.Name)
	log.Printf("Email: %s", admin.Email)
	log.Printf("Password: %s", adminPassword)
}
