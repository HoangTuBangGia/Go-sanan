package initializers

import (
	"fmt"
	"gin-database-connect/models"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := strings.TrimSpace(os.Getenv("DB_PASS"))
	name := strings.TrimSpace(os.Getenv("DB_NAME"))

	if host == "" || port == "" || user == "" || name == "" {
		log.Fatalf("missing db env (host=%q port=%q user=%q name=%q)", host, port, user, name)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("cannot connect database: %v", err)
	}

	DB = db
}

func AutoMigrate() error {
	return DB.AutoMigrate(&models.User{})
}
