package services

import (
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	authService := NewAuthService("test-secret", 24*time.Hour)

	password := "testpassword123"
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	if hash == password {
		t.Fatal("Hash should not be the same as password")
	}
}

func TestVerifyPassword(t *testing.T) {
	authService := NewAuthService("test-secret", 24*time.Hour)

	password := "testpassword123"
	hash, err := authService.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	err = authService.VerifyPassword(hash, password)
	if err != nil {
		t.Errorf("Failed to verify correct password: %v", err)
	}

	// Test incorrect password
	err = authService.VerifyPassword(hash, "wrongpassword")
	if err == nil {
		t.Error("Should fail to verify incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	authService := NewAuthService("test-secret", 24*time.Hour)

	userID := "user-123"
	email := "test@example.com"

	token, err := authService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	authService := NewAuthService("test-secret", 24*time.Hour)

	userID := "user-123"
	email := "test@example.com"

	// Generate token
	token, err := authService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	authService := NewAuthService("test-secret", 24*time.Hour)

	// Test invalid token
	_, err := authService.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Should fail to validate invalid token")
	}
}

func TestValidateExpiredToken(t *testing.T) {
	// Create service with short expiry
	authService := NewAuthService("test-secret", 1*time.Millisecond)

	userID := "user-123"
	email := "test@example.com"

	// Generate token
	token, err := authService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Validate expired token
	_, err = authService.ValidateToken(token)
	if err == nil {
		t.Error("Should fail to validate expired token")
	}

	if err != ErrTokenExpired {
		t.Errorf("Expected ErrTokenExpired, got %v", err)
	}
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
	authService1 := NewAuthService("secret1", 24*time.Hour)
	authService2 := NewAuthService("secret2", 24*time.Hour)

	userID := "user-123"
	email := "test@example.com"

	// Generate token with first service
	token, err := authService1.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Try to validate with second service (different secret)
	_, err = authService2.ValidateToken(token)
	if err == nil {
		t.Error("Should fail to validate token with wrong secret")
	}
}
