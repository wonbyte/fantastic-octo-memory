package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type BlueprintRepository struct {
	db *Database
}

func NewBlueprintRepository(db *Database) *BlueprintRepository {
	return &BlueprintRepository{db: db}
}

func (r *BlueprintRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Blueprint, error) {
	query := `
		SELECT id, project_id, filename, s3_key, file_size, mime_type, upload_status, created_at, updated_at
		FROM blueprints
		WHERE id = $1
	`

	var blueprint models.Blueprint
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&blueprint.ID,
		&blueprint.ProjectID,
		&blueprint.Filename,
		&blueprint.S3Key,
		&blueprint.FileSize,
		&blueprint.MimeType,
		&blueprint.UploadStatus,
		&blueprint.CreatedAt,
		&blueprint.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get blueprint: %w", err)
	}

	return &blueprint, nil
}

func (r *BlueprintRepository) Create(ctx context.Context, blueprint *models.Blueprint) error {
	query := `
		INSERT INTO blueprints (id, project_id, filename, s3_key, file_size, mime_type, upload_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		blueprint.ID,
		blueprint.ProjectID,
		blueprint.Filename,
		blueprint.S3Key,
		blueprint.FileSize,
		blueprint.MimeType,
		blueprint.UploadStatus,
		blueprint.CreatedAt,
		blueprint.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create blueprint: %w", err)
	}

	return nil
}

func (r *BlueprintRepository) Update(ctx context.Context, blueprint *models.Blueprint) error {
	query := `
		UPDATE blueprints
		SET file_size = $1, upload_status = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Pool.Exec(ctx, query,
		blueprint.FileSize,
		blueprint.UploadStatus,
		blueprint.UpdatedAt,
		blueprint.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update blueprint: %w", err)
	}

	return nil
}
