package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/services"
)

// GenerateBidRequest represents the request to generate a bid
type GenerateBidRequest struct {
	BlueprintID      uuid.UUID  `json:"blueprint_id"`
	MarkupPercentage float64    `json:"markup_percentage"`
	CompanyName      *string    `json:"company_name"`
	BidName          *string    `json:"bid_name"`
}

// GetProjectBids returns all bids for a project
func (h *Handler) GetProjectBids(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	bids, err := h.bidRepo.GetByProjectID(r.Context(), projectID)
	if err != nil {
		slog.Error("Failed to get bids", "project_id", projectID, "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to get bids")
		return
	}

	respondJSON(w, http.StatusOK, bids)
}

// GenerateBid generates a new bid for a project
func (h *Handler) GenerateBid(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req GenerateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate blueprint exists and belongs to project
	blueprint, err := h.blueprintRepo.GetByID(r.Context(), req.BlueprintID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Blueprint not found")
		return
	}

	if blueprint.ProjectID != projectID {
		respondError(w, http.StatusBadRequest, "Blueprint does not belong to this project")
		return
	}

	// Get blueprint analysis data
	if blueprint.AnalysisData == nil {
		respondError(w, http.StatusBadRequest, "Blueprint must be analyzed before generating bid")
		return
	}

	// Parse takeoff data
	pricingService := services.NewPricingService()
	takeoff, analysis, err := pricingService.ParseTakeoffData(*blueprint.AnalysisData)
	if err != nil {
		slog.Error("Failed to parse takeoff data", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to parse takeoff data")
		return
	}

	// Generate pricing summary
	pricingConfig := pricingService.GetDefaultPricingConfig()
	pricingSummary, err := pricingService.GeneratePricingSummary(takeoff, analysis, pricingConfig)
	if err != nil {
		slog.Error("Failed to generate pricing summary", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to generate pricing summary")
		return
	}

	// Prepare AI service request
	companyInfo := map[string]string{
		"name":      "Quality Construction Co.",
		"license":   "CA-123456",
		"insurance": "Fully insured and bonded",
	}
	if req.CompanyName != nil {
		companyInfo["name"] = *req.CompanyName
	}

	markupPercentage := req.MarkupPercentage
	if markupPercentage == 0 {
		markupPercentage = 20.0 // Default 20%
	}

	aiRequest := map[string]interface{}{
		"project_id":        projectID.String(),
		"blueprint_id":      req.BlueprintID.String(),
		"takeoff_data":      analysis,
		"pricing_rules": map[string]interface{}{
			"material_prices": pricingConfig.MaterialPrices,
			"labor_rates":     pricingConfig.LaborRates,
		},
		"company_info":      companyInfo,
		"markup_percentage": markupPercentage,
	}

	// Call AI service to generate bid
	slog.Info("Calling AI service to generate bid", "project_id", projectID)
	bidResponseJSON, err := h.aiService.GenerateBid(r.Context(), aiRequest)
	if err != nil {
		slog.Error("Failed to generate bid with AI service", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to generate bid")
		return
	}

	// Parse AI response
	var aiResponse models.GenerateBidResponse
	if err := json.Unmarshal([]byte(bidResponseJSON), &aiResponse); err != nil {
		slog.Error("Failed to parse AI response", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to parse bid response")
		return
	}

	// Create bid record
	bidID := uuid.New()
	now := time.Now()
	
	bidName := fmt.Sprintf("Bid-%s", time.Now().Format("20060102-150405"))
	if req.BidName != nil {
		bidName = *req.BidName
	}

	bid := &models.Bid{
		ID:               bidID,
		ProjectID:        projectID,
		Name:             &bidName,
		TotalCost:        &pricingSummary.Subtotal,
		LaborCost:        &aiResponse.LaborCost,
		MaterialCost:     &aiResponse.MaterialCost,
		MarkupPercentage: &markupPercentage,
		FinalPrice:       &aiResponse.TotalPrice,
		Status:           models.BidStatusDraft,
		BidData:          &bidResponseJSON,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := h.bidRepo.Create(r.Context(), bid); err != nil {
		slog.Error("Failed to create bid record", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to save bid")
		return
	}

	// Generate PDF
	project, err := h.projectRepo.GetByID(r.Context(), projectID)
	if err != nil {
		slog.Warn("Failed to get project for PDF generation", "error", err)
		project = &models.Project{Name: "Unknown Project"}
	}

	pdfService := services.NewPDFService()
	pdfBytes, err := pdfService.GenerateBidPDF(bid, &aiResponse, project.Name)
	if err != nil {
		slog.Error("Failed to generate PDF", "error", err)
		// Don't fail the request - PDF can be generated later
	} else {
		// Upload PDF to S3
		pdfKey := pdfService.GeneratePDFFilename(projectID, bidID)
		pdfURL, err := h.s3Service.UploadFile(r.Context(), pdfKey, pdfBytes, "application/pdf")
		if err != nil {
			slog.Error("Failed to upload PDF to S3", "error", err)
		} else {
			// Update bid with PDF URL
			bid.PDFURL = &pdfURL
			bid.PDFS3Key = &pdfKey
			bid.UpdatedAt = time.Now()
			if err := h.bidRepo.Update(r.Context(), bid); err != nil {
				slog.Error("Failed to update bid with PDF URL", "error", err)
			}
		}
	}

	slog.Info("Bid generated successfully", "bid_id", bidID, "project_id", projectID)
	respondJSON(w, http.StatusOK, bid)
}

// GetBid returns a specific bid
func (h *Handler) GetBid(w http.ResponseWriter, r *http.Request) {
	bidID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid bid ID")
		return
	}

	bid, err := h.bidRepo.GetByID(r.Context(), bidID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Bid not found")
		return
	}

	respondJSON(w, http.StatusOK, bid)
}

// GetBidPDF returns the PDF URL for a bid or generates it if not exists
func (h *Handler) GetBidPDF(w http.ResponseWriter, r *http.Request) {
	bidID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid bid ID")
		return
	}

	bid, err := h.bidRepo.GetByID(r.Context(), bidID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Bid not found")
		return
	}

	// If PDF already exists, return URL
	if bid.PDFURL != nil && *bid.PDFURL != "" {
		respondJSON(w, http.StatusOK, map[string]string{
			"pdf_url": *bid.PDFURL,
		})
		return
	}

	// Generate PDF if it doesn't exist
	if bid.BidData == nil {
		respondError(w, http.StatusInternalServerError, "Bid data not available")
		return
	}

	// Parse bid data
	pdfService := services.NewPDFService()
	bidResponse, err := pdfService.ParseBidDataFromJSON(*bid.BidData)
	if err != nil {
		slog.Error("Failed to parse bid data", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to parse bid data")
		return
	}

	// Get project name
	project, err := h.projectRepo.GetByID(r.Context(), bid.ProjectID)
	if err != nil {
		slog.Warn("Failed to get project", "error", err)
		project = &models.Project{Name: "Unknown Project"}
	}

	// Generate PDF
	pdfBytes, err := pdfService.GenerateBidPDF(bid, bidResponse, project.Name)
	if err != nil {
		slog.Error("Failed to generate PDF", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	// Upload to S3
	pdfKey := pdfService.GeneratePDFFilename(bid.ProjectID, bidID)
	pdfURL, err := h.s3Service.UploadFile(r.Context(), pdfKey, pdfBytes, "application/pdf")
	if err != nil {
		slog.Error("Failed to upload PDF to S3", "error", err)
		respondError(w, http.StatusInternalServerError, "Failed to upload PDF")
		return
	}

	// Update bid with PDF URL
	bid.PDFURL = &pdfURL
	bid.PDFS3Key = &pdfKey
	bid.UpdatedAt = time.Now()
	if err := h.bidRepo.Update(r.Context(), bid); err != nil {
		slog.Error("Failed to update bid with PDF URL", "error", err)
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"pdf_url": pdfURL,
	})
}

// GetPricingSummary returns the pricing summary for a blueprint
func (h *Handler) GetPricingSummary(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	blueprintIDStr := r.URL.Query().Get("blueprint_id")
	if blueprintIDStr == "" {
		respondError(w, http.StatusBadRequest, "blueprint_id query parameter required")
		return
	}

	blueprintID, err := uuid.Parse(blueprintIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid blueprint ID")
		return
	}

	// Get blueprint
	blueprint, err := h.blueprintRepo.GetByID(r.Context(), blueprintID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Blueprint not found")
		return
	}

	if blueprint.ProjectID != projectID {
		respondError(w, http.StatusBadRequest, "Blueprint does not belong to this project")
		return
	}

	if blueprint.AnalysisData == nil {
		respondError(w, http.StatusBadRequest, "Blueprint must be analyzed first")
		return
	}

	// Parse and generate pricing
	pricingService := services.NewPricingService()
	takeoff, analysis, err := pricingService.ParseTakeoffData(*blueprint.AnalysisData)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to parse takeoff data")
		return
	}

	pricingConfig := pricingService.GetDefaultPricingConfig()
	pricingSummary, err := pricingService.GeneratePricingSummary(takeoff, analysis, pricingConfig)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate pricing summary")
		return
	}

	respondJSON(w, http.StatusOK, pricingSummary)
}
