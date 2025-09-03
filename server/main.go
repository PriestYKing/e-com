package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"server/config"
	"server/routes"

	"github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    // Initialize database
    config.InitDB()
    defer config.CloseDB()

    if err := config.InitRedis(); err != nil {
        log.Fatal("Failed to initialize Redis:", err)
    }
    
    // Setup routes
    mux := routes.SetupRoutes()
    
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    
    fmt.Printf("Server listening on :%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}
