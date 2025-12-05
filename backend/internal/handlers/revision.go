package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

// GetBlueprintRevisions returns all revisions for a blueprint
func (h *Handler) GetBlueprintRevisions(w http.ResponseWriter, r *http.Request) {
	blueprintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid blueprint ID")
		return
	}

	revisions, err := h.blueprintRevisionRepo.GetByBlueprintID(r.Context(), blueprintID)
	if err != nil {
		slog.Error("Failed to get blueprint revisions", "blueprint_id", blueprintID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get blueprint revisions")
		return
	}

	respondJSON(w, http.StatusOK, revisions)
}

// CompareBlueprintRevisions compares two blueprint versions and returns the differences
func (h *Handler) CompareBlueprintRevisions(w http.ResponseWriter, r *http.Request) {
	blueprintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid blueprint ID")
		return
	}

	fromVersionStr := r.URL.Query().Get("from")
	toVersionStr := r.URL.Query().Get("to")

	if fromVersionStr == "" || toVersionStr == "" {
		respondError(w, http.StatusBadRequest, "from and to version query parameters are required")
		return
	}

	fromVersion, err := strconv.Atoi(fromVersionStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid from version")
		return
	}

	toVersion, err := strconv.Atoi(toVersionStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid to version")
		return
	}

	// Get revisions
	fromRevision, err := h.blueprintRevisionRepo.GetByVersion(r.Context(), blueprintID, fromVersion)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("From version %d not found", fromVersion))
		return
	}

	toRevision, err := h.blueprintRevisionRepo.GetByVersion(r.Context(), blueprintID, toVersion)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("To version %d not found", toVersion))
		return
	}

	// Compare revisions
	comparisonService := services.NewComparisonService()
	comparison, err := comparisonService.CompareBlueprintRevisions(fromRevision, toRevision)
	if err != nil {
		slog.Error("Failed to compare blueprint revisions", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to compare revisions")
		return
	}

	respondJSON(w, http.StatusOK, comparison)
}

// CreateBlueprintRevision creates a new revision snapshot when a blueprint is updated
func (h *Handler) CreateBlueprintRevision(w http.ResponseWriter, r *http.Request) {
	blueprintID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid blueprint ID")
		return
	}

	// Get current blueprint
	blueprint, err := h.blueprintRepo.GetByID(r.Context(), blueprintID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Blueprint not found")
		return
	}

	// Get next version number
	latestVersion, err := h.blueprintRevisionRepo.GetLatestVersion(r.Context(), blueprintID)
	if err != nil {
		slog.Error("Failed to get latest version", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get latest version")
		return
	}

	newVersion := latestVersion + 1

	// Create revision from current blueprint
	revision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  blueprintID,
		Version:      newVersion,
		Filename:     blueprint.Filename,
		S3Key:        blueprint.S3Key,
		FileSize:     blueprint.FileSize,
		MimeType:     blueprint.MimeType,
		AnalysisData: blueprint.AnalysisData,
		CreatedAt:    time.Now(),
	}

	// Get user ID from context if available
	userID := getUserID(r.Context())
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			revision.CreatedBy = &uid
		}
	}

	// Compare with previous version if exists
	if latestVersion > 0 {
		prevRevision, err := h.blueprintRevisionRepo.GetByVersion(r.Context(), blueprintID, latestVersion)
		if err == nil {
			comparisonService := services.NewComparisonService()
			comparison, err := comparisonService.CompareBlueprintRevisions(prevRevision, revision)
			if err == nil {
				// Store changes summary
				summaryJSON, _ := json.Marshal(comparison)
				summaryStr := string(summaryJSON)
				revision.ChangesSummary = &summaryStr
			}
		}
	}

	if err := h.blueprintRevisionRepo.Create(r.Context(), revision); err != nil {
		slog.Error("Failed to create blueprint revision", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to create revision")
		return
	}

	// Update blueprint version
	blueprint.Version = newVersion
	blueprint.UpdatedAt = time.Now()
	if err := h.blueprintRepo.Update(r.Context(), blueprint); err != nil {
		slog.Warn("Failed to update blueprint version", "error", err)
	}

	respondJSON(w, http.StatusCreated, revision)
}

