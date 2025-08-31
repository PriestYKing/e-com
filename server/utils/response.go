package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
    }
}

func WriteError(w http.ResponseWriter, status int, message string) {
    errorResp := ErrorResponse{
        Error:   http.StatusText(status),
        Message: message,
    }
    WriteJSON(w, status, errorResp)
}
