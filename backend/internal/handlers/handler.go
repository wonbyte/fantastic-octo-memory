package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

type Handler struct {
	db            *repository.Database
	projectRepo   *repository.ProjectRepository
	blueprintRepo *repository.BlueprintRepository
	jobRepo       *repository.JobRepository
	s3Service     *services.S3Service
	aiService     *services.AIService
}

func NewHandler(
	db *repository.Database,
	projectRepo *repository.ProjectRepository,
	blueprintRepo *repository.BlueprintRepository,
	jobRepo *repository.JobRepository,
	s3Service *services.S3Service,
	aiService *services.AIService,
) *Handler {
	return &Handler{
		db:            db,
		projectRepo:   projectRepo,
		blueprintRepo: blueprintRepo,
		jobRepo:       jobRepo,
		s3Service:     s3Service,
		aiService:     aiService,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	// Check database health
	if err := h.db.Health(r.Context()); err != nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  "unhealthy",
			"version": "1.0.0",
			"error":   "database unavailable",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"version": "1.0.0",
	})
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Construction Estimation & Bidding Automation API",
		"version": "1.0.0",
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log encoding error but don't panic - response has already been written
		slog.Error("Failed to encode JSON response", "error", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
