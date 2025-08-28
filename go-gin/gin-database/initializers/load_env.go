package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load()
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME"}

	for _, k := range required {
		if os.Getenv(k) == "" {
			log.Printf("warning: missing env %s", k)
		}
	}
}
