// models/user.go
package models

import (
	"context"
	"database/sql"
	"server/cache"
	"server/config"
	"time"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"password,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

type UserSession struct {
    ID           int       `json:"id"`
    UserID       int       `json:"user_id"`
    IPAddress    string    `json:"ip_address"`
    Device       string    `json:"device"`
    DeviceID     string    `json:"device_id"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    CreatedAt    time.Time `json:"created_at"`
    IsActive     bool      `json:"is_active"`
}

type BlacklistedToken struct {
    ID        int       `json:"id"`
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}

// Create user with device session
func CreateUser(name, email string, password []byte, ipAddress, device, deviceID string) (*User, *UserSession, error) {
    tx, err := config.DB.Begin()
    if err != nil {
        return nil, nil, err
    }
    defer tx.Rollback()

    var user User
    err = tx.QueryRow(
        "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, created_at",
        name, email, password,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    
    if err != nil {
        return nil, nil, err
    }

    // Create initial session
    session, err := CreateUserSession(tx, user.ID, ipAddress, device, deviceID)
    if err != nil {
        return nil, nil, err
    }

    if err = tx.Commit(); err != nil {
        return nil, nil, err
    }

    return &user, session, nil
}

func GetUserByEmail(email string) (*User, error) {
    ctx := context.Background()
    
    // Try cache first
   // cacheKey := fmt.Sprintf("email:%s", email)
    var cachedUser User
    if found, err := cache.GetCachedUser(ctx, 0, &cachedUser); err == nil && found {
        return &cachedUser, nil
    }

    // Cache miss - query database
    var user User
    err := config.DB.QueryRow(
        "SELECT id, name, email, password, created_at FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    user.Email = email
    
    // Cache the user (async)
    go func() {
        cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        cache.CacheUser(cacheCtx, user.ID, user)
    }()
    
    return &user, nil
} 

func GetUserByID(userID int) (*User, error) {
    ctx := context.Background()
    
    // Try cache first
    var cachedUser User
    if found, err := cache.GetCachedUser(ctx, userID, &cachedUser); err == nil && found {
        return &cachedUser, nil
    }

    // Cache miss - query database
    var user User
    err := config.DB.QueryRow(
        "SELECT id, name, email, created_at FROM users WHERE id = $1",
        userID,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    // Cache the user (async)
    go func() {
        cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        cache.CacheUser(cacheCtx, user.ID, user)
    }()
    
    return &user, nil
}

func CreateUserSession(tx *sql.Tx, userID int, ipAddress, device, deviceID string) (*UserSession, error) {
    var session UserSession
    err := tx.QueryRow(`
        INSERT INTO user_sessions (user_id, ip_address, device, device_id, expires_at, is_active)
        VALUES ($1, $2, $3, $4, $5, true)
        RETURNING id, user_id, ip_address, device, device_id, expires_at, created_at, is_active`,
        userID, ipAddress, device, deviceID, time.Now().Add(7*24*time.Hour), // 7 days for refresh token
    ).Scan(&session.ID, &session.UserID, &session.IPAddress, &session.Device, 
           &session.DeviceID, &session.ExpiresAt, &session.CreatedAt, &session.IsActive)
    
    return &session, err
}

func GetActiveSession(userID int, ipAddress, device, deviceID string) (*UserSession, error) {
    var session UserSession
    err := config.DB.QueryRow(`
        SELECT id, user_id, ip_address, device, device_id, refresh_token, expires_at, created_at, is_active
        FROM user_sessions 
        WHERE user_id = $1 AND ip_address = $2 AND device = $3 AND device_id = $4 AND is_active = true
        ORDER BY created_at DESC LIMIT 1`,
        userID, ipAddress, device, deviceID,
    ).Scan(&session.ID, &session.UserID, &session.IPAddress, &session.Device,
           &session.DeviceID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.IsActive)
    
    return &session, err
}

func UpdateSessionRefreshToken(sessionID int, refreshToken string) error {
    _, err := config.DB.Exec(
        "UPDATE user_sessions SET refresh_token = $1 WHERE id = $2",
        refreshToken, sessionID,
    )
    return err
}

func GetSessionByRefreshToken(refreshToken string) (*UserSession, error) {
    var session UserSession
    err := config.DB.QueryRow(`
        SELECT id, user_id, ip_address, device, device_id, refresh_token, expires_at, created_at, is_active
        FROM user_sessions 
        WHERE refresh_token = $1 AND is_active = true AND expires_at > NOW()`,
        refreshToken,
    ).Scan(&session.ID, &session.UserID, &session.IPAddress, &session.Device,
           &session.DeviceID, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.IsActive)
    
    return &session, err
}

func DeactivateSession(sessionID int) error {
    _, err := config.DB.Exec(
        "UPDATE user_sessions SET is_active = false WHERE id = $1",
        sessionID,
    )
    return err
}

func BlacklistToken(token string, expiresAt time.Time) error {
    _, err := config.DB.Exec(
        "INSERT INTO blacklisted_tokens (token, expires_at) VALUES ($1, $2)",
        token, expiresAt,
    )
    return err
}

func IsTokenBlacklisted(token string) (bool, error) {
    var count int
    err := config.DB.QueryRow(
        "SELECT COUNT(*) FROM blacklisted_tokens WHERE token = $1 AND expires_at > NOW()",
        token,
    ).Scan(&count)
    
    return count > 0, err
}



func DeactivateAllUserSessions(userID int) error {
    _, err := config.DB.Exec(
        "UPDATE user_sessions SET is_active = false WHERE user_id = $1",
        userID,
    )
    return err
}

func GetAllActiveUserSessions(userID int) ([]UserSession, error) {
    rows, err := config.DB.Query(`
        SELECT id, user_id, ip_address, device, device_id, refresh_token, expires_at, created_at, is_active
        FROM user_sessions 
        WHERE user_id = $1 AND is_active = true`,
        userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var sessions []UserSession
    for rows.Next() {
        var session UserSession
        var refreshToken sql.NullString
        
        err := rows.Scan(&session.ID, &session.UserID, &session.IPAddress, &session.Device,
                        &session.DeviceID, &refreshToken, &session.ExpiresAt, 
                        &session.CreatedAt, &session.IsActive)
        if err != nil {
            return nil, err
        }
        
        if refreshToken.Valid {
            session.RefreshToken = refreshToken.String
        }
        
        sessions = append(sessions, session)
    }
    
    return sessions, rows.Err()
}

