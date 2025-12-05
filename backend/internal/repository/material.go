package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type MaterialRepository struct {
	db *pgxpool.Pool
}

func NewMaterialRepository(db *pgxpool.Pool) *MaterialRepository {
	return &MaterialRepository{db: db}
}

// GetAll returns all materials, optionally filtered by category and region
func (r *MaterialRepository) GetAll(ctx context.Context, category, region *string) ([]models.MaterialCost, error) {
	query := `
		SELECT id, name, description, category, unit, base_price, source, source_id, region, 
		       last_updated, created_at, updated_at
		FROM materials
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if category != nil {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, *category)
		argCount++
	}

	if region != nil {
		query += fmt.Sprintf(" AND (region = $%d OR region = 'national' OR region IS NULL)", argCount)
		args = append(args, *region)
	}

	query += " ORDER BY category, name"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []models.MaterialCost
	for rows.Next() {
		var m models.MaterialCost
		err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Category, &m.Unit, &m.BasePrice,
			&m.Source, &m.SourceID, &m.Region, &m.LastUpdated, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}

	return materials, rows.Err()
}

// GetByID returns a material by ID
func (r *MaterialRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.MaterialCost, error) {
	query := `
		SELECT id, name, description, category, unit, base_price, source, source_id, region,
		       last_updated, created_at, updated_at
		FROM materials
		WHERE id = $1
	`

	var m models.MaterialCost
	err := r.db.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.Description, &m.Category, &m.Unit, &m.BasePrice,
		&m.Source, &m.SourceID, &m.Region, &m.LastUpdated, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// GetByName returns a material by name and optional region
func (r *MaterialRepository) GetByName(ctx context.Context, name string, region *string) (*models.MaterialCost, error) {
	query := `
		SELECT id, name, description, category, unit, base_price, source, source_id, region,
		       last_updated, created_at, updated_at
		FROM materials
		WHERE name = $1
	`
	args := []interface{}{name}

	if region != nil {
		query += " AND (region = $2 OR region = 'national' OR region IS NULL) ORDER BY CASE WHEN region = $2 THEN 1 ELSE 2 END LIMIT 1"
		args = append(args, *region)
	} else {
		query += " AND (region = 'national' OR region IS NULL) LIMIT 1"
	}

	var m models.MaterialCost
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&m.ID, &m.Name, &m.Description, &m.Category, &m.Unit, &m.BasePrice,
		&m.Source, &m.SourceID, &m.Region, &m.LastUpdated, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Create creates a new material
func (r *MaterialRepository) Create(ctx context.Context, material *models.MaterialCost) error {
	query := `
		INSERT INTO materials (id, name, description, category, unit, base_price, source, source_id, region, last_updated, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		material.ID, material.Name, material.Description, material.Category, material.Unit,
		material.BasePrice, material.Source, material.SourceID, material.Region,
		material.LastUpdated, material.CreatedAt, material.UpdatedAt,
	)
	return err
}

// Update updates a material
func (r *MaterialRepository) Update(ctx context.Context, material *models.MaterialCost) error {
	query := `
		UPDATE materials
		SET name = $2, description = $3, category = $4, unit = $5, base_price = $6,
		    source = $7, source_id = $8, region = $9, last_updated = $10, updated_at = $11
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		material.ID, material.Name, material.Description, material.Category, material.Unit,
		material.BasePrice, material.Source, material.SourceID, material.Region,
		material.LastUpdated, material.UpdatedAt,
	)
	return err
}

// Delete deletes a material
func (r *MaterialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM materials WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
