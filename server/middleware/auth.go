// middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"server/utils"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            cookie, err := r.Cookie("access_token")
            if err != nil {
                utils.WriteError(w, http.StatusUnauthorized, "Access token required")
                return
            }
            token = cookie.Value
        }

        // Remove "Bearer " prefix if present
        token = strings.TrimPrefix(token, "Bearer ")

        claims, err := utils.ValidateToken(token)
        if err != nil || claims.TokenType != "access" {
            utils.WriteError(w, http.StatusUnauthorized, "Invalid access token")
            return
        }

        // Add user info to context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "email", claims.Email)
        ctx = context.WithValue(ctx, "session_id", claims.SessionID)

        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
