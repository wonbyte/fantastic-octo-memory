package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type JobRepository struct {
	db *Database
}

func NewJobRepository(db *Database) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	query := `
		SELECT id, blueprint_id, job_type, status, started_at, completed_at, error_message, result_data, created_at, updated_at, retry_count
		FROM jobs
		WHERE id = $1
	`

	var job models.Job
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.BlueprintID,
		&job.JobType,
		&job.Status,
		&job.StartedAt,
		&job.CompletedAt,
		&job.ErrorMessage,
		&job.ResultData,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.RetryCount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return &job, nil
}

func (r *JobRepository) Create(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (id, blueprint_id, job_type, status, started_at, completed_at, error_message, result_data, created_at, updated_at, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		job.ID,
		job.BlueprintID,
		job.JobType,
		job.Status,
		job.StartedAt,
		job.CompletedAt,
		job.ErrorMessage,
		job.ResultData,
		job.CreatedAt,
		job.UpdatedAt,
		job.RetryCount,
	)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

func (r *JobRepository) Update(ctx context.Context, job *models.Job) error {
	query := `
		UPDATE jobs
		SET status = $1, started_at = $2, completed_at = $3, error_message = $4, result_data = $5, updated_at = $6, retry_count = $7
		WHERE id = $8
	`

	_, err := r.db.Pool.Exec(ctx, query,
		job.Status,
		job.StartedAt,
		job.CompletedAt,
		job.ErrorMessage,
		job.ResultData,
		job.UpdatedAt,
		job.RetryCount,
		job.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	return nil
}

func (r *JobRepository) GetQueuedJobs(ctx context.Context, limit int) ([]*models.Job, error) {
	query := `
		SELECT id, blueprint_id, job_type, status, started_at, completed_at, error_message, result_data, created_at, updated_at, retry_count
		FROM jobs
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, models.JobStatusQueued, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get queued jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(
			&job.ID,
			&job.BlueprintID,
			&job.JobType,
			&job.Status,
			&job.StartedAt,
			&job.CompletedAt,
			&job.ErrorMessage,
			&job.ResultData,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.RetryCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}
