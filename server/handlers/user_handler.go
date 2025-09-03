package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"server/models"
	"server/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
    utils.WriteJSON(w, http.StatusOK, map[string]string{
        "message": "Hello World - PostgreSQL Connected!",
        "status":  "healthy",
    })
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    if req.Name == "" || req.Email == "" {
        utils.WriteError(w, http.StatusBadRequest, "Name and Email are required")
        return
    }
    
    existingUser,err := models.GetUserByEmail(req.Email)
    if err != nil && err != sql.ErrNoRows {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to check existing user: "+err.Error())
        return
    }

    if existingUser != nil {
        utils.WriteError(w, http.StatusConflict, "User already exists")
        return
    }

hashedPassword,err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)

if err != nil {
    utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password: "+err.Error())
    return
}

    user, err := models.CreateUser(req.Name, req.Email, hashedPassword)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
        return
    }
    
token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
    "email": user.Email,
    "exp": time.Now().Add(24 * time.Hour).Unix(),
})
tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
 if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to sign token: "+err.Error())
        return
    }

     http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    tokenString,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, 
        Expires:  time.Now().Add(24 * time.Hour),
        SameSite: http.SameSiteLaxMode,
    })

    utils.WriteJSON(w, http.StatusCreated, "User Created")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    user, err := models.GetUserByEmail(req.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.WriteError(w, http.StatusNotFound, "User not found")
            return
        }
        utils.WriteError(w, http.StatusInternalServerError, "Database error: "+err.Error())
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        utils.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
    "email": user.Email,
    "exp": time.Now().Add(24 * time.Hour).Unix(),
})
tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
 if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to sign token: "+err.Error())
        return
    }

     http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    tokenString,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, 
        Expires:  time.Now().Add(24 * time.Hour),
        SameSite: http.SameSiteLaxMode,
    })
    
    user.Password = ""

    utils.WriteJSON(w, http.StatusOK, user)
}

