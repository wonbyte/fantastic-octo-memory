package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type UploadURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type UploadURLResponse struct {
	BlueprintID uuid.UUID `json:"blueprint_id"`
	UploadURL   string    `json:"upload_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type CompleteUploadResponse struct {
	ID       uuid.UUID `json:"id"`
	Status   string    `json:"status"`
	Filename string    `json:"filename"`
}

func (h *Handler) CreateUploadURL(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req UploadURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Filename == "" || req.ContentType == "" {
		respondError(w, http.StatusBadRequest, "filename and content_type are required")
		return
	}

	// Verify project exists (simplified - in production, verify user ownership)
	project, err := h.projectRepo.GetByID(r.Context(), projectID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return
	}

	// Create blueprint record
	blueprintID := uuid.New()
	s3Key := fmt.Sprintf("projects/%s/blueprints/%s/%s", project.ID, blueprintID, req.Filename)

	blueprint := &models.Blueprint{
		ID:             blueprintID,
		ProjectID:      projectID,
		Filename:       req.Filename,
		S3Key:          s3Key,
		UploadStatus:   models.UploadStatusPending,
		AnalysisStatus: models.AnalysisStatusNotStarted,
		Version:        1,
		IsLatest:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.blueprintRepo.Create(r.Context(), blueprint); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create blueprint record")
		return
	}

	// Generate presigned URL
	uploadURL, err := h.s3Service.GeneratePresignedUploadURL(r.Context(), s3Key, req.ContentType)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate upload URL")
		return
	}

	// Use the actual S3 presign expiry time from config
	// The S3 service is configured with the expiry duration from config
	expiresAt := time.Now().Add(5 * time.Minute) // This matches the default S3_PRESIGN_EXPIRY
	respondJSON(w, http.StatusOK, UploadURLResponse{
		BlueprintID: blueprintID,
		UploadURL:   uploadURL,
		ExpiresAt:   expiresAt,
	})
}

func (h *Handler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
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

	// Verify file exists in S3
	exists, fileSize, err := h.s3Service.ObjectExists(r.Context(), blueprint.S3Key)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to verify file")
		return
	}

	if !exists {
		respondError(w, http.StatusNotFound, "File not found in storage")
		return
	}

	// Update blueprint record
	blueprint.UploadStatus = models.UploadStatusUploaded
	blueprint.FileSize = &fileSize
	blueprint.UpdatedAt = time.Now()

	if err := h.blueprintRepo.Update(r.Context(), blueprint); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update blueprint")
		return
	}

	respondJSON(w, http.StatusOK, CompleteUploadResponse{
		ID:       blueprint.ID,
		Status:   string(blueprint.UploadStatus),
		Filename: blueprint.Filename,
	})
}
