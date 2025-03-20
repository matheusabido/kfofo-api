package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/matheusabido/kfofo-api/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env não encontrado. Continuando sem ele.")
	}

	db.SetupDB()
	SetupRoutes()
}
