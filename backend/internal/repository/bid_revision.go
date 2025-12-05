package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type BidRevisionRepository struct {
	db *Database
}

func NewBidRevisionRepository(db *Database) *BidRevisionRepository {
	return &BidRevisionRepository{db: db}
}

func (r *BidRevisionRepository) Create(ctx context.Context, revision *models.BidRevision) error {
	query := `
		INSERT INTO bid_revisions (id, bid_id, version, name, total_cost, labor_cost, 
		                          material_cost, markup_percentage, final_price, status, 
		                          bid_data, changes_summary, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		revision.ID,
		revision.BidID,
		revision.Version,
		revision.Name,
		revision.TotalCost,
		revision.LaborCost,
		revision.MaterialCost,
		revision.MarkupPercentage,
		revision.FinalPrice,
		revision.Status,
		revision.BidData,
		revision.ChangesSummary,
		revision.CreatedBy,
		revision.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create bid revision: %w", err)
	}

	return nil
}

func (r *BidRevisionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BidRevision, error) {
	query := `
		SELECT id, bid_id, version, name, total_cost, labor_cost, material_cost, 
		       markup_percentage, final_price, status, bid_data, changes_summary, 
		       created_by, created_at
		FROM bid_revisions
		WHERE id = $1
	`

	var revision models.BidRevision
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&revision.ID,
		&revision.BidID,
		&revision.Version,
		&revision.Name,
		&revision.TotalCost,
		&revision.LaborCost,
		&revision.MaterialCost,
		&revision.MarkupPercentage,
		&revision.FinalPrice,
		&revision.Status,
		&revision.BidData,
		&revision.ChangesSummary,
		&revision.CreatedBy,
		&revision.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get bid revision: %w", err)
	}

	return &revision, nil
}

func (r *BidRevisionRepository) GetByBidID(ctx context.Context, bidID uuid.UUID) ([]*models.BidRevision, error) {
	query := `
		SELECT id, bid_id, version, name, total_cost, labor_cost, material_cost, 
		       markup_percentage, final_price, status, bid_data, changes_summary, 
		       created_by, created_at
		FROM bid_revisions
		WHERE bid_id = $1
		ORDER BY version DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, bidID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid revisions: %w", err)
	}
	defer rows.Close()

	var revisions []*models.BidRevision
	for rows.Next() {
		var revision models.BidRevision
		err := rows.Scan(
			&revision.ID,
			&revision.BidID,
			&revision.Version,
			&revision.Name,
			&revision.TotalCost,
			&revision.LaborCost,
			&revision.MaterialCost,
			&revision.MarkupPercentage,
			&revision.FinalPrice,
			&revision.Status,
			&revision.BidData,
			&revision.ChangesSummary,
			&revision.CreatedBy,
			&revision.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bid revision: %w", err)
		}
		revisions = append(revisions, &revision)
	}

	return revisions, nil
}

func (r *BidRevisionRepository) GetByVersion(ctx context.Context, bidID uuid.UUID, version int) (*models.BidRevision, error) {
	query := `
		SELECT id, bid_id, version, name, total_cost, labor_cost, material_cost, 
		       markup_percentage, final_price, status, bid_data, changes_summary, 
		       created_by, created_at
		FROM bid_revisions
		WHERE bid_id = $1 AND version = $2
	`

	var revision models.BidRevision
	err := r.db.Pool.QueryRow(ctx, query, bidID, version).Scan(
		&revision.ID,
		&revision.BidID,
		&revision.Version,
		&revision.Name,
		&revision.TotalCost,
		&revision.LaborCost,
		&revision.MaterialCost,
		&revision.MarkupPercentage,
		&revision.FinalPrice,
		&revision.Status,
		&revision.BidData,
		&revision.ChangesSummary,
		&revision.CreatedBy,
		&revision.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get bid revision by version: %w", err)
	}

	return &revision, nil
}

func (r *BidRevisionRepository) GetLatestVersion(ctx context.Context, bidID uuid.UUID) (int, error) {
	query := `
		SELECT COALESCE(MAX(version), 0)
		FROM bid_revisions
		WHERE bid_id = $1
	`

	var version int
	err := r.db.Pool.QueryRow(ctx, query, bidID).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest bid version: %w", err)
	}

	return version, nil
}
