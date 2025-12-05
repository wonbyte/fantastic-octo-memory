package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/middleware"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

// CostIntegrationServiceInterface defines the interface for cost integration service
type CostIntegrationServiceInterface interface {
	SyncMaterials(ctx context.Context, providerName, region string) error
	SyncLaborRates(ctx context.Context, providerName, region string) error
	SyncRegionalAdjustment(ctx context.Context, providerName, region string) error
	SyncAll(ctx context.Context, region string) error
}

// CostDataServiceInterface defines the interface for cost data retrieval (with or without cache)
type CostDataServiceInterface interface {
	GetMaterials(ctx context.Context, category, region *string) ([]models.MaterialCost, error)
	GetLaborRates(ctx context.Context, trade, region *string) ([]models.LaborRate, error)
	GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error)
}

type Handler struct {
	db                       *repository.Database
	projectRepo              *repository.ProjectRepository
	blueprintRepo            *repository.BlueprintRepository
	blueprintRevisionRepo    *repository.BlueprintRevisionRepository
	jobRepo                  *repository.JobRepository
	bidRepo                  *repository.BidRepository
	bidRevisionRepo          *repository.BidRevisionRepository
	userRepo                 *repository.UserRepository
	materialRepo             *repository.MaterialRepository
	laborRateRepo            *repository.LaborRateRepository
	regionalRepo             *repository.RegionalAdjustmentRepository
	companyOverrideRepo      *repository.CompanyPricingOverrideRepository
	s3Service                *services.S3Service
	aiService                *services.AIService
	authService              *services.AuthService
	costIntegrationService   CostIntegrationServiceInterface
	costDataService          CostDataServiceInterface
}

func NewHandler(
	db *repository.Database,
	projectRepo *repository.ProjectRepository,
	blueprintRepo *repository.BlueprintRepository,
	blueprintRevisionRepo *repository.BlueprintRevisionRepository,
	jobRepo *repository.JobRepository,
	bidRepo *repository.BidRepository,
	bidRevisionRepo *repository.BidRevisionRepository,
	userRepo *repository.UserRepository,
	materialRepo *repository.MaterialRepository,
	laborRateRepo *repository.LaborRateRepository,
	regionalRepo *repository.RegionalAdjustmentRepository,
	companyOverrideRepo *repository.CompanyPricingOverrideRepository,
	s3Service *services.S3Service,
	aiService *services.AIService,
	authService *services.AuthService,
	costIntegrationService CostIntegrationServiceInterface,
) *Handler {
	// Use costIntegrationService as costDataService if it supports the interface
	var costDataService CostDataServiceInterface
	if cds, ok := costIntegrationService.(CostDataServiceInterface); ok {
		costDataService = cds
	} else {
		// Fallback to nil - handlers will use repositories directly
		slog.Warn("CostIntegrationService does not implement CostDataServiceInterface, handlers will use direct repository access")
	}
	
	return &Handler{
		db:                       db,
		projectRepo:              projectRepo,
		blueprintRepo:            blueprintRepo,
		blueprintRevisionRepo:    blueprintRevisionRepo,
		jobRepo:                  jobRepo,
		bidRepo:                  bidRepo,
		bidRevisionRepo:          bidRevisionRepo,
		userRepo:                 userRepo,
		materialRepo:             materialRepo,
		laborRateRepo:            laborRateRepo,
		regionalRepo:             regionalRepo,
		companyOverrideRepo:      companyOverrideRepo,
		s3Service:                s3Service,
		aiService:                aiService,
		authService:              authService,
		costIntegrationService:   costIntegrationService,
		costDataService:          costDataService,
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
