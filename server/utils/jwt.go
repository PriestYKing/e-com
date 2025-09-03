// utils/jwt.go
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"server/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID    int    `json:"user_id"`
    Email     string `json:"email"`
    SessionID int    `json:"session_id"`
    TokenType string `json:"token_type"` // "access" or "refresh"
    jwt.RegisteredClaims
}

func GenerateTokenPair(userID int, email string, sessionID int) (string, string, error) {
    // Access token (15 minutes)
    accessClaims := &Claims{
        UserID:    userID,
        Email:     email,
        SessionID: sessionID,
        TokenType: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return "", "", err
    }

    // Refresh token (7 days)
    refreshClaims := &Claims{
        UserID:    userID,
        Email:     email,
        SessionID: sessionID,
        TokenType: "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
    // Check if token is blacklisted
    isBlacklisted, err := models.IsTokenBlacklisted(tokenString)
    if err != nil {
        return nil, err
    }
    if isBlacklisted {
        return nil, errors.New("token is blacklisted")
    }

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

func GenerateDeviceID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func GetClientInfo(r *http.Request) (string, string) {
    // Get IP address
    ip := r.Header.Get("X-Forwarded-For")
    if ip == "" {
        ip = r.Header.Get("X-Real-IP")
    }
    if ip == "" {
        ip = r.RemoteAddr
    }

    // Get User Agent as device info
    device := r.Header.Get("User-Agent")
    if device == "" {
        device = "Unknown Device"
    }

    return ip, device
}
