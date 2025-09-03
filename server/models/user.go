package models

import (
	"time"

	"server/config"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"created_at"`
}


func CreateUser(name string, email string, password []byte) (*User, error) {
    var user User
    err := config.DB.QueryRow(
        "INSERT INTO users (name,email,password) VALUES ($1,$2,$3) RETURNING id, name, created_at,email",
        name,
        email,
        password,
    ).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.Email)

    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
    var user User
    err := config.DB.QueryRow(
        "SELECT id, name, created_at,password FROM users WHERE email = $1",
        email,
    ).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.Password)

    if err != nil {
        return nil, err
    }
    
    return &user, nil
}


