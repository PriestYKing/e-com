package models

import (
	"database/sql"
	"time"

	"server/config"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
    Name string `json:"name"`
}

func CreateUser(name string) (*User, error) {
    var user User
    err := config.DB.QueryRow(
        "INSERT INTO users (name) VALUES ($1) RETURNING id, name, created_at",
        name,
    ).Scan(&user.ID, &user.Name, &user.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func GetUserByID(id int) (*User, error) {
    var user User
    err := config.DB.QueryRow(
        "SELECT id, name, created_at FROM users WHERE id = $1",
        id,
    ).Scan(&user.ID, &user.Name, &user.CreatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func GetAllUsers() ([]User, error) {
    rows, err := config.DB.Query("SELECT id, name, created_at FROM users ORDER BY id")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}

func DeleteUser(id int) error {
    result, err := config.DB.Exec("DELETE FROM users WHERE id = $1", id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    
    return nil
}
