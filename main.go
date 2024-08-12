package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/Hailemari/task_manager/data"
	"github.com/Hailemari/task_manager/router"

)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize the database connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	client, err := data.ConnectDB(mongoURI)  
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Initialize the user collection
	data.InitUserCollection(client)

	// Set up the router and start the server
	r := router.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
