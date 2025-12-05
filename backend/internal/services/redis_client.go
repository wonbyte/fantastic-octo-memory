package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client with connection management
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient() (*RedisClient, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	
	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     redisPassword,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis connection failed, caching will be disabled", "error", err, "addr", addr)
		// Don't return error - allow app to run without cache
		return &RedisClient{client: nil}, nil
	}

	slog.Info("Redis client initialized successfully", "addr", addr)
	return &RedisClient{client: client}, nil
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("redis client not available")
	}
	return r.client.Get(ctx, key).Result()
}

// Set stores a value in Redis with TTL
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from Redis
func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}
	return r.client.Del(ctx, keys...).Err()
}

// DeletePattern deletes all keys matching a pattern
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
	if r.client == nil {
		return fmt.Errorf("redis client not available")
	}

	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// IsAvailable checks if Redis is available
func (r *RedisClient) IsAvailable() bool {
	return r.client != nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}
