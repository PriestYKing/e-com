package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
    RedisClient *redis.Client
    RedisCacheClient *redis.Client
    ctx = context.Background()
)

type RedisConfig struct {
    Host         string
    Port         string
    Password     string
    DBCache      int
    DBRateLimit  int
    MaxRetries   int
    PoolSize     int
    MinIdleConns int
}

func InitRedis() error {
    config := RedisConfig{
        Host:         getEnv("REDIS_HOST", "localhost"),
        Port:         getEnv("REDIS_PORT", "6379"),
        Password:     getEnv("REDIS_PASSWORD", "abc"),
        DBCache:      0, // Database 0 for caching
        DBRateLimit:  1, // Database 1 for rate limiting
        MaxRetries:   3,
        PoolSize:     10,
        MinIdleConns: 5,
    }

    // Redis client for rate limiting
    RedisClient = redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
        Password:     config.Password,
        DB:           config.DBRateLimit,
        MaxRetries:   config.MaxRetries,
        PoolSize:     config.PoolSize,
        MinIdleConns: config.MinIdleConns,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolTimeout:  4 * time.Second,
       
    })

    // Redis client for caching
    RedisCacheClient = redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
        Password:     config.Password,
        DB:           config.DBCache,
        MaxRetries:   config.MaxRetries,
        PoolSize:     config.PoolSize,
        MinIdleConns: config.MinIdleConns,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolTimeout:  4 * time.Second,
       
    })

    // Test connections
    if err := testRedisConnection(RedisClient, "Rate Limiting"); err != nil {
        return err
    }

    if err := testRedisConnection(RedisCacheClient, "Caching"); err != nil {
        return err
    }

    log.Println("Redis connections established successfully")
    return nil
}

func testRedisConnection(client *redis.Client, purpose string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return fmt.Errorf("failed to connect to Redis for %s: %w", purpose, err)
    }

    log.Printf("Redis connection for %s: OK", purpose)
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
