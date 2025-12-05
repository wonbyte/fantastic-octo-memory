package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type RegionalAdjustmentRepository struct {
	db *pgxpool.Pool
}

func NewRegionalAdjustmentRepository(db *pgxpool.Pool) *RegionalAdjustmentRepository {
	return &RegionalAdjustmentRepository{db: db}
}

// GetAll returns all regional adjustments
func (r *RegionalAdjustmentRepository) GetAll(ctx context.Context) ([]models.RegionalAdjustment, error) {
	query := `
		SELECT id, region, state_code, city, adjustment_factor, cost_of_living_index, source,
		       last_updated, created_at, updated_at
		FROM regional_adjustments
		ORDER BY region
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var adjustments []models.RegionalAdjustment
	for rows.Next() {
		var ra models.RegionalAdjustment
		err := rows.Scan(&ra.ID, &ra.Region, &ra.StateCode, &ra.City, &ra.AdjustmentFactor,
			&ra.CostOfLivingIndex, &ra.Source, &ra.LastUpdated, &ra.CreatedAt, &ra.UpdatedAt)
		if err != nil {
			return nil, err
		}
		adjustments = append(adjustments, ra)
	}

	return adjustments, rows.Err()
}

// GetByID returns a regional adjustment by ID
func (r *RegionalAdjustmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.RegionalAdjustment, error) {
	query := `
		SELECT id, region, state_code, city, adjustment_factor, cost_of_living_index, source,
		       last_updated, created_at, updated_at
		FROM regional_adjustments
		WHERE id = $1
	`

	var ra models.RegionalAdjustment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ra.ID, &ra.Region, &ra.StateCode, &ra.City, &ra.AdjustmentFactor,
		&ra.CostOfLivingIndex, &ra.Source, &ra.LastUpdated, &ra.CreatedAt, &ra.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ra, nil
}

// GetByRegion returns a regional adjustment by region name
func (r *RegionalAdjustmentRepository) GetByRegion(ctx context.Context, region string) (*models.RegionalAdjustment, error) {
	query := `
		SELECT id, region, state_code, city, adjustment_factor, cost_of_living_index, source,
		       last_updated, created_at, updated_at
		FROM regional_adjustments
		WHERE region = $1
	`

	var ra models.RegionalAdjustment
	err := r.db.QueryRow(ctx, query, region).Scan(
		&ra.ID, &ra.Region, &ra.StateCode, &ra.City, &ra.AdjustmentFactor,
		&ra.CostOfLivingIndex, &ra.Source, &ra.LastUpdated, &ra.CreatedAt, &ra.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ra, nil
}

// Create creates a new regional adjustment
func (r *RegionalAdjustmentRepository) Create(ctx context.Context, adjustment *models.RegionalAdjustment) error {
	query := `
		INSERT INTO regional_adjustments (id, region, state_code, city, adjustment_factor, cost_of_living_index, source, last_updated, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		adjustment.ID, adjustment.Region, adjustment.StateCode, adjustment.City,
		adjustment.AdjustmentFactor, adjustment.CostOfLivingIndex, adjustment.Source,
		adjustment.LastUpdated, adjustment.CreatedAt, adjustment.UpdatedAt,
	)
	return err
}

// Update updates a regional adjustment
func (r *RegionalAdjustmentRepository) Update(ctx context.Context, adjustment *models.RegionalAdjustment) error {
	query := `
		UPDATE regional_adjustments
		SET region = $2, state_code = $3, city = $4, adjustment_factor = $5,
		    cost_of_living_index = $6, source = $7, last_updated = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		adjustment.ID, adjustment.Region, adjustment.StateCode, adjustment.City,
		adjustment.AdjustmentFactor, adjustment.CostOfLivingIndex, adjustment.Source,
		adjustment.LastUpdated, adjustment.UpdatedAt,
	)
	return err
}

// Delete deletes a regional adjustment
func (r *RegionalAdjustmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM regional_adjustments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
