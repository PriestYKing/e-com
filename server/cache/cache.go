package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"server/config"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
    client *redis.Client
}

type CacheConfig struct {
    TTL               time.Duration
    KeyPrefix         string
    EnableCompression bool
    MaxSize           int64
}

type CacheItem struct {
    Data      interface{} `json:"data"`
    ExpiresAt time.Time   `json:"expires_at"`
    Version   int         `json:"version"`
    Tags      []string    `json:"tags,omitempty"`
}

type CacheStats struct {
    Hits       int64   `json:"hits"`
    Misses     int64   `json:"misses"`
    HitRate    float64 `json:"hit_rate"`
    Size       int64   `json:"size"`
    Memory     string  `json:"memory_usage"`
}


// Initialize DefaultCache after Redis is initialized
func InitCache() error {
    if config.RedisCacheClient == nil {
        return fmt.Errorf("Redis cache client not initialized")
    }
    
    DefaultCache = &Cache{
        client: config.RedisCacheClient,
    }
    
    log.Println("Cache initialized successfully")
    return nil
}
var (
    DefaultCache = &Cache{client: config.RedisCacheClient}
    
    // Different cache configurations for different data types
    UserCacheConfig = CacheConfig{
        TTL:       30 * time.Minute,
        KeyPrefix: "user",
        MaxSize:   1000,
    }

    SessionCacheConfig = CacheConfig{
        TTL:       24 * time.Hour,
        KeyPrefix: "session",
        MaxSize:   5000,
    }

    APIResponseCacheConfig = CacheConfig{
        TTL:       5 * time.Minute,
        KeyPrefix: "api_response",
        MaxSize:   10000,
    }

    DatabaseQueryCacheConfig = CacheConfig{
        TTL:       15 * time.Minute,
        KeyPrefix: "db_query",
        MaxSize:   2000,
    }
)

func NewCache() *Cache {
    return &Cache{
        client: config.RedisCacheClient,
    }
}

// Generic cache operations
func (c *Cache) Set(ctx context.Context, key string, value interface{}, config CacheConfig) error {
    fullKey := fmt.Sprintf("%s:%s", config.KeyPrefix, key)
    
    item := CacheItem{
        Data:      value,
        ExpiresAt: time.Now().Add(config.TTL),
        Version:   1,
    }
    
    data, err := json.Marshal(item)
    if err != nil {
        return fmt.Errorf("failed to marshal cache item: %w", err)
    }

    pipe := c.client.Pipeline()
    pipe.Set(ctx, fullKey, data, config.TTL)
    
    // Update cache statistics
    pipe.Incr(ctx, "cache:stats:sets")
    pipe.HIncrBy(ctx, "cache:stats:size", config.KeyPrefix, 1)
    
    _, err = pipe.Exec(ctx)
    if err != nil {
        log.Printf("Cache set error for key %s: %v", fullKey, err)
        return err
    }

    return nil
}

func (c *Cache) Get(ctx context.Context, key string, config CacheConfig, dest interface{}) (bool, error) {
    fullKey := fmt.Sprintf("%s:%s", config.KeyPrefix, key)
    
    data, err := c.client.Get(ctx, fullKey).Result()
    if err != nil {
        if err == redis.Nil {
            // Cache miss
            c.client.Incr(ctx, "cache:stats:misses")
            return false, nil
        }
        return false, fmt.Errorf("cache get error: %w", err)
    }

    var item CacheItem
    if err := json.Unmarshal([]byte(data), &item); err != nil {
        return false, fmt.Errorf("failed to unmarshal cache item: %w", err)
    }

    // Check if expired (additional safety check)
    if time.Now().After(item.ExpiresAt) {
        c.Delete(ctx, key, config)
        c.client.Incr(ctx, "cache:stats:misses")
        return false, nil
    }

    // Unmarshal the actual data
    itemData, err := json.Marshal(item.Data)
    if err != nil {
        return false, fmt.Errorf("failed to marshal item data: %w", err)
    }

    if err := json.Unmarshal(itemData, dest); err != nil {
        return false, fmt.Errorf("failed to unmarshal to destination: %w", err)
    }

    // Cache hit
    c.client.Incr(ctx, "cache:stats:hits")
    return true, nil
}

func (c *Cache) Delete(ctx context.Context, key string, config CacheConfig) error {
    fullKey := fmt.Sprintf("%s:%s", config.KeyPrefix, key)
    
    pipe := c.client.Pipeline()
    pipe.Del(ctx, fullKey)
    pipe.HIncrBy(ctx, "cache:stats:size", config.KeyPrefix, -1)
    
    _, err := pipe.Exec(ctx)
    return err
}

