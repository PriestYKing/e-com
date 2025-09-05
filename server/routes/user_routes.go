package routes

import (
	"net/http"
	"time"

	"server/handlers"
	"server/middleware"
)

func SetupRoutes() http.Handler {
    mux := http.NewServeMux()

    // Auth routes - NO authentication required, only rate limiting
    mux.HandleFunc("/register", methodGuard("POST", 
        applyMiddleware(handlers.RegisterUser, 
            middleware.AuthRateLimitMiddleware(),
        ),
    ))
    
    mux.HandleFunc("/login", methodGuard("POST", 
        applyMiddleware(handlers.LoginUser, 
            middleware.AuthRateLimitMiddleware(),
        ),
    ))
    
    // Logout - requires authentication (user must be logged in to logout)
    mux.HandleFunc("/logout", methodGuard("POST", 
        applyMiddleware(handlers.LogoutUser, 
            middleware.AuthMiddleware,
            middleware.AuthRateLimitMiddleware(),
        ),
    ))

    // Public API routes - can be cached, no auth required
    mux.HandleFunc("/products", methodGuard("GET", 
        applyMiddleware(handlers.GetAllProducts, 
            middleware.APIRateLimitMiddleware(),
            middleware.APICacheMiddleware(10*time.Minute),
        ),
    ))

    mux.HandleFunc("/me", methodGuard("GET", 
        applyMiddleware(handlers.GetCurrentUser, 
            middleware.AuthMiddleware,
            middleware.AuthRateLimitMiddleware(),
        ),
    ))

    // Apply CORS middleware and return the handler
    return middleware.EnableCORS(mux);
}

// Helper function to enforce HTTP methods
func methodGuard(method string, handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != method {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        handler.ServeHTTP(w, r)
    }
}

// Apply middleware in reverse order so the first one wraps the innermost
func applyMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
