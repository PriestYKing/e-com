package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"server/models"
	"server/utils"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
    utils.WriteJSON(w, http.StatusOK, map[string]string{
        "message": "Hello World - PostgreSQL Connected!",
        "status":  "healthy",
    })
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }
    
    if req.Name == "" {
        utils.WriteError(w, http.StatusBadRequest, "Name is required")
        return
    }
    
    user, err := models.CreateUser(req.Name,req.Email,req.Password)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
        return
    }
    
    utils.WriteJSON(w, http.StatusCreated, user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }
    
    user, err := models.GetUserByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.WriteError(w, http.StatusNotFound, "User not found")
            return
        }
        utils.WriteError(w, http.StatusInternalServerError, "Database error: "+err.Error())
        return
    }
    
    utils.WriteJSON(w, http.StatusOK, user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := models.GetAllUsers()
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "Database error: "+err.Error())
        return
    }
    
    utils.WriteJSON(w, http.StatusOK, users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        utils.WriteError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }
    
    err = models.DeleteUser(id)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.WriteError(w, http.StatusNotFound, "User not found")
            return
        }
        utils.WriteError(w, http.StatusInternalServerError, "Database error: "+err.Error())
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}
