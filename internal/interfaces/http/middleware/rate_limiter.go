package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitorInfo tracks limiter and last access time
type visitorInfo struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// RateLimiter implements token bucket algorithm for rate limiting
type RateLimiter struct {
	visitors map[string]*visitorInfo
	mu       sync.RWMutex
	r        rate.Limit // requests per second
	b        int        // burst size
}

// NewRateLimiter creates a new rate limiter
// r: requests per second
// b: burst size (maximum tokens)
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	limiter := &RateLimiter{
		visitors: make(map[string]*visitorInfo),
		r:        r,
		b:        b,
	}

	// Cleanup old entries every minute
	go limiter.cleanupVisitors()

	return limiter
}

// getVisitor returns the rate limiter for the given key (IP or user)
func (rl *RateLimiter) getVisitor(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, exists := rl.visitors[key]
	if !exists {
		info = &visitorInfo{
			limiter:    rate.NewLimiter(rl.r, rl.b),
			lastAccess: time.Now(),
		}
		rl.visitors[key] = info
	} else {
		info.lastAccess = time.Now()
	}

	return info.limiter
}

// cleanupVisitors removes inactive limiters (not accessed in 3 minutes)
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		threshold := time.Now().Add(-3 * time.Minute)
		for key, info := range rl.visitors {
			if info.lastAccess.Before(threshold) {
				delete(rl.visitors, key)
			}
		}
		rl.mu.Unlock()
	}
}

// IPRateLimiter middleware limits requests per IP address
func IPRateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := NewRateLimiter(r, b)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiterForIP := limiter.getVisitor(ip)

		if !limiterForIP.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimiter middleware limits requests per authenticated user
func UserRateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := NewRateLimiter(r, b)

	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// If not authenticated, fall back to IP-based limiting
			ip := c.ClientIP()
			limiterForIP := limiter.getVisitor(ip)

			if !limiterForIP.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded. Please try again later.",
					"code":  "RATE_LIMIT_EXCEEDED",
				})
				c.Abort()
				return
			}

			c.Next()
			return
		}

		// Use user ID for rate limiting
		limiterForUser := limiter.getVisitor(userID.(string))

		if !limiterForUser.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshTokenRateLimiter limits refresh endpoint to 10 requests per 5 minutes per IP
// This uses a simple counter-based approach for longer time windows
func RefreshTokenRateLimiter() gin.HandlerFunc {
	type requestCounter struct {
		count     int
		resetTime time.Time
	}

	requests := make(map[string]*requestCounter)
	mu := sync.RWMutex{}
	limit := 10
	window := 5 * time.Minute

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for key, counter := range requests {
				if now.After(counter.resetTime) {
					delete(requests, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		counter, exists := requests[clientIP]

		if !exists || now.After(counter.resetTime) {
			requests[clientIP] = &requestCounter{
				count:     1,
				resetTime: now.Add(window),
			}
			c.Next()
			return
		}

		if counter.count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many refresh requests. Please try again in a few minutes.",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		counter.count++
		c.Next()
	}
}
