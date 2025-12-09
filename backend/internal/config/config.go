package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	S3       S3Config
	AI       AIConfig
	Worker   WorkerConfig
	Auth     AuthConfig
	RateLimit RateLimitConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL            string
	MaxConnections int
	MaxIdleConns   int
}

type S3Config struct {
	Endpoint       string
	AccessKey      string
	SecretKey      string
	Bucket         string
	Region         string
	UsePathStyle   bool
	PresignExpiry  time.Duration
}

type AIConfig struct {
	ServiceURL string
	Timeout    time.Duration
}

type WorkerConfig struct {
	PollInterval time.Duration
	MaxRetries   int
}

type AuthConfig struct {
	JWTSecret   string
	TokenExpiry time.Duration
}

type RateLimitConfig struct {
	Enabled               bool
	IPRequestsPerMinute   int
	UserRequestsPerMinute int
}

type SecurityConfig struct {
	EnableSecurityHeaders bool
	EnableHSTS           bool
	HSTSMaxAge           int
	EnableCSP            bool
	CSPDirectives        string
	CORSAllowedOrigins   []string
	MaxRequestBodyBytes  int64
}

func Load() (*Config, error) {
	// Try to load .env file (optional in production)
	_ = godotenv.Load()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/construction_db?sslmode=disable")
	viper.SetDefault("S3_ENDPOINT", "http://localhost:9000")
	viper.SetDefault("S3_ACCESS_KEY", "minioadmin")
	viper.SetDefault("S3_SECRET_KEY", "minioadmin")
	viper.SetDefault("S3_BUCKET", "blueprints")
	viper.SetDefault("S3_REGION", "us-east-1")
	viper.SetDefault("S3_USE_PATH_STYLE", true)
	viper.SetDefault("S3_PRESIGN_EXPIRY", "5m")
	viper.SetDefault("AI_SERVICE_URL", "http://localhost:8000")
	viper.SetDefault("AI_SERVICE_TIMEOUT", "30s")
	viper.SetDefault("JOB_POLL_INTERVAL", "5s")
	viper.SetDefault("WORKER_MAX_RETRIES", 3)
	viper.SetDefault("DB_MAX_CONNECTIONS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNECTIONS", 5)
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_TOKEN_EXPIRY", "24h")
	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_IP_REQUESTS_PER_MIN", 100)
	viper.SetDefault("RATE_LIMIT_USER_REQUESTS_PER_MIN", 200)
	viper.SetDefault("ENABLE_SECURITY_HEADERS", true)
	viper.SetDefault("ENABLE_HSTS", true)
	viper.SetDefault("HSTS_MAX_AGE", 31536000)
	viper.SetDefault("ENABLE_CSP", true)
	viper.SetDefault("CSP_DIRECTIVES", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:19006")
	viper.SetDefault("MAX_REQUEST_BODY_BYTES", 10485760) // 10MB default

	// Auto bind environment variables
	viper.AutomaticEnv()

	// Parse durations
	presignExpiry, err := time.ParseDuration(viper.GetString("S3_PRESIGN_EXPIRY"))
	if err != nil {
		presignExpiry = 5 * time.Minute
		log.Printf("Warning: Invalid S3_PRESIGN_EXPIRY, using default: %s", presignExpiry)
	}

	aiTimeout, err := time.ParseDuration(viper.GetString("AI_SERVICE_TIMEOUT"))
	if err != nil {
		aiTimeout = 30 * time.Second
		log.Printf("Warning: Invalid AI_SERVICE_TIMEOUT, using default: %s", aiTimeout)
	}

	pollInterval, err := time.ParseDuration(viper.GetString("JOB_POLL_INTERVAL"))
	if err != nil {
		pollInterval = 5 * time.Second
		log.Printf("Warning: Invalid JOB_POLL_INTERVAL, using default: %s", pollInterval)
	}

	tokenExpiry, err := time.ParseDuration(viper.GetString("JWT_TOKEN_EXPIRY"))
	if err != nil {
		tokenExpiry = 24 * time.Hour
		log.Printf("Warning: Invalid JWT_TOKEN_EXPIRY, using default: %s", tokenExpiry)
	}

	// Parse CORS allowed origins
	corsOriginsStr := viper.GetString("CORS_ALLOWED_ORIGINS")
	corsOrigins := []string{}
	if corsOriginsStr != "" {
		for _, origin := range splitAndTrim(corsOriginsStr, ",") {
			if origin != "" {
				corsOrigins = append(corsOrigins, origin)
			}
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Env:  viper.GetString("ENV"),
		},
		Database: DatabaseConfig{
			URL:            viper.GetString("DATABASE_URL"),
			MaxConnections: viper.GetInt("DB_MAX_CONNECTIONS"),
			MaxIdleConns:   viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
		},
		S3: S3Config{
			Endpoint:      viper.GetString("S3_ENDPOINT"),
			AccessKey:     viper.GetString("S3_ACCESS_KEY"),
			SecretKey:     viper.GetString("S3_SECRET_KEY"),
			Bucket:        viper.GetString("S3_BUCKET"),
			Region:        viper.GetString("S3_REGION"),
			UsePathStyle:  viper.GetBool("S3_USE_PATH_STYLE"),
			PresignExpiry: presignExpiry,
		},
		AI: AIConfig{
			ServiceURL: viper.GetString("AI_SERVICE_URL"),
			Timeout:    aiTimeout,
		},
		Worker: WorkerConfig{
			PollInterval: pollInterval,
			MaxRetries:   viper.GetInt("WORKER_MAX_RETRIES"),
		},
		Auth: AuthConfig{
			JWTSecret:   viper.GetString("JWT_SECRET"),
			TokenExpiry: tokenExpiry,
		},
		RateLimit: RateLimitConfig{
			Enabled:               viper.GetBool("RATE_LIMIT_ENABLED"),
			IPRequestsPerMinute:   viper.GetInt("RATE_LIMIT_IP_REQUESTS_PER_MIN"),
			UserRequestsPerMinute: viper.GetInt("RATE_LIMIT_USER_REQUESTS_PER_MIN"),
		},
		Security: SecurityConfig{
			EnableSecurityHeaders: viper.GetBool("ENABLE_SECURITY_HEADERS"),
			EnableHSTS:           viper.GetBool("ENABLE_HSTS"),
			HSTSMaxAge:           viper.GetInt("HSTS_MAX_AGE"),
			EnableCSP:            viper.GetBool("ENABLE_CSP"),
			CSPDirectives:        viper.GetString("CSP_DIRECTIVES"),
			CORSAllowedOrigins:   corsOrigins,
			MaxRequestBodyBytes:  viper.GetInt64("MAX_REQUEST_BODY_BYTES"),
		},
	}

	// Validate required fields
	if config.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if config.Auth.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required - please set a secure secret in environment variables")
	}

	return config, nil
}

// splitAndTrim splits a string by delimiter and trims whitespace from each part
func splitAndTrim(s, delimiter string) []string {
	parts := []string{}
	for _, part := range splitString(s, delimiter) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, delimiter string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for i := 0; i < len(s); i++ {
		if i+len(delimiter) <= len(s) && s[i:i+len(delimiter)] == delimiter {
			result = append(result, current)
			current = ""
			i += len(delimiter) - 1
		} else {
			current += string(s[i])
		}
	}
	result = append(result, current)
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
