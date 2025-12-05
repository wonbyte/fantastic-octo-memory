package services

import (
	"context"
	"testing"
	"time"
)

func TestRedisClient_NewRedisClient(t *testing.T) {
	// Test that NewRedisClient doesn't fail even without Redis
	// It should return a client with nil internal client if Redis is unavailable
	client, err := NewRedisClient()
	if err != nil {
		t.Fatalf("NewRedisClient failed: %v", err)
	}
	
	// Verify client is created (may or may not be available)
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestRedisClient_Operations_WithoutRedis(t *testing.T) {
	// Create a client without Redis connection
	client := &RedisClient{client: nil}
	ctx := context.Background()
	
	// Test Get without Redis
	_, err := client.Get(ctx, "test:key")
	if err == nil {
		t.Error("Expected error when getting without Redis")
	}
	
	// Test Set without Redis
	err = client.Set(ctx, "test:key", "value", time.Minute)
	if err == nil {
		t.Error("Expected error when setting without Redis")
	}
	
	// Test Delete without Redis
	err = client.Delete(ctx, "test:key")
	if err == nil {
		t.Error("Expected error when deleting without Redis")
	}
	
	// Test DeletePattern without Redis
	err = client.DeletePattern(ctx, "test:*")
	if err == nil {
		t.Error("Expected error when deleting pattern without Redis")
	}
	
	// Test IsAvailable
	if client.IsAvailable() {
		t.Error("Expected IsAvailable to return false")
	}
	
	// Test Close (should not fail)
	err = client.Close()
	if err != nil {
		t.Errorf("Close should not fail: %v", err)
	}
}
