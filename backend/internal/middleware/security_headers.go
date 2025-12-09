package middleware

import (
	"fmt"
	"net/http"
)

// SecurityHeadersConfig holds configuration for security headers
type SecurityHeadersConfig struct {
	EnableHSTS           bool
	HSTSMaxAge           int
	EnableCSP            bool
	CSPDirectives        string
	EnableXFrameOptions  bool
	XFrameOptionsValue   string
	EnableXContentType   bool
	EnableReferrerPolicy bool
	ReferrerPolicyValue  string
}

// DefaultSecurityHeadersConfig returns default security headers configuration
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		EnableHSTS:           true,
		HSTSMaxAge:           31536000, // 1 year
		EnableCSP:            true,
		CSPDirectives:        "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';",
		EnableXFrameOptions:  true,
		XFrameOptionsValue:   "DENY",
		EnableXContentType:   true,
		EnableReferrerPolicy: true,
		ReferrerPolicyValue:  "strict-origin-when-cross-origin",
	}
}

// SecurityHeaders adds security headers to HTTP responses
func SecurityHeaders(config SecurityHeadersConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// HSTS (HTTP Strict Transport Security)
			// Forces browsers to use HTTPS for future requests
			if config.EnableHSTS {
				hstsValue := fmt.Sprintf("max-age=%d; includeSubDomains; preload", config.HSTSMaxAge)
				w.Header().Set("Strict-Transport-Security", hstsValue)
			}

			// CSP (Content Security Policy)
			// Prevents XSS attacks by controlling what resources can be loaded
			if config.EnableCSP {
				w.Header().Set("Content-Security-Policy", config.CSPDirectives)
			}

			// X-Frame-Options
			// Prevents clickjacking attacks by controlling iframe embedding
			if config.EnableXFrameOptions {
				w.Header().Set("X-Frame-Options", config.XFrameOptionsValue)
			}

			// X-Content-Type-Options
			// Prevents MIME type sniffing
			if config.EnableXContentType {
				w.Header().Set("X-Content-Type-Options", "nosniff")
			}

			// Referrer-Policy
			// Controls how much referrer information is included with requests
			if config.EnableReferrerPolicy {
				w.Header().Set("Referrer-Policy", config.ReferrerPolicyValue)
			}

			// X-XSS-Protection (legacy, but still good to have for older browsers)
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Permissions-Policy (formerly Feature-Policy)
			// Controls which browser features can be used
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			next.ServeHTTP(w, r)
		})
	}
}
