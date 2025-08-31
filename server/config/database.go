package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    host := getEnv("DB_HOST","")
    port := getEnv("DB_PORT", "")
    user := getEnv("DB_USER", "")
    password := getEnv("DB_PASSWORD", "")
    dbname := getEnv("DB_NAME", "")
    
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    var err error
    DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Failed to open database connection:", err)
    }
    
    if err = DB.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    fmt.Println("Successfully connected to PostgreSQL database")
}

func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}


func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
