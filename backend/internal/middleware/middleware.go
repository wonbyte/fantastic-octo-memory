package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type contextKey string

const (
	ContextKeyUserID        contextKey = "user_id"
	ContextKeyEmail         contextKey = "email"
	ContextKeyCorrelationID contextKey = "correlation_id"
)

// CorrelationID middleware adds a correlation ID to each request
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for existing correlation ID in header
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add correlation ID to response header
		w.Header().Set("X-Correlation-ID", correlationID)

		// Add correlation ID to context
		ctx := context.WithValue(r.Context(), ContextKeyCorrelationID, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get correlation ID from context
		correlationID := ""
		if val := r.Context().Value(ContextKeyCorrelationID); val != nil {
			correlationID = val.(string)
		}

		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request
		duration := time.Since(start)
		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
			"correlation_id", correlationID,
		)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In production, this should be configured via environment variables
		// For now, allow localhost origins for development
		origin := r.Header.Get("Origin")
		if origin != "" {
			// Allow common development origins
			// TODO: Make this configurable via environment variable
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CORSWithConfig creates a CORS middleware with configurable origins
func CORSWithConfig(allowedOrigins []string) func(http.Handler) http.Handler {
	// Create a map for faster origin lookups
	originMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originMap[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Check if origin is in allowed list
			if origin != "" && originMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else if len(allowedOrigins) == 0 {
				// If no origins configured, use permissive CORS for development
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Correlation-ID")
			w.Header().Set("Access-Control-Expose-Headers", "X-Correlation-ID, X-RateLimit-Limit, X-RateLimit-Remaining")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Get correlation ID from context
				correlationID := ""
				if val := r.Context().Value(ContextKeyCorrelationID); val != nil {
					correlationID = val.(string)
				}

				slog.Error("Panic recovered",
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
					"correlation_id", correlationID,
				)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RequestBodyLimit limits the size of request bodies
func RequestBodyLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Limit request body size
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			
			next.ServeHTTP(w, r)
		})
	}
}

// Auth middleware validates JWT tokens and adds user info to context
func Auth(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get correlation ID from context
			correlationID := ""
			if val := r.Context().Value(ContextKeyCorrelationID); val != nil {
				correlationID = val.(string)
			}

			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.Warn("Missing authorization header",
					"path", r.URL.Path,
					"correlation_id", correlationID)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Missing authorization header"}`))
				return
			}

			// Extract bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				slog.Warn("Invalid authorization header format",
					"path", r.URL.Path,
					"correlation_id", correlationID)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Invalid authorization header format"}`))
				return
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				slog.Warn("Invalid token",
					"error", err,
					"path", r.URL.Path,
					"correlation_id", correlationID)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"Invalid or expired token"}`))
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
