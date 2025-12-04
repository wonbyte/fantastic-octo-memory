package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
)

type SignupRequest struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Name        *string `json:"name"`
	CompanyName *string `json:"company_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Name        *string `json:"name"`
	CompanyName *string `json:"company_name"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Signup handles user registration
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID := getCorrelationID(ctx)

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode signup request",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		slog.Error("Failed to hash password",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		CompanyName:  req.CompanyName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.userRepo.CreateUser(ctx, user); err != nil {
		if err == repository.ErrEmailAlreadyExists {
			respondError(w, http.StatusConflict, "Email already exists")
			return
		}
		slog.Error("Failed to create user",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		slog.Error("Failed to generate token",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	slog.Info("User registered successfully",
		"user_id", user.ID,
		"email", user.Email,
		"correlation_id", correlationID)

	respondJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			Name:        user.Name,
			CompanyName: user.CompanyName,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// Login handles user authentication
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID := getCorrelationID(ctx)

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode login request",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user by email
	user, err := h.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			respondError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		slog.Error("Failed to get user by email",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	// Verify password
	if err := h.authService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		slog.Error("Failed to generate token",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	slog.Info("User logged in successfully",
		"user_id", user.ID,
		"email", user.Email,
		"correlation_id", correlationID)

	respondJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User: UserResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			Name:        user.Name,
			CompanyName: user.CompanyName,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// GetCurrentUser returns the authenticated user's information
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID := getCorrelationID(ctx)
	userID := getUserID(ctx)

	if userID == "" {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("Failed to parse user ID",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetUserByID(ctx, uid)
	if err != nil {
		if err == repository.ErrUserNotFound {
			respondError(w, http.StatusNotFound, "User not found")
			return
		}
		slog.Error("Failed to get user",
			"error", err,
			"correlation_id", correlationID)
		respondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	respondJSON(w, http.StatusOK, UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		Name:        user.Name,
		CompanyName: user.CompanyName,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	})
}
