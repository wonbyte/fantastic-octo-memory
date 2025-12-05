package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

// GetMaterials returns all materials, optionally filtered by category and region
func (h *Handler) GetMaterials(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	region := r.URL.Query().Get("region")

	var categoryPtr, regionPtr *string
	if category != "" {
		categoryPtr = &category
	}
	if region != "" {
		regionPtr = &region
	}

	materials, err := h.materialRepo.GetAll(r.Context(), categoryPtr, regionPtr)
	if err != nil {
		slog.Error("Failed to get materials", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get materials")
		return
	}

	respondJSON(w, http.StatusOK, materials)
}

// GetLaborRates returns all labor rates, optionally filtered by trade and region
func (h *Handler) GetLaborRates(w http.ResponseWriter, r *http.Request) {
	trade := r.URL.Query().Get("trade")
	region := r.URL.Query().Get("region")

	var tradePtr, regionPtr *string
	if trade != "" {
		tradePtr = &trade
	}
	if region != "" {
		regionPtr = &region
	}

	rates, err := h.laborRateRepo.GetAll(r.Context(), tradePtr, regionPtr)
	if err != nil {
		slog.Error("Failed to get labor rates", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get labor rates")
		return
	}

	respondJSON(w, http.StatusOK, rates)
}

// GetRegionalAdjustments returns all regional adjustments
func (h *Handler) GetRegionalAdjustments(w http.ResponseWriter, r *http.Request) {
	adjustments, err := h.regionalRepo.GetAll(r.Context())
	if err != nil {
		slog.Error("Failed to get regional adjustments", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get regional adjustments")
		return
	}

	respondJSON(w, http.StatusOK, adjustments)
}

// GetCompanyPricingOverrides returns all pricing overrides for the authenticated user
func (h *Handler) GetCompanyPricingOverrides(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	overrides, err := h.companyOverrideRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to get pricing overrides", "user_id", userID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get pricing overrides")
		return
	}

	respondJSON(w, http.StatusOK, overrides)
}

// CreateCompanyPricingOverrideRequest represents a request to create a pricing override
type CreateCompanyPricingOverrideRequest struct {
	OverrideType  string  `json:"override_type"`
	ItemKey       string  `json:"item_key"`
	OverrideValue float64 `json:"override_value"`
	IsPercentage  bool    `json:"is_percentage"`
	Notes         *string `json:"notes"`
}

// CreateCompanyPricingOverride creates a new pricing override for the authenticated user
func (h *Handler) CreateCompanyPricingOverride(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)

	var req CreateCompanyPricingOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate override type
	validTypes := map[string]bool{
		"material":      true,
		"labor":         true,
		"overhead":      true,
		"profit_margin": true,
	}
	if !validTypes[req.OverrideType] {
		respondError(w, http.StatusBadRequest, "Invalid override type")
		return
	}

	// Check if override already exists
	existing, err := h.companyOverrideRepo.GetByUserIDTypeAndKey(r.Context(), userID, req.OverrideType, req.ItemKey)
	if err == nil && existing != nil {
		respondError(w, http.StatusConflict, "Override already exists for this item")
		return
	}

	now := time.Now()
	override := &models.CompanyPricingOverride{
		ID:            uuid.New(),
		UserID:        userID,
		OverrideType:  req.OverrideType,
		ItemKey:       req.ItemKey,
		OverrideValue: req.OverrideValue,
		IsPercentage:  req.IsPercentage,
		Notes:         req.Notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := h.companyOverrideRepo.Create(r.Context(), override); err != nil {
		slog.Error("Failed to create pricing override", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to create pricing override")
		return
	}

	respondJSON(w, http.StatusCreated, override)
}

// UpdateCompanyPricingOverrideRequest represents a request to update a pricing override
type UpdateCompanyPricingOverrideRequest struct {
	OverrideValue float64 `json:"override_value"`
	IsPercentage  bool    `json:"is_percentage"`
	Notes         *string `json:"notes"`
}

// UpdateCompanyPricingOverride updates a pricing override
func (h *Handler) UpdateCompanyPricingOverride(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)
	overrideID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid override ID")
		return
	}

	// Get existing override
	override, err := h.companyOverrideRepo.GetByID(r.Context(), overrideID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Override not found")
		return
	}

	// Verify ownership
	if override.UserID != userID {
		respondError(w, http.StatusForbidden, "You don't have permission to update this override")
		return
	}

	var req UpdateCompanyPricingOverrideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update fields
	override.OverrideValue = req.OverrideValue
	override.IsPercentage = req.IsPercentage
	override.Notes = req.Notes
	override.UpdatedAt = time.Now()

	if err := h.companyOverrideRepo.Update(r.Context(), override); err != nil {
		slog.Error("Failed to update pricing override", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to update pricing override")
		return
	}

	respondJSON(w, http.StatusOK, override)
}

// DeleteCompanyPricingOverride deletes a pricing override
func (h *Handler) DeleteCompanyPricingOverride(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)
	overrideID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid override ID")
		return
	}

	// Get existing override
	override, err := h.companyOverrideRepo.GetByID(r.Context(), overrideID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Override not found")
		return
	}

	// Verify ownership
	if override.UserID != userID {
		respondError(w, http.StatusForbidden, "You don't have permission to delete this override")
		return
	}

	if err := h.companyOverrideRepo.Delete(r.Context(), overrideID); err != nil {
		slog.Error("Failed to delete pricing override", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete pricing override")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SyncCostDataRequest represents a request to sync cost data from external providers
type SyncCostDataRequest struct {
	Provider string `json:"provider"`
	Region   string `json:"region"`
}

// SyncCostData syncs cost data from external providers (admin only)
func (h *Handler) SyncCostData(w http.ResponseWriter, r *http.Request) {
	var req SyncCostDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Region == "" {
		req.Region = "national"
	}

	// Sync based on provider
	switch req.Provider {
	case "all":
		if err := h.costIntegrationService.SyncAll(r.Context(), req.Region); err != nil {
			slog.Error("Failed to sync all cost data", "error", err)
			respondError(w, http.StatusInternalServerError, "Failed to sync cost data")
			return
		}
	case "rsmeans", "homedepot", "lowes":
		if err := h.costIntegrationService.SyncMaterials(r.Context(), req.Provider, req.Region); err != nil {
			slog.Error("Failed to sync materials", "provider", req.Provider, "error", err)
			respondError(w, http.StatusInternalServerError, "Failed to sync materials")
			return
		}
		if err := h.costIntegrationService.SyncLaborRates(r.Context(), req.Provider, req.Region); err != nil {
			slog.Error("Failed to sync labor rates", "provider", req.Provider, "error", err)
			respondError(w, http.StatusInternalServerError, "Failed to sync labor rates")
			return
		}
		if err := h.costIntegrationService.SyncRegionalAdjustment(r.Context(), req.Provider, req.Region); err != nil {
			slog.Error("Failed to sync regional adjustment", "provider", req.Provider, "error", err)
			respondError(w, http.StatusInternalServerError, "Failed to sync regional adjustment")
			return
		}
	default:
		respondError(w, http.StatusBadRequest, "Invalid provider")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Cost data synced successfully",
	})
}
