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
    mux.HandleFunc("POST /register", handlers.RegisterUser)
    mux.HandleFunc("POST /login", handlers.LoginUser)


    //Product routes
    mux.HandleFunc("GET /products", handlers.GetAllProducts)
    
    // Apply CORS middleware and return the handler
    return middleware.EnableCORS(mux)
}
