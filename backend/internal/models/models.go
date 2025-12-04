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
	AnalysisData *string      `json:"analysis_data"` // JSONB stored as string
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

// Analysis models - match Python AI service response and TypeScript frontend

type Room struct {
	Name       string  `json:"name"`
	Dimensions string  `json:"dimensions"`
	Area       float64 `json:"area"`
	RoomType   *string `json:"room_type,omitempty"`
}

type Opening struct {
	OpeningType string  `json:"opening_type"`
	Count       int     `json:"count"`
	Size        string  `json:"size"`
	Details     *string `json:"details,omitempty"`
}

type Fixture struct {
	FixtureType string  `json:"fixture_type"`
	Category    string  `json:"category"`
	Count       int     `json:"count"`
	Details     *string `json:"details,omitempty"`
}

type Measurement struct {
	MeasurementType string  `json:"measurement_type"`
	Value           float64 `json:"value"`
	Unit            string  `json:"unit"`
	Location        *string `json:"location,omitempty"`
}

type Material struct {
	MaterialName   string  `json:"material_name"`
	Quantity       float64 `json:"quantity"`
	Unit           string  `json:"unit"`
	Specifications *string `json:"specifications,omitempty"`
}

type AnalysisResult struct {
	BlueprintID      string        `json:"blueprint_id"`
	Status           string        `json:"status"`
	Rooms            []Room        `json:"rooms"`
	Openings         []Opening     `json:"openings"`
	Fixtures         []Fixture     `json:"fixtures"`
	Measurements     []Measurement `json:"measurements"`
	Materials        []Material    `json:"materials"`
	RawOCRText       *string       `json:"raw_ocr_text,omitempty"`
	ConfidenceScore  float64       `json:"confidence_score"`
	ProcessingTimeMs int           `json:"processing_time_ms"`
}

// TakeoffSummary represents aggregated takeoff calculations
type TakeoffSummary struct {
	TotalArea       float64            `json:"total_area"`        // Sum of all room areas (SF)
	TotalPerimeter  float64            `json:"total_perimeter"`   // Sum of all room perimeters (LF)
	OpeningCounts   map[string]int     `json:"opening_counts"`    // Count by opening type (door, window)
	FixtureCounts   map[string]int     `json:"fixture_counts"`    // Count by fixture category
	RoomCount       int                `json:"room_count"`        // Total number of rooms
	RoomBreakdown   []RoomSummary      `json:"room_breakdown"`    // Per-room details
	OpeningBreakdown []OpeningSummary  `json:"opening_breakdown"` // Per-opening details
	FixtureBreakdown []FixtureSummary  `json:"fixture_breakdown"` // Per-fixture details
}

type RoomSummary struct {
	Name       string  `json:"name"`
	RoomType   *string `json:"room_type,omitempty"`
	Area       float64 `json:"area"`
	Dimensions string  `json:"dimensions"`
}

type OpeningSummary struct {
	OpeningType string `json:"opening_type"`
	Count       int    `json:"count"`
	Size        string `json:"size"`
}

type FixtureSummary struct {
	FixtureType string `json:"fixture_type"`
	Category    string `json:"category"`
	Count       int    `json:"count"`
}
