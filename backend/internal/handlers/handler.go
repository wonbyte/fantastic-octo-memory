package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/middleware"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

type Handler struct {
	db                     *repository.Database
	projectRepo            *repository.ProjectRepository
	blueprintRepo          *repository.BlueprintRepository
	jobRepo                *repository.JobRepository
	bidRepo                *repository.BidRepository
	userRepo               *repository.UserRepository
	materialRepo           *repository.MaterialRepository
	laborRateRepo          *repository.LaborRateRepository
	regionalRepo           *repository.RegionalAdjustmentRepository
	companyOverrideRepo    *repository.CompanyPricingOverrideRepository
	s3Service              *services.S3Service
	aiService              *services.AIService
	authService            *services.AuthService
	costIntegrationService *services.CostIntegrationService
}

func NewHandler(
	db *repository.Database,
	projectRepo *repository.ProjectRepository,
	blueprintRepo *repository.BlueprintRepository,
	jobRepo *repository.JobRepository,
	bidRepo *repository.BidRepository,
	userRepo *repository.UserRepository,
	materialRepo *repository.MaterialRepository,
	laborRateRepo *repository.LaborRateRepository,
	regionalRepo *repository.RegionalAdjustmentRepository,
	companyOverrideRepo *repository.CompanyPricingOverrideRepository,
	s3Service *services.S3Service,
	aiService *services.AIService,
	authService *services.AuthService,
	costIntegrationService *services.CostIntegrationService,
) *Handler {
	return &Handler{
		db:                     db,
		projectRepo:            projectRepo,
		blueprintRepo:          blueprintRepo,
		jobRepo:                jobRepo,
		bidRepo:                bidRepo,
		userRepo:               userRepo,
		materialRepo:           materialRepo,
		laborRateRepo:          laborRateRepo,
		regionalRepo:           regionalRepo,
		companyOverrideRepo:    companyOverrideRepo,
		s3Service:              s3Service,
		aiService:              aiService,
		authService:            authService,
		costIntegrationService: costIntegrationService,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	healthStatus := map[string]interface{}{
		"status":  "ok",
		"version": "1.0.0",
	}

	// Check database health
	if err := h.db.Health(ctx); err != nil {
		healthStatus["status"] = "unhealthy"
		healthStatus["database"] = "unavailable"
		healthStatus["error"] = "database unavailable"
		respondJSON(w, http.StatusServiceUnavailable, healthStatus)
		return
	}
	healthStatus["database"] = "ok"

	// Check AI service health (optional - don't fail health check if AI service is down)
	if err := h.aiService.Health(ctx); err != nil {
		healthStatus["ai_service"] = "degraded"
	} else {
		healthStatus["ai_service"] = "ok"
	}

	respondJSON(w, http.StatusOK, healthStatus)
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

// Helper functions to extract values from context
func getUserID(ctx context.Context) string {
	if val := ctx.Value(middleware.ContextKeyUserID); val != nil {
		return val.(string)
	}
	return ""
}

func getEmail(ctx context.Context) string {
	if val := ctx.Value(middleware.ContextKeyEmail); val != nil {
		return val.(string)
	}
	return ""
}

func getCorrelationID(ctx context.Context) string {
	if val := ctx.Value(middleware.ContextKeyCorrelationID); val != nil {
		return val.(string)
	}
	return ""
}
