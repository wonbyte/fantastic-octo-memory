package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Name         *string    `json:"name"`
	CompanyName  *string    `json:"company_name"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ProjectStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusArchived  ProjectStatus = "archived"
)

type Project struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.UUID     `json:"user_id"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Status      ProjectStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type UploadStatus string

const (
	UploadStatusPending  UploadStatus = "pending"
	UploadStatusUploaded UploadStatus = "uploaded"
	UploadStatusFailed   UploadStatus = "failed"
)

type Blueprint struct {
	ID           uuid.UUID    `json:"id"`
	ProjectID    uuid.UUID    `json:"project_id"`
	Filename     string       `json:"filename"`
	S3Key        string       `json:"s3_key"`
	FileSize     *int64       `json:"file_size"`
	MimeType     *string      `json:"mime_type"`
	UploadStatus UploadStatus `json:"upload_status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type JobType string

const (
	JobTypeTakeoff       JobType = "takeoff"
	JobTypeEstimate      JobType = "estimate"
	JobTypeBidGeneration JobType = "bid_generation"
)

type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID           uuid.UUID  `json:"id"`
	BlueprintID  uuid.UUID  `json:"blueprint_id"`
	JobType      JobType    `json:"job_type"`
	Status       JobStatus  `json:"status"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	ErrorMessage *string    `json:"error_message"`
	ResultData   *string    `json:"result_data"` // JSONB stored as string
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	RetryCount   int        `json:"retry_count"`
}

type BidStatus string

const (
	BidStatusDraft    BidStatus = "draft"
	BidStatusSent     BidStatus = "sent"
	BidStatusAccepted BidStatus = "accepted"
	BidStatusRejected BidStatus = "rejected"
)

type Bid struct {
	ID               uuid.UUID `json:"id"`
	ProjectID        uuid.UUID `json:"project_id"`
	JobID            *uuid.UUID `json:"job_id"`
	Name             *string    `json:"name"`
	TotalCost        *float64   `json:"total_cost"`
	LaborCost        *float64   `json:"labor_cost"`
	MaterialCost     *float64   `json:"material_cost"`
	MarkupPercentage *float64   `json:"markup_percentage"`
	FinalPrice       *float64   `json:"final_price"`
	Status           BidStatus  `json:"status"`
	BidData          *string    `json:"bid_data"` // JSONB stored as string
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
