package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type BlueprintRevisionRepository struct {
	db *Database
}

func NewBlueprintRevisionRepository(db *Database) *BlueprintRevisionRepository {
	return &BlueprintRevisionRepository{db: db}
}

func (r *BlueprintRevisionRepository) Create(ctx context.Context, revision *models.BlueprintRevision) error {
	query := `
		INSERT INTO blueprint_revisions (id, blueprint_id, version, filename, s3_key, 
		                                 file_size, mime_type, analysis_data, changes_summary, 
		                                 created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		revision.ID,
		revision.BlueprintID,
		revision.Version,
		revision.Filename,
		revision.S3Key,
		revision.FileSize,
		revision.MimeType,
		revision.AnalysisData,
		revision.ChangesSummary,
		revision.CreatedBy,
		revision.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create blueprint revision: %w", err)
	}

	return nil
}

func (r *BlueprintRevisionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BlueprintRevision, error) {
	query := `
		SELECT id, blueprint_id, version, filename, s3_key, file_size, mime_type, 
		       analysis_data, changes_summary, created_by, created_at
		FROM blueprint_revisions
		WHERE id = $1
	`

	var revision models.BlueprintRevision
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&revision.ID,
		&revision.BlueprintID,
		&revision.Version,
		&revision.Filename,
		&revision.S3Key,
		&revision.FileSize,
		&revision.MimeType,
		&revision.AnalysisData,
		&revision.ChangesSummary,
		&revision.CreatedBy,
		&revision.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get blueprint revision: %w", err)
	}

	return &revision, nil
}

func (r *BlueprintRevisionRepository) GetByBlueprintID(ctx context.Context, blueprintID uuid.UUID) ([]*models.BlueprintRevision, error) {
	query := `
		SELECT id, blueprint_id, version, filename, s3_key, file_size, mime_type, 
		       analysis_data, changes_summary, created_by, created_at
		FROM blueprint_revisions
		WHERE blueprint_id = $1
		ORDER BY version DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, blueprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blueprint revisions: %w", err)
	}
	defer rows.Close()

	var revisions []*models.BlueprintRevision
	for rows.Next() {
		var revision models.BlueprintRevision
		err := rows.Scan(
			&revision.ID,
			&revision.BlueprintID,
			&revision.Version,
			&revision.Filename,
			&revision.S3Key,
			&revision.FileSize,
			&revision.MimeType,
			&revision.AnalysisData,
			&revision.ChangesSummary,
			&revision.CreatedBy,
			&revision.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blueprint revision: %w", err)
		}
		revisions = append(revisions, &revision)
	}

	return revisions, nil
}

func (r *BlueprintRevisionRepository) GetByVersion(ctx context.Context, blueprintID uuid.UUID, version int) (*models.BlueprintRevision, error) {
	query := `
		SELECT id, blueprint_id, version, filename, s3_key, file_size, mime_type, 
		       analysis_data, changes_summary, created_by, created_at
		FROM blueprint_revisions
		WHERE blueprint_id = $1 AND version = $2
	`

	var revision models.BlueprintRevision
	err := r.db.Pool.QueryRow(ctx, query, blueprintID, version).Scan(
		&revision.ID,
		&revision.BlueprintID,
		&revision.Version,
		&revision.Filename,
		&revision.S3Key,
		&revision.FileSize,
		&revision.MimeType,
		&revision.AnalysisData,
		&revision.ChangesSummary,
		&revision.CreatedBy,
		&revision.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get blueprint revision by version: %w", err)
	}

	return &revision, nil
}

func (r *BlueprintRevisionRepository) GetLatestVersion(ctx context.Context, blueprintID uuid.UUID) (int, error) {
	query := `
		SELECT COALESCE(MAX(version), 0)
		FROM blueprint_revisions
		WHERE blueprint_id = $1
	`

	var version int
	err := r.db.Pool.QueryRow(ctx, query, blueprintID).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest blueprint version: %w", err)
	}

	return version, nil
}
