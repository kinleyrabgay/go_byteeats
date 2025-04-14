package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"go_byteeats/internal/database"
	"go_byteeats/pkg/logger"

	"github.com/gin-gonic/gin"
)

type CacheConfig struct {
	Expiration time.Duration
	KeyPrefix  string
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func CacheMiddleware(redis *database.RedisClient, config CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip caching for non-GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Generate cache key
		key := generateCacheKey(config.KeyPrefix, c.Request)

		// Try to get from cache
		var cachedResponse []byte
		err := redis.Get(c.Request.Context(), key, &cachedResponse)
		if err == nil && len(cachedResponse) > 0 {
			c.Data(http.StatusOK, "application/json", cachedResponse)
			c.Abort()
			return
		}

		// Create custom response writer to capture response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = w

		// Process request
		c.Next()

		// Cache response if status is 200
		if c.Writer.Status() == http.StatusOK {
			err := redis.Set(
				context.Background(),
				key,
				w.body.Bytes(),
				config.Expiration,
			)
			if err != nil {
				logger.Error(err, "Failed to cache response")
			}
		}
	}
}

func generateCacheKey(prefix string, r *http.Request) string {
	// Create unique key based on URL path and query parameters
	h := sha256.New()
	io.WriteString(h, r.URL.Path)
	io.WriteString(h, r.URL.RawQuery)

	hash := hex.EncodeToString(h.Sum(nil))
	return prefix + ":" + hash
}
