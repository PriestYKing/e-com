package routes

import (
	"net/http"

	"server/handlers"
	"server/middleware"
)

func SetupRoutes() http.Handler {
    mux := http.NewServeMux()
    
   
 
    // User routes
    mux.HandleFunc("POST /register", handlers.RegisterUser)
    mux.HandleFunc("POST /login", handlers.LoginUser)
    mux.HandleFunc("POST /logout", handlers.LogoutUser)

    //Product routes
    mux.HandleFunc("GET /products", handlers.GetAllProducts)
    
    // Apply CORS middleware and return the handler
    return middleware.EnableCORS(mux)
}
