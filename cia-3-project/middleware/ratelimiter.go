package middleware

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itcodehery/irctc-simulator/models"
)

// TokenBucket implements a per-IP token bucket rate limiter.
type TokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// RateLimiter manages rate limiting across multiple clients.
type RateLimiter struct {
	mu       sync.RWMutex
	buckets  map[string]*TokenBucket
	maxReqs  int
	window   time.Duration
	stopChan chan struct{}
}

// NewRateLimiter creates a new rate limiter with the given parameters.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		maxReqs:  maxRequests,
		window:   window,
		stopChan: make(chan struct{}),
	}

	// Background goroutine to clean up stale entries
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given IP is allowed.
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[ip]
	if !exists {
		bucket = &TokenBucket{
			tokens:     float64(rl.maxReqs),
			maxTokens:  float64(rl.maxReqs),
			refillRate: float64(rl.maxReqs) / rl.window.Seconds(),
			lastRefill: time.Now(),
		}
		rl.buckets[ip] = bucket
	}

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	// Check if we have tokens available
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup removes stale entries every 5 minutes.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, bucket := range rl.buckets {
				if now.Sub(bucket.lastRefill) > 10*time.Minute {
					delete(rl.buckets, ip)
				}
			}
			rl.mu.Unlock()
			slog.Info("Rate limiter cleanup completed", "active_buckets", len(rl.buckets))
		case <-rl.stopChan:
			return
		}
	}
}

// Stop shuts down the rate limiter's background goroutine.
func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// RateLimitMiddleware returns a Gin middleware that applies rate limiting.
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.Allow(ip) {
			slog.Warn("Rate limit exceeded",
				"client_ip", ip,
				"path", c.Request.URL.Path,
			)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error:   "rate limit exceeded",
				Code:    http.StatusTooManyRequests,
				Details: "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
