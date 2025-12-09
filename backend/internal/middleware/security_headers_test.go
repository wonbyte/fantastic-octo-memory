package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	config := DefaultSecurityHeadersConfig()

	handler := SecurityHeaders(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	tests := []struct {
		header   string
		expected string
	}{
		{"Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"},
		{"Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';"},
		{"X-Frame-Options", "DENY"},
		{"X-Content-Type-Options", "nosniff"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{"X-XSS-Protection", "1; mode=block"},
		{"Permissions-Policy", "geolocation=(), microphone=(), camera=()"},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			got := w.Header().Get(tt.header)
			if got != tt.expected {
				t.Errorf("Header %s: expected %q, got %q", tt.header, tt.expected, got)
			}
		})
	}
}

func TestSecurityHeadersDisabled(t *testing.T) {
	config := SecurityHeadersConfig{
		EnableHSTS:           false,
		EnableCSP:            false,
		EnableXFrameOptions:  false,
		EnableXContentType:   false,
		EnableReferrerPolicy: false,
	}

	handler := SecurityHeaders(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// These headers should not be set when disabled
	if w.Header().Get("Strict-Transport-Security") != "" {
		t.Error("HSTS header should not be set when disabled")
	}
	if w.Header().Get("Content-Security-Policy") != "" {
		t.Error("CSP header should not be set when disabled")
	}
	if w.Header().Get("X-Frame-Options") != "" {
		t.Error("X-Frame-Options should not be set when disabled")
	}
	if w.Header().Get("X-Content-Type-Options") != "" {
		t.Error("X-Content-Type-Options should not be set when disabled")
	}
	if w.Header().Get("Referrer-Policy") != "" {
		t.Error("Referrer-Policy should not be set when disabled")
	}

	// These headers should always be set
	if w.Header().Get("X-XSS-Protection") == "" {
		t.Error("X-XSS-Protection should always be set")
	}
	if w.Header().Get("Permissions-Policy") == "" {
		t.Error("Permissions-Policy should always be set")
	}
}

func TestSecurityHeadersCustom(t *testing.T) {
	config := SecurityHeadersConfig{
		EnableHSTS:           true,
		HSTSMaxAge:           3600,
		EnableCSP:            true,
		CSPDirectives:        "default-src 'none';",
		EnableXFrameOptions:  true,
		XFrameOptionsValue:   "SAMEORIGIN",
		EnableXContentType:   true,
		EnableReferrerPolicy: true,
		ReferrerPolicyValue:  "no-referrer",
	}

	handler := SecurityHeaders(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if got := w.Header().Get("Strict-Transport-Security"); got != "max-age=3600; includeSubDomains; preload" {
		t.Errorf("Expected custom HSTS, got %q", got)
	}

	if got := w.Header().Get("Content-Security-Policy"); got != "default-src 'none';" {
		t.Errorf("Expected custom CSP, got %q", got)
	}

	if got := w.Header().Get("X-Frame-Options"); got != "SAMEORIGIN" {
		t.Errorf("Expected SAMEORIGIN, got %q", got)
	}

	if got := w.Header().Get("Referrer-Policy"); got != "no-referrer" {
		t.Errorf("Expected no-referrer, got %q", got)
	}
}
