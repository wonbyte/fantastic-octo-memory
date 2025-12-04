package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type BidRepository struct {
	db *Database
}

func NewBidRepository(db *Database) *BidRepository {
	return &BidRepository{db: db}
}

func (r *BidRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Bid, error) {
	query := `
		SELECT id, project_id, job_id, name, total_cost, labor_cost, material_cost, 
		       markup_percentage, final_price, status, bid_data, pdf_url, pdf_s3_key, 
		       created_at, updated_at
		FROM bids
		WHERE id = $1
	`

	var bid models.Bid
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&bid.ID,
		&bid.ProjectID,
		&bid.JobID,
		&bid.Name,
		&bid.TotalCost,
		&bid.LaborCost,
		&bid.MaterialCost,
		&bid.MarkupPercentage,
		&bid.FinalPrice,
		&bid.Status,
		&bid.BidData,
		&bid.PDFURL,
		&bid.PDFS3Key,
		&bid.CreatedAt,
		&bid.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get bid: %w", err)
	}

	return &bid, nil
}

func (r *BidRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Bid, error) {
	query := `
		SELECT id, project_id, job_id, name, total_cost, labor_cost, material_cost, 
		       markup_percentage, final_price, status, bid_data, pdf_url, pdf_s3_key, 
		       created_at, updated_at
		FROM bids
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bids by project: %w", err)
	}
	defer rows.Close()

	var bids []*models.Bid
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(
			&bid.ID,
			&bid.ProjectID,
			&bid.JobID,
			&bid.Name,
			&bid.TotalCost,
			&bid.LaborCost,
			&bid.MaterialCost,
			&bid.MarkupPercentage,
			&bid.FinalPrice,
			&bid.Status,
			&bid.BidData,
			&bid.PDFURL,
			&bid.PDFS3Key,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}

func (r *BidRepository) Create(ctx context.Context, bid *models.Bid) error {
	query := `
		INSERT INTO bids (id, project_id, job_id, name, total_cost, labor_cost, material_cost, 
		                  markup_percentage, final_price, status, bid_data, pdf_url, pdf_s3_key, 
		                  created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		bid.ID,
		bid.ProjectID,
		bid.JobID,
		bid.Name,
		bid.TotalCost,
		bid.LaborCost,
		bid.MaterialCost,
		bid.MarkupPercentage,
		bid.FinalPrice,
		bid.Status,
		bid.BidData,
		bid.PDFURL,
		bid.PDFS3Key,
		bid.CreatedAt,
		bid.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create bid: %w", err)
	}

	return nil
}

func (r *BidRepository) Update(ctx context.Context, bid *models.Bid) error {
	query := `
		UPDATE bids
		SET name = $1, total_cost = $2, labor_cost = $3, material_cost = $4, 
		    markup_percentage = $5, final_price = $6, status = $7, bid_data = $8, 
		    pdf_url = $9, pdf_s3_key = $10, updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.Pool.Exec(ctx, query,
		bid.Name,
		bid.TotalCost,
		bid.LaborCost,
		bid.MaterialCost,
		bid.MarkupPercentage,
		bid.FinalPrice,
		bid.Status,
		bid.BidData,
		bid.PDFURL,
		bid.PDFS3Key,
		bid.UpdatedAt,
		bid.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update bid: %w", err)
	}

	return nil
}