// GetBidRevisions returns all revisions for a bid
func (h *Handler) GetBidRevisions(w http.ResponseWriter, r *http.Request) {
	bidID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid bid ID")
		return
	}

	revisions, err := h.bidRevisionRepo.GetByBidID(r.Context(), bidID)
	if err != nil {
		slog.Error("Failed to get bid revisions", "bid_id", bidID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get bid revisions")
		return
	}

	respondJSON(w, http.StatusOK, revisions)
}

// CompareBidRevisions compares two bid versions and returns the differences
func (h *Handler) CompareBidRevisions(w http.ResponseWriter, r *http.Request) {
	bidID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid bid ID")
		return
	}

	fromVersionStr := r.URL.Query().Get("from")
	toVersionStr := r.URL.Query().Get("to")

	if fromVersionStr == "" || toVersionStr == "" {
		respondError(w, http.StatusBadRequest, "from and to version query parameters are required")
		return
	}

	fromVersion, err := strconv.Atoi(fromVersionStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid from version")
		return
	}

	toVersion, err := strconv.Atoi(toVersionStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid to version")
		return
	}

	// Get revisions
	fromRevision, err := h.bidRevisionRepo.GetByVersion(r.Context(), bidID, fromVersion)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("From version %d not found", fromVersion))
		return
	}

	toRevision, err := h.bidRevisionRepo.GetByVersion(r.Context(), bidID, toVersion)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("To version %d not found", toVersion))
		return
	}

	// Compare revisions
	comparisonService := services.NewComparisonService()
	comparison, err := comparisonService.CompareBidRevisions(fromRevision, toRevision)
	if err != nil {
		slog.Error("Failed to compare bid revisions", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to compare revisions")
		return
	}

	respondJSON(w, http.StatusOK, comparison)
}

// CreateBidRevision creates a new revision snapshot when a bid is updated
func (h *Handler) CreateBidRevision(w http.ResponseWriter, r *http.Request) {
	bidID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid bid ID")
		return
	}

	// Get current bid
	bid, err := h.bidRepo.GetByID(r.Context(), bidID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Bid not found")
		return
	}

	// Get next version number
	latestVersion, err := h.bidRevisionRepo.GetLatestVersion(r.Context(), bidID)
	if err != nil {
		slog.Error("Failed to get latest version", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get latest version")
		return
	}

	newVersion := latestVersion + 1

	// Create revision from current bid
	revision := &models.BidRevision{
		ID:               uuid.New(),
		BidID:            bidID,
		Version:          newVersion,
		Name:             bid.Name,
		TotalCost:        bid.TotalCost,
		LaborCost:        bid.LaborCost,
		MaterialCost:     bid.MaterialCost,
		MarkupPercentage: bid.MarkupPercentage,
		FinalPrice:       bid.FinalPrice,
		Status:           bid.Status,
		BidData:          bid.BidData,
		CreatedAt:        time.Now(),
	}

	// Get user ID from context if available
	userID := getUserID(r.Context())
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			revision.CreatedBy = &uid
		}
	}

	// Compare with previous version if exists
	if latestVersion > 0 {
		prevRevision, err := h.bidRevisionRepo.GetByVersion(r.Context(), bidID, latestVersion)
		if err == nil {
			comparisonService := services.NewComparisonService()
			comparison, err := comparisonService.CompareBidRevisions(prevRevision, revision)
			if err == nil {
				// Store changes summary
				summaryJSON, _ := json.Marshal(comparison)
				summaryStr := string(summaryJSON)
				revision.ChangesSummary = &summaryStr
			}
		}
	}

	if err := h.bidRevisionRepo.Create(r.Context(), revision); err != nil {
		slog.Error("Failed to create bid revision", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to create revision")
		return
	}

	// Update bid version
	bid.Version = newVersion
	bid.UpdatedAt = time.Now()
	if err := h.bidRepo.Update(r.Context(), bid); err != nil {
		slog.Warn("Failed to update bid version", "error", err)
	}

	respondJSON(w, http.StatusCreated, revision)
}
