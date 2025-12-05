package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type LaborRateRepository struct {
	db *pgxpool.Pool
}

func NewLaborRateRepository(db *pgxpool.Pool) *LaborRateRepository {
	return &LaborRateRepository{db: db}
}

// GetAll returns all labor rates, optionally filtered by trade and region
func (r *LaborRateRepository) GetAll(ctx context.Context, trade, region *string) ([]models.LaborRate, error) {
	query := `
		SELECT id, trade, description, hourly_rate, source, source_id, region,
		       last_updated, created_at, updated_at
		FROM labor_rates
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if trade != nil {
		query += fmt.Sprintf(" AND trade = $%d", argCount)
		args = append(args, *trade)
		argCount++
	}

	if region != nil {
		query += fmt.Sprintf(" AND (region = $%d OR region = 'national' OR region IS NULL)", argCount)
		args = append(args, *region)
	}

	query += " ORDER BY trade"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []models.LaborRate
	for rows.Next() {
		var lr models.LaborRate
		err := rows.Scan(&lr.ID, &lr.Trade, &lr.Description, &lr.HourlyRate, &lr.Source,
			&lr.SourceID, &lr.Region, &lr.LastUpdated, &lr.CreatedAt, &lr.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rates = append(rates, lr)
	}

	return rates, rows.Err()
}

// GetByID returns a labor rate by ID
func (r *LaborRateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.LaborRate, error) {
	query := `
		SELECT id, trade, description, hourly_rate, source, source_id, region,
		       last_updated, created_at, updated_at
		FROM labor_rates
		WHERE id = $1
	`

	var lr models.LaborRate
	err := r.db.QueryRow(ctx, query, id).Scan(
		&lr.ID, &lr.Trade, &lr.Description, &lr.HourlyRate, &lr.Source,
		&lr.SourceID, &lr.Region, &lr.LastUpdated, &lr.CreatedAt, &lr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &lr, nil
}

// GetByTrade returns a labor rate by trade and optional region
func (r *LaborRateRepository) GetByTrade(ctx context.Context, trade string, region *string) (*models.LaborRate, error) {
	query := `
		SELECT id, trade, description, hourly_rate, source, source_id, region,
		       last_updated, created_at, updated_at
		FROM labor_rates
		WHERE trade = $1
	`
	args := []interface{}{trade}

	if region != nil {
		query += " AND (region = $2 OR region = 'national' OR region IS NULL) ORDER BY CASE WHEN region = $2 THEN 1 ELSE 2 END LIMIT 1"
		args = append(args, *region)
	} else {
		query += " AND (region = 'national' OR region IS NULL) LIMIT 1"
	}

	var lr models.LaborRate
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&lr.ID, &lr.Trade, &lr.Description, &lr.HourlyRate, &lr.Source,
		&lr.SourceID, &lr.Region, &lr.LastUpdated, &lr.CreatedAt, &lr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &lr, nil
}

// Create creates a new labor rate
func (r *LaborRateRepository) Create(ctx context.Context, rate *models.LaborRate) error {
	query := `
		INSERT INTO labor_rates (id, trade, description, hourly_rate, source, source_id, region, last_updated, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		rate.ID, rate.Trade, rate.Description, rate.HourlyRate, rate.Source,
		rate.SourceID, rate.Region, rate.LastUpdated, rate.CreatedAt, rate.UpdatedAt,
	)
	return err
}

// Update updates a labor rate
func (r *LaborRateRepository) Update(ctx context.Context, rate *models.LaborRate) error {
	query := `
		UPDATE labor_rates
		SET trade = $2, description = $3, hourly_rate = $4, source = $5,
		    source_id = $6, region = $7, last_updated = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		rate.ID, rate.Trade, rate.Description, rate.HourlyRate, rate.Source,
		rate.SourceID, rate.Region, rate.LastUpdated, rate.UpdatedAt,
	)
	return err
}

// Delete deletes a labor rate
func (r *LaborRateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM labor_rates WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
