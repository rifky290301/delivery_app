package main

import (
	"delivery_app/config"
	"delivery_app/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize the database connection
	config.InitDB()

	// Initialize the routes
	r := routes.InitRoutes()

	// Start the server
	log.Println("Server starting on port 8000...")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
