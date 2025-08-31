package routes

import (
	"net/http"

	"server/handlers"
	"server/middleware"
)

func SetupRoutes() http.Handler {
    mux := http.NewServeMux()
    
    // Root route
    mux.HandleFunc("/", handlers.HandleRoot)
    
    // User routes
    mux.HandleFunc("POST /users", handlers.CreateUser)
    mux.HandleFunc("GET /users/{id}", handlers.GetUser)
    mux.HandleFunc("GET /users", handlers.GetAllUsers)
    mux.HandleFunc("DELETE /users/{id}", handlers.DeleteUser)
    
    // Apply CORS middleware and return the handler
    return middleware.EnableCORS(mux)
}
