package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	config := RateLimitConfig{
		IPRequestsPerMinute:   5,
		UserRequestsPerMinute: 10,
		Enabled:               true,
	}

	handler := RateLimit(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Test IP rate limiting
	t.Run("IP rate limit", func(t *testing.T) {
		// Make 5 requests (should all succeed)
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Request %d: expected status 200, got %d", i+1, w.Code)
			}
		}

		// 6th request should be rate limited
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected rate limit (429), got %d", w.Code)
		}
	})

	t.Run("Different IPs not limited", func(t *testing.T) {
		// Requests from different IPs should not affect each other
		req1 := httptest.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = "192.168.1.2:12345"
		w1 := httptest.NewRecorder()

		handler.ServeHTTP(w1, req1)

		if w1.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w1.Code)
		}

		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = "192.168.1.3:12345"
		w2 := httptest.NewRecorder()

		handler.ServeHTTP(w2, req2)

		if w2.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w2.Code)
		}
	})
}

func TestRateLimitDisabled(t *testing.T) {
	config := RateLimitConfig{
		IPRequestsPerMinute:   5,
		UserRequestsPerMinute: 10,
		Enabled:               false,
	}

	handler := RateLimit(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Make many requests - none should be rate limited
	for i := 0; i < 20; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d (rate limiting should be disabled)", i+1, w.Code)
		}
	}
}

func TestTokenBucket(t *testing.T) {
	t.Run("Basic token bucket", func(t *testing.T) {
		bucket := NewTokenBucket(3, 1) // 3 tokens, 1 per second

		// Should allow 3 requests immediately
		for i := 0; i < 3; i++ {
			if !bucket.Allow() {
				t.Errorf("Request %d should be allowed", i+1)
			}
		}

		// 4th request should be blocked
		if bucket.Allow() {
			t.Error("4th request should be blocked")
		}
	})

	t.Run("Token refill", func(t *testing.T) {
		bucket := NewTokenBucket(2, 10) // 2 tokens, 10 per second

		// Use both tokens
		bucket.Allow()
		bucket.Allow()

		// Should be out of tokens
		if bucket.Allow() {
			t.Error("Should be out of tokens")
		}

		// Wait for tokens to refill (0.2 seconds = 2 tokens at 10/sec)
		time.Sleep(250 * time.Millisecond)

		// Should have tokens again
		if !bucket.Allow() {
			t.Error("Should have refilled tokens")
		}
	})
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		xRealIP        string
		expectedIP     string
	}{
		{
			name:       "RemoteAddr only",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1",
		},
		{
			name:          "X-Forwarded-For takes precedence",
			remoteAddr:    "192.168.1.1:12345",
			xForwardedFor: "10.0.0.1",
			expectedIP:    "10.0.0.1",
		},
		{
			name:          "X-Forwarded-For with multiple IPs",
			remoteAddr:    "192.168.1.1:12345",
			xForwardedFor: "10.0.0.1, 10.0.0.2, 10.0.0.3",
			expectedIP:    "10.0.0.1",
		},
		{
			name:          "X-Forwarded-For with spaces",
			remoteAddr:    "192.168.1.1:12345",
			xForwardedFor: "  10.0.0.1  , 10.0.0.2",
			expectedIP:    "10.0.0.1",
		},
		{
			name:       "X-Real-IP takes precedence over RemoteAddr",
			remoteAddr: "192.168.1.1:12345",
			xRealIP:    "10.0.0.2",
			expectedIP: "10.0.0.2",
		},
		{
			name:          "X-Forwarded-For over X-Real-IP",
			remoteAddr:    "192.168.1.1:12345",
			xForwardedFor: "10.0.0.1",
			xRealIP:       "10.0.0.2",
			expectedIP:    "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			ip := getClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}