func (c *Cache) InvalidateByPattern(ctx context.Context, pattern string) error {
    keys, err := c.client.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }

    if len(keys) > 0 {
        return c.client.Del(ctx, keys...).Err()
    }
    return nil
}

func (c *Cache) InvalidateByTag(ctx context.Context, tag string) error {
    tagKey := fmt.Sprintf("tag:%s", tag)
    keys, err := c.client.SMembers(ctx, tagKey).Result()
    if err != nil {
        return err
    }

    if len(keys) > 0 {
        pipe := c.client.Pipeline()
        pipe.Del(ctx, keys...)
        pipe.Del(ctx, tagKey)
        _, err = pipe.Exec(ctx)
    }
    return err
}

func (c *Cache) GetStats(ctx context.Context) (*CacheStats, error) {
    pipe := c.client.Pipeline()
    hitsCmd := pipe.Get(ctx, "cache:stats:hits")
    missesCmd := pipe.Get(ctx, "cache:stats:misses")
    sizeCmd := pipe.HLen(ctx, "cache:stats:size")
    infoCmd := pipe.Info(ctx, "memory")
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return nil, err
    }

    hits, _ := hitsCmd.Int64()
    misses, _ := missesCmd.Int64()
    size, _ := sizeCmd.Result()
    memInfo, _ := infoCmd.Result()

    var hitRate float64
    if hits+misses > 0 {
        hitRate = float64(hits) / float64(hits+misses)
    }

    // Extract memory usage from INFO command
    memoryUsage := extractMemoryUsage(memInfo)

    return &CacheStats{
        Hits:    hits,
        Misses:  misses,
        HitRate: hitRate,
        Size:    size,
        Memory:  memoryUsage,
    }, nil
}

// High-level cache functions for specific use cases
func CacheUser(ctx context.Context, userID int, user interface{}) error {
    key := fmt.Sprintf("id:%d", userID)
    return DefaultCache.Set(ctx, key, user, UserCacheConfig)
}

func GetCachedUser(ctx context.Context, userID int, dest interface{}) (bool, error) {
    key := fmt.Sprintf("id:%d", userID)
    return DefaultCache.Get(ctx, key, UserCacheConfig, dest)
}

func CacheSession(ctx context.Context, sessionID int, session interface{}) error {
    key := fmt.Sprintf("id:%d", sessionID)
    return DefaultCache.Set(ctx, key, session, SessionCacheConfig)
}

func GetCachedSession(ctx context.Context, sessionID int, dest interface{}) (bool, error) {
    key := fmt.Sprintf("id:%d", sessionID)
    return DefaultCache.Get(ctx, key, SessionCacheConfig, dest)
}

func CacheAPIResponse(ctx context.Context, endpoint, params string, response interface{}) error {
    // Create a hash of endpoint + parameters as key
    h := sha256.New()
    h.Write([]byte(endpoint + params))
    key := hex.EncodeToString(h.Sum(nil))
    return DefaultCache.Set(ctx, key, response, APIResponseCacheConfig)
}

func GetCachedAPIResponse(ctx context.Context, endpoint, params string, dest interface{}) (bool, error) {
    h := sha256.New()
    h.Write([]byte(endpoint + params))
    key := hex.EncodeToString(h.Sum(nil))
    return DefaultCache.Get(ctx, key, APIResponseCacheConfig, dest)
}

// Cache-aside pattern for database queries
func CacheQuery(ctx context.Context, query string, args []interface{}, result interface{}) error {
    key := generateQueryKey(query, args)
    return DefaultCache.Set(ctx, key, result, DatabaseQueryCacheConfig)
}

func GetCachedQuery(ctx context.Context, query string, args []interface{}, dest interface{}) (bool, error) {
    key := generateQueryKey(query, args)
    return DefaultCache.Get(ctx, key, DatabaseQueryCacheConfig, dest)
}

// Helper functions
func generateQueryKey(query string, args []interface{}) string {
    h := sha256.New()
    h.Write([]byte(query))
    for _, arg := range args {
        h.Write([]byte(fmt.Sprintf("%v", arg)))
    }
    return hex.EncodeToString(h.Sum(nil))[:16]
}

func extractMemoryUsage(info string) string {
    lines := strings.Split(info, "\r\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "used_memory_human:") {
            return strings.TrimPrefix(line, "used_memory_human:")
        }
    }
    return "Unknown"
}
