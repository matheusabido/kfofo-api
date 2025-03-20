package main

import (
	"fmt"
	"os"

	"github.com/matheusabido/kfofo-api/db"
)

func main() {
	// if err := godotenv.Load(); err != nil && os.Getenv("OCI_PRIVATE_KEY_BASE64") == "" {
	// 	log.Fatal("Failed to load .env")
	// }
	fmt.Println(os.Getenv("DB_URL"))

	db.SetupDB()
	SetupRoutes()
}
