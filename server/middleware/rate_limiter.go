package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"server/config"
	"server/utils"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
    client *redis.Client
}

type RateLimitConfig struct {
    RequestsPerMinute int
    RequestsPerHour   int
    RequestsPerDay    int
    BurstSize         int
    WindowSize        time.Duration
    KeyPrefix         string
}

type RateLimitResult struct {
    Allowed           bool
    Limit             int
    Remaining         int
    ResetTime         time.Time
    RetryAfter        time.Duration
}

var (
    // Different rate limits for different endpoints
    DefaultRateLimit = RateLimitConfig{
        RequestsPerMinute: 60,
        RequestsPerHour:   1000,
        RequestsPerDay:    10000,
        BurstSize:         10,
        WindowSize:        time.Minute,
        KeyPrefix:         "rate_limit",
    }

    AuthRateLimit = RateLimitConfig{
        RequestsPerMinute: 10,
        RequestsPerHour:   100,
        RequestsPerDay:    500,
        BurstSize:         3,
        WindowSize:        time.Minute,
        KeyPrefix:         "auth_rate_limit",
    }

    // Sliding window rate limiter Lua script
    slidingWindowScript = `
        local key = KEYS[1]
        local window = tonumber(ARGV[1])
        local limit = tonumber(ARGV[2])
        local current_time = tonumber(ARGV[3])
        
        -- Remove expired entries
        redis.call('ZREMRANGEBYSCORE', key, 0, current_time - window)
        
        -- Count current requests in window
        local current_requests = redis.call('ZCARD', key)
        
        if current_requests < limit then
            -- Add current request
            redis.call('ZADD', key, current_time, current_time .. ':' .. math.random())
            redis.call('EXPIRE', key, window)
            return {1, limit - current_requests - 1, current_time + window}
        else
            local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
            local reset_time = current_time + window
            if #oldest > 0 then
                reset_time = tonumber(oldest[2]) + window
            end
            return {0, 0, reset_time}
        end
    `

    // Token bucket rate limiter Lua script
    tokenBucketScript = `
        local key = KEYS[1]
        local capacity = tonumber(ARGV[1])
        local tokens = tonumber(ARGV[2])
        local interval = tonumber(ARGV[3])
        local current_time = tonumber(ARGV[4])
        
        local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
        local tokens_count = tonumber(bucket[1]) or capacity
        local last_refill = tonumber(bucket[2]) or current_time
        
        -- Calculate tokens to add
        local time_passed = current_time - last_refill
        local tokens_to_add = math.floor(time_passed / interval) * tokens
        tokens_count = math.min(capacity, tokens_count + tokens_to_add)
        
        if tokens_count >= 1 then
            tokens_count = tokens_count - 1
            redis.call('HMSET', key, 'tokens', tokens_count, 'last_refill', current_time)
            redis.call('EXPIRE', key, interval * capacity)
            return {1, tokens_count, current_time + interval}
        else
            redis.call('HMSET', key, 'tokens', tokens_count, 'last_refill', current_time)
            redis.call('EXPIRE', key, interval * capacity)
            return {0, tokens_count, current_time + interval}
        end
    `
)

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        client: config.RedisClient,
    }
}

func (rl *RateLimiter) CheckRateLimit(ctx context.Context, identifier string, config RateLimitConfig) (*RateLimitResult, error) {
    // Use sliding window algorithm
    return rl.slidingWindowCheck(ctx, identifier, config)
}

