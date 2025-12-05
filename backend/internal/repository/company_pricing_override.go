package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type CompanyPricingOverrideRepository struct {
	db *pgxpool.Pool
}

func NewCompanyPricingOverrideRepository(db *pgxpool.Pool) *CompanyPricingOverrideRepository {
	return &CompanyPricingOverrideRepository{db: db}
}

// GetByUserID returns all pricing overrides for a user
func (r *CompanyPricingOverrideRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.CompanyPricingOverride, error) {
	query := `
		SELECT id, user_id, override_type, item_key, override_value, is_percentage, notes,
		       created_at, updated_at
		FROM company_pricing_overrides
		WHERE user_id = $1
		ORDER BY override_type, item_key
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overrides []models.CompanyPricingOverride
	for rows.Next() {
		var cpo models.CompanyPricingOverride
		err := rows.Scan(&cpo.ID, &cpo.UserID, &cpo.OverrideType, &cpo.ItemKey, &cpo.OverrideValue,
			&cpo.IsPercentage, &cpo.Notes, &cpo.CreatedAt, &cpo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		overrides = append(overrides, cpo)
	}

	return overrides, rows.Err()
}

// GetByUserIDAndType returns pricing overrides for a user filtered by type
func (r *CompanyPricingOverrideRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, overrideType string) ([]models.CompanyPricingOverride, error) {
	query := `
		SELECT id, user_id, override_type, item_key, override_value, is_percentage, notes,
		       created_at, updated_at
		FROM company_pricing_overrides
		WHERE user_id = $1 AND override_type = $2
		ORDER BY item_key
	`

	rows, err := r.db.Query(ctx, query, userID, overrideType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overrides []models.CompanyPricingOverride
	for rows.Next() {
		var cpo models.CompanyPricingOverride
		err := rows.Scan(&cpo.ID, &cpo.UserID, &cpo.OverrideType, &cpo.ItemKey, &cpo.OverrideValue,
			&cpo.IsPercentage, &cpo.Notes, &cpo.CreatedAt, &cpo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		overrides = append(overrides, cpo)
	}

	return overrides, rows.Err()
}

// GetByID returns a pricing override by ID
func (r *CompanyPricingOverrideRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CompanyPricingOverride, error) {
	query := `
		SELECT id, user_id, override_type, item_key, override_value, is_percentage, notes,
		       created_at, updated_at
		FROM company_pricing_overrides
		WHERE id = $1
	`

	var cpo models.CompanyPricingOverride
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cpo.ID, &cpo.UserID, &cpo.OverrideType, &cpo.ItemKey, &cpo.OverrideValue,
		&cpo.IsPercentage, &cpo.Notes, &cpo.CreatedAt, &cpo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &cpo, nil
}

// GetByUserIDTypeAndKey returns a specific pricing override
func (r *CompanyPricingOverrideRepository) GetByUserIDTypeAndKey(ctx context.Context, userID uuid.UUID, overrideType, itemKey string) (*models.CompanyPricingOverride, error) {
	query := `
		SELECT id, user_id, override_type, item_key, override_value, is_percentage, notes,
		       created_at, updated_at
		FROM company_pricing_overrides
		WHERE user_id = $1 AND override_type = $2 AND item_key = $3
	`

	var cpo models.CompanyPricingOverride
	err := r.db.QueryRow(ctx, query, userID, overrideType, itemKey).Scan(
		&cpo.ID, &cpo.UserID, &cpo.OverrideType, &cpo.ItemKey, &cpo.OverrideValue,
		&cpo.IsPercentage, &cpo.Notes, &cpo.CreatedAt, &cpo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &cpo, nil
}

// Create creates a new pricing override
func (r *CompanyPricingOverrideRepository) Create(ctx context.Context, override *models.CompanyPricingOverride) error {
	query := `
		INSERT INTO company_pricing_overrides (id, user_id, override_type, item_key, override_value, is_percentage, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		override.ID, override.UserID, override.OverrideType, override.ItemKey,
		override.OverrideValue, override.IsPercentage, override.Notes,
		override.CreatedAt, override.UpdatedAt,
	)
	return err
}

// Update updates a pricing override
func (r *CompanyPricingOverrideRepository) Update(ctx context.Context, override *models.CompanyPricingOverride) error {
	query := `
		UPDATE company_pricing_overrides
		SET override_type = $2, item_key = $3, override_value = $4, is_percentage = $5, notes = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		override.ID, override.OverrideType, override.ItemKey, override.OverrideValue,
		override.IsPercentage, override.Notes, override.UpdatedAt,
	)
	return err
}

// Delete deletes a pricing override
func (r *CompanyPricingOverrideRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM company_pricing_overrides WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
