package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type AnalyzeResponse struct {
	JobID  uuid.UUID `json:"job_id"`
	Status string    `json:"status"`
}

type JobStatusResponse struct {
	ID           uuid.UUID  `json:"id"`
	BlueprintID  uuid.UUID  `json:"blueprint_id"`
	JobType      string     `json:"job_type"`
	Status       string     `json:"status"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	ErrorMessage *string    `json:"error_message"`
	ResultData   *string    `json:"result_data"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (h *Handler) AnalyzeBlueprint(w http.ResponseWriter, r *http.Request) {
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

	// Verify blueprint is uploaded
	if blueprint.UploadStatus != models.UploadStatusUploaded {
		respondError(w, http.StatusBadRequest, "Blueprint must be uploaded before analysis")
		return
	}

	// Create job record
	jobID := uuid.New()
	job := &models.Job{
		ID:          jobID,
		BlueprintID: blueprintID,
		JobType:     models.JobTypeTakeoff,
		Status:      models.JobStatusQueued,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		RetryCount:  0,
	}

	if err := h.jobRepo.Create(r.Context(), job); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create job")
		return
	}

	// Update blueprint analysis status to queued
	blueprint.AnalysisStatus = models.AnalysisStatusQueued
	blueprint.UpdatedAt = time.Now()
	if err := h.blueprintRepo.Update(r.Context(), blueprint); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update blueprint status")
		return
	}

	respondJSON(w, http.StatusOK, AnalyzeResponse{
		JobID:  jobID,
		Status: string(models.JobStatusQueued),
	})
}

func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	// Get job record
	job, err := h.jobRepo.GetByID(r.Context(), jobID)
	if err != nil {
		respondError(w, http.StatusNotFound, "Job not found")
		return
	}

	respondJSON(w, http.StatusOK, JobStatusResponse{
		ID:           job.ID,
		BlueprintID:  job.BlueprintID,
		JobType:      string(job.JobType),
		Status:       string(job.Status),
		StartedAt:    job.StartedAt,
		CompletedAt:  job.CompletedAt,
		ErrorMessage: job.ErrorMessage,
		ResultData:   job.ResultData,
		CreatedAt:    job.CreatedAt,
		UpdatedAt:    job.UpdatedAt,
	})
}
