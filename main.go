package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/matheusabido/kfofo-api/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env")
	}

	db.SetupDB()
	SetupRoutes()
}
