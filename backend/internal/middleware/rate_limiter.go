package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// TokenBucket implements a token bucket algorithm for rate limiting
type TokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed based on available tokens
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(tb.capacity, tb.tokens+elapsed*tb.refillRate)
	tb.lastRefill = now

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// RateLimiter manages rate limiting for different IPs and users
type RateLimiter struct {
	ipBuckets     map[string]*TokenBucket
	userBuckets   map[string]*TokenBucket
	mu            sync.RWMutex
	ipCapacity    float64
	ipRefillRate  float64
	userCapacity  float64
	userRefillRate float64
	cleanupInterval time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(ipRequestsPerMin, userRequestsPerMin int) *RateLimiter {
	rl := &RateLimiter{
		ipBuckets:       make(map[string]*TokenBucket),
		userBuckets:     make(map[string]*TokenBucket),
		ipCapacity:      float64(ipRequestsPerMin),
		ipRefillRate:    float64(ipRequestsPerMin) / 60.0, // tokens per second
		userCapacity:    float64(userRequestsPerMin),
		userRefillRate:  float64(userRequestsPerMin) / 60.0,
		cleanupInterval: 10 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes old buckets to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// Clean up IP buckets that haven't been used in a while
		for ip, bucket := range rl.ipBuckets {
			bucket.mu.Lock()
			if time.Since(bucket.lastRefill) > rl.cleanupInterval {
				delete(rl.ipBuckets, ip)
			}
			bucket.mu.Unlock()
		}
		// Clean up user buckets that haven't been used in a while
		for userID, bucket := range rl.userBuckets {
			bucket.mu.Lock()
			if time.Since(bucket.lastRefill) > rl.cleanupInterval {
				delete(rl.userBuckets, userID)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// getIPBucket gets or creates a token bucket for an IP address
func (rl *RateLimiter) getIPBucket(ip string) *TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.ipBuckets[ip]
	rl.mu.RUnlock()

	if exists {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	bucket, exists = rl.ipBuckets[ip]
	if exists {
		return bucket
	}

	bucket = NewTokenBucket(rl.ipCapacity, rl.ipRefillRate)
	rl.ipBuckets[ip] = bucket
	return bucket
}

// getUserBucket gets or creates a token bucket for a user
func (rl *RateLimiter) getUserBucket(userID string) *TokenBucket {
	rl.mu.RLock()
	bucket, exists := rl.userBuckets[userID]
	rl.mu.RUnlock()

	if exists {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Double-check after acquiring write lock
	bucket, exists = rl.userBuckets[userID]
	if exists {
		return bucket
	}

	bucket = NewTokenBucket(rl.userCapacity, rl.userRefillRate)
	rl.userBuckets[userID] = bucket
	return bucket
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	IPRequestsPerMinute   int
	UserRequestsPerMinute int
	Enabled               bool
}

// RateLimit creates a rate limiting middleware
func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	if !config.Enabled {
		// Return a no-op middleware if rate limiting is disabled
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	limiter := NewRateLimiter(config.IPRequestsPerMinute, config.UserRequestsPerMinute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get correlation ID from context
			correlationID := ""
			if val := r.Context().Value(ContextKeyCorrelationID); val != nil {
				if id, ok := val.(string); ok {
					correlationID = id
				}
			}

			// Extract client IP
			clientIP := getClientIP(r)

			// Check IP-based rate limit
			ipBucket := limiter.getIPBucket(clientIP)
			if !ipBucket.Allow() {
				slog.Warn("Rate limit exceeded for IP",
					"ip", clientIP,
					"path", r.URL.Path,
					"correlation_id", correlationID)

				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.IPRequestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"Rate limit exceeded. Please try again later."}`))
				return
			}

			// Check user-based rate limit if user is authenticated
			userID := ""
			if val := r.Context().Value(ContextKeyUserID); val != nil {
				if id, ok := val.(string); ok {
					userID = id
				}
			}

			if userID != "" {
				userBucket := limiter.getUserBucket(userID)
				if !userBucket.Allow() {
					slog.Warn("Rate limit exceeded for user",
						"user_id", userID,
						"ip", clientIP,
						"path", r.URL.Path,
						"correlation_id", correlationID)

					w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.UserRequestsPerMinute))
					w.Header().Set("X-RateLimit-Remaining", "0")
					w.Header().Set("Retry-After", "60")
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte(`{"error":"Rate limit exceeded. Please try again later."}`))
					return
				}
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.IPRequestsPerMinute))

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		// If no port, just use the value as-is
		return xff
	}

	// Check X-Real-IP header (set by some proxies)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
