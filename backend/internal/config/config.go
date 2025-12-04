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
