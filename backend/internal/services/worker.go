package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/config"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
)

type Worker struct {
	jobRepo       *repository.JobRepository
	blueprintRepo *repository.BlueprintRepository
	aiService     *AIService
	config        *config.WorkerConfig
	stopChan      chan struct{}
	doneChan      chan struct{}
}

func NewWorker(
	jobRepo *repository.JobRepository,
	blueprintRepo *repository.BlueprintRepository,
	aiService *AIService,
	cfg *config.Config,
) *Worker {
	return &Worker{
		jobRepo:       jobRepo,
		blueprintRepo: blueprintRepo,
		aiService:     aiService,
		config:        &cfg.Worker,
		stopChan:      make(chan struct{}),
		doneChan:      make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	slog.Info("Worker started", "poll_interval", w.config.PollInterval)

	ticker := time.NewTicker(w.config.PollInterval)
	defer ticker.Stop()

	go func() {
		defer close(w.doneChan)

		for {
			select {
			case <-ctx.Done():
				slog.Info("Worker stopping due to context cancellation")
				return
			case <-w.stopChan:
				slog.Info("Worker stopping due to stop signal")
				return
			case <-ticker.C:
				w.processJobs(ctx)
			}
		}
	}()
}

func (w *Worker) Stop() {
	slog.Info("Worker stop requested")
	close(w.stopChan)
	<-w.doneChan
	slog.Info("Worker stopped")
}

func (w *Worker) processJobs(ctx context.Context) {
	jobs, err := w.jobRepo.GetQueuedJobs(ctx, 10)
	if err != nil {
		slog.Error("Failed to get queued jobs", "error", err)
		return
	}

	for _, job := range jobs {
		if err := w.processJob(ctx, job); err != nil {
			slog.Error("Failed to process job", "job_id", job.ID, "error", err)
		}
	}
}

func (w *Worker) processJob(ctx context.Context, job *models.Job) error {
	slog.Info("Processing job", "job_id", job.ID, "job_type", job.JobType)

	// Update job to processing
	now := time.Now()
	job.Status = models.JobStatusProcessing
	job.StartedAt = &now
	job.UpdatedAt = now

	if err := w.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	// Get blueprint
	blueprint, err := w.blueprintRepo.GetByID(ctx, job.BlueprintID)
	if err != nil {
		return w.failJob(ctx, job, fmt.Sprintf("failed to get blueprint: %v", err))
	}

	// Call AI service
	resultData, err := w.aiService.AnalyzeBlueprint(ctx, blueprint.ID, blueprint.S3Key)
	if err != nil {
		// Check if we should retry
		if job.RetryCount < w.config.MaxRetries {
			job.RetryCount++
			job.Status = models.JobStatusQueued
			job.StartedAt = nil
			job.UpdatedAt = time.Now()
			
			if updateErr := w.jobRepo.Update(ctx, job); updateErr != nil {
				slog.Error("Failed to requeue job", "job_id", job.ID, "error", updateErr)
			} else {
				slog.Info("Job requeued for retry", "job_id", job.ID, "retry_count", job.RetryCount)
			}
			return err
		}

		return w.failJob(ctx, job, fmt.Sprintf("AI service error: %v", err))
	}

	// Update job to completed
	completedAt := time.Now()
	job.Status = models.JobStatusCompleted
	job.CompletedAt = &completedAt
	job.ResultData = &resultData
	job.UpdatedAt = completedAt

	if err := w.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job to completed: %w", err)
	}

	slog.Info("Job completed successfully", "job_id", job.ID)
	return nil
}

func (w *Worker) failJob(ctx context.Context, job *models.Job, errorMsg string) error {
	completedAt := time.Now()
	job.Status = models.JobStatusFailed
	job.CompletedAt = &completedAt
	job.ErrorMessage = &errorMsg
	job.UpdatedAt = completedAt

	if err := w.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job to failed: %w", err)
	}

	slog.Error("Job failed", "job_id", job.ID, "error", errorMsg)
	return fmt.Errorf("job failed: %s", errorMsg)
}
