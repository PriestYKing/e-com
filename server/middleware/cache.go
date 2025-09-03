package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"server/cache"
)

type ResponseCache struct {
    StatusCode int                 `json:"status_code"`
    Headers    map[string][]string `json:"headers"`
    Body       []byte              `json:"body"`
    CachedAt   time.Time           `json:"cached_at"`
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
    body       *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    rw.body.Write(b)
    return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
    rw.statusCode = statusCode
    rw.ResponseWriter.WriteHeader(statusCode)
}

func APICacheMiddleware(cacheDuration time.Duration) func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // Only cache GET requests
            if r.Method != http.MethodGet {
                next.ServeHTTP(w, r)
                return
            }

            // Skip caching for authenticated requests with user-specific data
            if r.Header.Get("Authorization") != "" {
                next.ServeHTTP(w, r)
                return
            }

            ctx := r.Context()
            cacheKey := generateCacheKey(r)

            // Try to get from cache
            var cachedResponse ResponseCache
            found, err := cache.DefaultCache.Get(ctx, cacheKey, cache.APIResponseCacheConfig, &cachedResponse)
            
            if err == nil && found {
                // Serve from cache
                w.Header().Set("X-Cache", "HIT")
                w.Header().Set("X-Cache-Date", cachedResponse.CachedAt.Format(time.RFC3339))
                
                // Set original headers
                for key, values := range cachedResponse.Headers {
                    for _, value := range values {
                        w.Header().Add(key, value)
                    }
                }
                
                w.WriteHeader(cachedResponse.StatusCode)
                w.Write(cachedResponse.Body)
                return
            }

            // Cache miss - execute request
            w.Header().Set("X-Cache", "MISS")
            rw := &responseWriter{
                ResponseWriter: w,
                statusCode:     http.StatusOK,
                body:           &bytes.Buffer{},
            }

            next.ServeHTTP(rw, r)

            // Only cache successful responses
            if rw.statusCode >= 200 && rw.statusCode < 300 {
                responseToCache := ResponseCache{
                    StatusCode: rw.statusCode,
                    Headers:    w.Header().Clone(),
                    Body:       rw.body.Bytes(),
                    CachedAt:   time.Now(),
                }

                // Cache the response asynchronously
                go func() {
                    cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                    defer cancel()
                    
                    if err := cache.DefaultCache.Set(cacheCtx, cacheKey, responseToCache, cache.APIResponseCacheConfig); err != nil {
                        // Log cache error but don't affect the response
                        fmt.Printf("Failed to cache response: %v\n", err)
                    }
                }()
            }
        }
    }
}

func generateCacheKey(r *http.Request) string {
    h := sha256.New()
    h.Write([]byte(r.URL.Path))
    h.Write([]byte(r.URL.RawQuery))
    
    // Include relevant headers in cache key
    cacheHeaders := []string{"Accept", "Accept-Language", "Accept-Encoding"}
    for _, header := range cacheHeaders {
        if value := r.Header.Get(header); value != "" {
            h.Write([]byte(header + ":" + value))
        }
    }
    
    return hex.EncodeToString(h.Sum(nil))[:16]
}
