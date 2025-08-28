package main

import (
	"gin-database-connect/initializers"
	"log"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDatabase()
	initializers.AutoMigrate()
}

func main() {
	if err := initializers.AutoMigrate(); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}
	log.Println("migration done")
}