func (rl *RateLimiter) slidingWindowCheck(ctx context.Context, identifier string, config RateLimitConfig) (*RateLimitResult, error) {
    key := fmt.Sprintf("%s:%s:%s", config.KeyPrefix, identifier, "sliding")
    currentTime := time.Now().Unix()
    windowSeconds := int64(config.WindowSize.Seconds())

    result, err := rl.client.Eval(ctx, slidingWindowScript, []string{key}, 
        windowSeconds, config.RequestsPerMinute, currentTime).Result()
    
    if err != nil {
        log.Printf("Redis rate limit error: %v", err)
        // Fail open - allow request if Redis is down
        return &RateLimitResult{
            Allowed:   true,
            Limit:     config.RequestsPerMinute,
            Remaining: config.RequestsPerMinute - 1,
            ResetTime: time.Now().Add(config.WindowSize),
        }, nil
    }

    resultSlice := result.([]interface{})
    allowed := resultSlice[0].(int64) == 1
    remaining := int(resultSlice[1].(int64))
    resetTime := time.Unix(resultSlice[2].(int64), 0)

    return &RateLimitResult{
        Allowed:    allowed,
        Limit:      config.RequestsPerMinute,
        Remaining:  remaining,
        ResetTime:  resetTime,
        RetryAfter: time.Until(resetTime),
    }, nil
}

func (rl *RateLimiter) tokenBucketCheck(ctx context.Context, identifier string, config RateLimitConfig) (*RateLimitResult, error) {
    key := fmt.Sprintf("%s:%s:%s", config.KeyPrefix, identifier, "bucket")
    currentTime := time.Now().Unix()
    intervalSeconds := int64(config.WindowSize.Seconds()) / int64(config.RequestsPerMinute)

    result, err := rl.client.Eval(ctx, tokenBucketScript, []string{key}, 
        config.BurstSize, 1, intervalSeconds, currentTime).Result()
    
    if err != nil {
        log.Printf("Redis rate limit error: %v", err)
        return &RateLimitResult{
            Allowed:   true,
            Limit:     config.BurstSize,
            Remaining: config.BurstSize - 1,
            ResetTime: time.Now().Add(time.Duration(intervalSeconds) * time.Second),
        }, nil
    }

    resultSlice := result.([]interface{})
    allowed := resultSlice[0].(int64) == 1
    remaining := int(resultSlice[1].(int64))
    resetTime := time.Unix(resultSlice[2].(int64), 0)

    return &RateLimitResult{
        Allowed:    allowed,
        Limit:      config.BurstSize,
        Remaining:  remaining,
        ResetTime:  resetTime,
        RetryAfter: time.Until(resetTime),
    }, nil
}

func getClientIdentifier(r *http.Request) string {
    // Priority: User ID > IP + User Agent hash > IP only
    userID := r.Context().Value("user_id")
    if userID != nil {
        return fmt.Sprintf("user:%v", userID)
    }

    ip := utils.GetClientIP(r)
	
    userAgent := r.Header.Get("User-Agent")
    
    if userAgent != "" {
        hash := sha256.Sum256([]byte(userAgent))
        return fmt.Sprintf("ip:%s:ua:%s", ip, hex.EncodeToString(hash[:8]))
    }

    return fmt.Sprintf("ip:%s", ip)
}

// Middleware functions
func RateLimitMiddleware(config RateLimitConfig) func(http.HandlerFunc) http.HandlerFunc {
    limiter := NewRateLimiter()
    
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
            defer cancel()

            identifier := getClientIdentifier(r)
            result, err := limiter.CheckRateLimit(ctx, identifier, config)
            
            if err != nil {
                log.Printf("Rate limiting error: %v", err)
                // Fail open - continue with request
                next.ServeHTTP(w, r)
                return
            }

            // Set rate limit headers
            w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
            w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
            w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))

            if !result.Allowed {
                w.Header().Set("Retry-After", strconv.Itoa(int(result.RetryAfter.Seconds())))
                utils.WriteError(w, http.StatusTooManyRequests, fmt.Sprintf(
                    "Rate limit exceeded. Try again in %d seconds", 
                    int(result.RetryAfter.Seconds())))
                return
            }

            next.ServeHTTP(w, r)
        }
    }
}

// Specific middleware for auth endpoints
func AuthRateLimitMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return RateLimitMiddleware(AuthRateLimit)
}

// Default rate limiting for API endpoints
func APIRateLimitMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return RateLimitMiddleware(DefaultRateLimit)
}
