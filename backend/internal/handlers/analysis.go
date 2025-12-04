package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

// GetBlueprintAnalysis returns the normalized analysis data for a blueprint
func (h *Handler) GetBlueprintAnalysis(w http.ResponseWriter, r *http.Request) {
	blueprintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid blueprint ID")
		return
	}

	// Get blueprint record
	blueprint, err := h.blueprintRepo.GetByID(r.Context(), blueprintID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Blueprint not found")
		return
	}

	// Check if analysis data exists
	if blueprint.AnalysisData == nil || *blueprint.AnalysisData == "" {
		respondError(w, http.StatusNotFound, "Analysis data not available")
		return
	}

	// Parse analysis data
	var analysisResult models.AnalysisResult
	if err := json.Unmarshal([]byte(*blueprint.AnalysisData), &analysisResult); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to parse analysis data")
		return
	}

	respondJSON(w, http.StatusOK, analysisResult)
}

// GetBlueprintTakeoffSummary returns the calculated takeoff summary for a blueprint
func (h *Handler) GetBlueprintTakeoffSummary(w http.ResponseWriter, r *http.Request) {
	blueprintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid blueprint ID")
		return
	}

	// Get blueprint record
	blueprint, err := h.blueprintRepo.GetByID(r.Context(), blueprintID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Blueprint not found")
		return
	}

	// Check if analysis data exists
	if blueprint.AnalysisData == nil || *blueprint.AnalysisData == "" {
		respondError(w, http.StatusNotFound, "Analysis data not available")
		return
	}

	// Parse analysis data
	takeoffService := services.NewTakeoffService()
	analysisResult, err := takeoffService.ParseAnalysisData(*blueprint.AnalysisData)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to parse analysis data")
		return
	}

	// Calculate takeoff summary
	summary, err := takeoffService.CalculateTakeoffSummary(analysisResult)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to calculate takeoff summary")
		return
	}

	respondJSON(w, http.StatusOK, summary)
}
