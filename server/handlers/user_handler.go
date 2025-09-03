// handlers/user_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server/config"
	"server/models"
	"server/utils"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    DeviceID string `json:"device_id,omitempty"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    DeviceID string `json:"device_id,omitempty"`
}

type TokenResponse struct {
    AccessToken  string      `json:"access_token"`
    RefreshToken string      `json:"refresh_token"`
    User         models.User `json:"user"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    if req.Name == "" || req.Email == "" || req.Password == "" {
        utils.WriteError(w, http.StatusBadRequest, "Name, Email and Password are required")
        return
    }

    // Check if user already exists
    existingUser, err := models.GetUserByEmail(req.Email)
    if err != nil && err != sql.ErrNoRows {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to check existing user")
        return
    }
    if existingUser != nil {
        utils.WriteError(w, http.StatusConflict, "User already exists")
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
        return
    }

    // Get client info
    ipAddress, device := utils.GetClientInfo(r)
    deviceID := req.DeviceID
    if deviceID == "" {
        deviceID = utils.GenerateDeviceID()
    }

    // Create user with session
    user, session, err := models.CreateUser(req.Name, req.Email, hashedPassword, ipAddress, device, deviceID)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }

    // Generate token pair
    accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, user.Email, session.ID)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to generate tokens")
        return
    }

    // Update session with refresh token
    if err := models.UpdateSessionRefreshToken(session.ID, refreshToken); err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to update session")
        return
    }

    // Set cookies
    setTokenCookies(w, accessToken, refreshToken)

    user.Password = ""
    response := TokenResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         *user,
    }

    utils.WriteJSON(w, http.StatusCreated, response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    // Get user
    user, err := models.GetUserByEmail(req.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
            return
        }
        utils.WriteError(w, http.StatusInternalServerError, "Database error: "+err.Error())
        return
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        utils.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    // Get client info
    ipAddress, device := utils.GetClientInfo(r)
    deviceID := req.DeviceID
    if deviceID == "" {
        deviceID = utils.GenerateDeviceID()
    }

    // Always create a new session for login
    tx, err := config.DB.Begin()
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Database error: "+err.Error())
        return
    }
    defer tx.Rollback()

    // Create new session
    session, err := models.CreateUserSession(tx, user.ID, ipAddress, device, deviceID)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to create session: "+err.Error())
        return
    }
    
    if err = tx.Commit(); err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to commit transaction: "+err.Error())
        return
    }

    // Generate token pair
    accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, user.Email, session.ID)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to generate tokens: "+err.Error())
        return
    }

    // Update session with new refresh token
    if err := models.UpdateSessionRefreshToken(session.ID, refreshToken); err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to update session: "+err.Error())
        return
    }

    // Set cookies
    setTokenCookies(w, accessToken, refreshToken)

    user.Password = ""
    response := TokenResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         *user,
    }

    utils.WriteJSON(w, http.StatusOK, response)
}


func RefreshToken(w http.ResponseWriter, r *http.Request) {
    refreshToken := r.Header.Get("Authorization")
    if refreshToken == "" {
        cookie, err := r.Cookie("refresh_token")
        if err != nil {
            utils.WriteError(w, http.StatusUnauthorized, "Refresh token required")
            return
        }
        refreshToken = cookie.Value
    }

    // Validate refresh token
    claims, err := utils.ValidateToken(refreshToken)
    if err != nil || claims.TokenType != "refresh" {
        utils.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
        return
    }

    // Get session
    session, err := models.GetSessionByRefreshToken(refreshToken)
    if err != nil {
        utils.WriteError(w, http.StatusUnauthorized, "Invalid session")
        return
    }

    // Generate new token pair
    accessToken, newRefreshToken, err := utils.GenerateTokenPair(claims.UserID, claims.Email, session.ID)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to generate tokens")
        return
    }

    // Update session with new refresh token
    if err := models.UpdateSessionRefreshToken(session.ID, newRefreshToken); err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to update session")
        return
    }

    // Blacklist old refresh token
    models.BlacklistToken(refreshToken, claims.ExpiresAt.Time)

    // Set new cookies
    setTokenCookies(w, accessToken, newRefreshToken)

    utils.WriteJSON(w, http.StatusOK, map[string]string{
        "access_token":  accessToken,
        "refresh_token": newRefreshToken,
    })
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
    var accessToken, refreshToken string
    
    // Get access token from header
    authHeader := r.Header.Get("Authorization")
    if authHeader != "" {
        accessToken = strings.TrimPrefix(authHeader, "Bearer ")
    }
    
    // Get access token from cookie if not in header
    if accessToken == "" {
        if cookie, err := r.Cookie("access_token"); err == nil {
            accessToken = cookie.Value
        }
    }

    // Get refresh token from cookie
    if cookie, err := r.Cookie("refresh_token"); err == nil {
        refreshToken = cookie.Value
    }

    log.Printf("Logout attempt - Access token present: %v, Refresh token present: %v", 
               accessToken != "", refreshToken != "")

    // Validate and get claims from access token
    if accessToken != "" {
        claims, err := utils.ValidateToken(accessToken)
        if err != nil {
            log.Printf("Failed to validate access token: %v", err)
        } else if claims.TokenType != "access" {
            log.Printf("Token is not an access token: %s", claims.TokenType)
        } else {
            log.Printf("Valid access token found - User ID: %d, Session ID: %d", 
                      claims.UserID, claims.SessionID)
            
            // Deactivate session
            if err := models.DeactivateSession(claims.SessionID); err != nil {
                log.Printf("Failed to deactivate session %d: %v", claims.SessionID, err)
            } else {
                log.Printf("Successfully deactivated session %d", claims.SessionID)
            }
            
            // Blacklist access token
            if err := models.BlacklistToken(accessToken, claims.ExpiresAt.Time); err != nil {
                log.Printf("Failed to blacklist access token: %v", err)
            } else {
                log.Printf("Successfully blacklisted access token")
            }
            
            // Blacklist refresh token if present
            if refreshToken != "" {
                refreshClaims, err := utils.ValidateToken(refreshToken)
                if err != nil {
                    log.Printf("Failed to validate refresh token: %v", err)
                } else {
                    if err := models.BlacklistToken(refreshToken, refreshClaims.ExpiresAt.Time); err != nil {
                        log.Printf("Failed to blacklist refresh token: %v", err)
                    } else {
                        log.Printf("Successfully blacklisted refresh token")
                    }
                }
            }
        }
    } else {
        log.Printf("No access token found for logout")
    }

    // Clear cookies
    clearTokenCookies(w)
    log.Printf("Cleared cookies")

    utils.WriteJSON(w, http.StatusOK, map[string]string{
        "message": "Logged out successfully",
    })
}

func setTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Set to true in production
        Expires:  time.Now().Add(15 * time.Minute),
        SameSite: http.SameSiteLaxMode,
    })

    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Set to true in production
        Expires:  time.Now().Add(7 * 24 * time.Hour),
        SameSite: http.SameSiteLaxMode,
    })
}

func clearTokenCookies(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Expires:  time.Now().Add(-time.Hour),
    })

    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Expires:  time.Now().Add(-time.Hour),
    })
}
