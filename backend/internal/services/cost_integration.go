package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
)

// CostProvider defines the interface for external cost data providers
type CostProvider interface {
	// GetMaterials retrieves material pricing data
	GetMaterials(ctx context.Context, region string) ([]models.MaterialCost, error)
	// GetLaborRates retrieves labor rate data
	GetLaborRates(ctx context.Context, region string) ([]models.LaborRate, error)
	// GetRegionalAdjustment retrieves regional cost adjustment factor
	GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error)
	// GetName returns the provider name
	GetName() string
}

// CostIntegrationService manages integration with external cost data providers
type CostIntegrationService struct {
	materialRepo  *repository.MaterialRepository
	laborRateRepo *repository.LaborRateRepository
	regionalRepo  *repository.RegionalAdjustmentRepository
	providers     map[string]CostProvider
}

func NewCostIntegrationService(
	materialRepo *repository.MaterialRepository,
	laborRateRepo *repository.LaborRateRepository,
	regionalRepo *repository.RegionalAdjustmentRepository,
) *CostIntegrationService {
	service := &CostIntegrationService{
		materialRepo:  materialRepo,
		laborRateRepo: laborRateRepo,
		regionalRepo:  regionalRepo,
		providers:     make(map[string]CostProvider),
	}

	// Register mock providers (replace with real implementations when API keys are available)
	service.RegisterProvider(&MockRSMeansProvider{})
	service.RegisterProvider(&MockHomeDepotProvider{})
	service.RegisterProvider(&MockLowesProvider{})

	return service
}

// RegisterProvider registers a cost data provider
func (s *CostIntegrationService) RegisterProvider(provider CostProvider) {
	s.providers[provider.GetName()] = provider
}

// SyncMaterials syncs material data from a provider to the database
func (s *CostIntegrationService) SyncMaterials(ctx context.Context, providerName, region string) error {
	provider, ok := s.providers[providerName]
	if !ok {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	materials, err := provider.GetMaterials(ctx, region)
	if err != nil {
		return fmt.Errorf("failed to get materials from provider: %w", err)
	}

	for _, material := range materials {
		// Check if material already exists, update or create
		existing, err := s.materialRepo.GetByName(ctx, material.Name, &region)
		if err == nil && existing != nil {
			// Update existing
			existing.BasePrice = material.BasePrice
			existing.Description = material.Description
			existing.LastUpdated = time.Now()
			existing.UpdatedAt = time.Now()
			if err := s.materialRepo.Update(ctx, existing); err != nil {
				return fmt.Errorf("failed to update material %s: %w", material.Name, err)
			}
		} else {
			// Create new
			material.ID = uuid.New()
			material.CreatedAt = time.Now()
			material.UpdatedAt = time.Now()
			material.LastUpdated = time.Now()
			if err := s.materialRepo.Create(ctx, &material); err != nil {
				return fmt.Errorf("failed to create material %s: %w", material.Name, err)
			}
		}
	}

	return nil
}

// SyncLaborRates syncs labor rate data from a provider to the database
func (s *CostIntegrationService) SyncLaborRates(ctx context.Context, providerName, region string) error {
	provider, ok := s.providers[providerName]
	if !ok {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	rates, err := provider.GetLaborRates(ctx, region)
	if err != nil {
		return fmt.Errorf("failed to get labor rates from provider: %w", err)
	}

	for _, rate := range rates {
		// Check if rate already exists, update or create
		existing, err := s.laborRateRepo.GetByTrade(ctx, rate.Trade, &region)
		if err == nil && existing != nil {
			// Update existing
			existing.HourlyRate = rate.HourlyRate
			existing.Description = rate.Description
			existing.LastUpdated = time.Now()
			existing.UpdatedAt = time.Now()
			if err := s.laborRateRepo.Update(ctx, existing); err != nil {
				return fmt.Errorf("failed to update labor rate %s: %w", rate.Trade, err)
			}
		} else {
			// Create new
			rate.ID = uuid.New()
			rate.CreatedAt = time.Now()
			rate.UpdatedAt = time.Now()
			rate.LastUpdated = time.Now()
			if err := s.laborRateRepo.Create(ctx, &rate); err != nil {
				return fmt.Errorf("failed to create labor rate %s: %w", rate.Trade, err)
			}
		}
	}

	return nil
}

// SyncRegionalAdjustment syncs regional adjustment data from a provider to the database
func (s *CostIntegrationService) SyncRegionalAdjustment(ctx context.Context, providerName, region string) error {
	provider, ok := s.providers[providerName]
	if !ok {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	adjustment, err := provider.GetRegionalAdjustment(ctx, region)
	if err != nil {
		return fmt.Errorf("failed to get regional adjustment from provider: %w", err)
	}

	// Check if adjustment already exists, update or create
	existing, err := s.regionalRepo.GetByRegion(ctx, region)
	if err == nil && existing != nil {
		// Update existing
		existing.AdjustmentFactor = adjustment.AdjustmentFactor
		existing.StateCode = adjustment.StateCode
		existing.City = adjustment.City
		existing.CostOfLivingIndex = adjustment.CostOfLivingIndex
		existing.LastUpdated = time.Now()
		existing.UpdatedAt = time.Now()
		if err := s.regionalRepo.Update(ctx, existing); err != nil {
			return fmt.Errorf("failed to update regional adjustment for %s: %w", region, err)
		}
	} else {
		// Create new
		adjustment.ID = uuid.New()
		adjustment.CreatedAt = time.Now()
		adjustment.UpdatedAt = time.Now()
		adjustment.LastUpdated = time.Now()
		if err := s.regionalRepo.Create(ctx, adjustment); err != nil {
			return fmt.Errorf("failed to create regional adjustment for %s: %w", region, err)
		}
	}

	return nil
}

// SyncAll syncs all cost data from all providers
func (s *CostIntegrationService) SyncAll(ctx context.Context, region string) error {
	for name := range s.providers {
		if err := s.SyncMaterials(ctx, name, region); err != nil {
			return err
		}
		if err := s.SyncLaborRates(ctx, name, region); err != nil {
			return err
		}
		if err := s.SyncRegionalAdjustment(ctx, name, region); err != nil {
			return err
		}
	}
	return nil
}

// Mock implementations for cost providers
// These should be replaced with real API implementations when keys are available

type MockRSMeansProvider struct{}

func (p *MockRSMeansProvider) GetName() string {
	return "rsmeans"
}

func (p *MockRSMeansProvider) GetMaterials(ctx context.Context, region string) ([]models.MaterialCost, error) {
	// Mock implementation - in production, this would call RSMeans API
	// RSMeans provides comprehensive construction cost data including materials and labor
	return []models.MaterialCost{
		{
			Name:        "Drywall 1/2\" - RSMeans",
			Description: strPtr("1/2 inch drywall - RSMeans standard"),
			Category:    "drywall",
			Unit:        "sq ft",
			BasePrice:   1.65,
			Source:      "rsmeans",
			SourceID:    strPtr("RSM-DRY-001"),
			Region:      &region,
		},
		{
			Name:        "Lumber 2x4 8' - RSMeans",
			Description: strPtr("2x4 lumber 8 feet - RSMeans standard"),
			Category:    "lumber",
			Unit:        "each",
			BasePrice:   7.50,
			Source:      "rsmeans",
			SourceID:    strPtr("RSM-LUM-001"),
			Region:      &region,
		},
	}, nil
}

func (p *MockRSMeansProvider) GetLaborRates(ctx context.Context, region string) ([]models.LaborRate, error) {
	// Mock implementation - RSMeans provides industry-standard labor rates
	return []models.LaborRate{
		{
			Trade:       "carpentry",
			Description: strPtr("Skilled carpentry - RSMeans standard"),
			HourlyRate:  78.00,
			Source:      "rsmeans",
			SourceID:    strPtr("RSM-LAB-CARP"),
			Region:      &region,
		},
		{
			Trade:       "electrical",
			Description: strPtr("Licensed electrician - RSMeans standard"),
			HourlyRate:  98.00,
			Source:      "rsmeans",
			SourceID:    strPtr("RSM-LAB-ELEC"),
			Region:      &region,
		},
	}, nil
}

func (p *MockRSMeansProvider) GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error) {
	// Mock implementation - RSMeans provides regional cost indices
	adjustments := map[string]float64{
		"california": 1.25,
		"new_york":   1.30,
		"texas":      0.95,
		"florida":    0.98,
		"national":   1.00,
	}

	factor := adjustments[region]
	if factor == 0 {
		factor = 1.00
	}

	return &models.RegionalAdjustment{
		Region:           region,
		AdjustmentFactor: factor,
		Source:           "rsmeans",
	}, nil
}

type MockHomeDepotProvider struct{}

func (p *MockHomeDepotProvider) GetName() string {
	return "homedepot"
}

func (p *MockHomeDepotProvider) GetMaterials(ctx context.Context, region string) ([]models.MaterialCost, error) {
	// Mock implementation - in production, this would call Home Depot API
	return []models.MaterialCost{
		{
			Name:        "Interior Paint Gallon - Home Depot",
			Description: strPtr("Premium interior latex paint - Home Depot"),
			Category:    "paint",
			Unit:        "gallon",
			BasePrice:   28.00,
			Source:      "homedepot",
			SourceID:    strPtr("HD-PAINT-001"),
			Region:      &region,
		},
		{
			Name:        "Vinyl Flooring - Home Depot",
			Description: strPtr("Luxury vinyl plank flooring - Home Depot"),
			Category:    "flooring",
			Unit:        "sq ft",
			BasePrice:   9.25,
			Source:      "homedepot",
			SourceID:    strPtr("HD-FLOOR-001"),
			Region:      &region,
		},
	}, nil
}

func (p *MockHomeDepotProvider) GetLaborRates(ctx context.Context, region string) ([]models.LaborRate, error) {
	// Home Depot doesn't typically provide labor rates, return empty
	return []models.LaborRate{}, nil
}

func (p *MockHomeDepotProvider) GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error) {
	// Home Depot pricing is already regional, so adjustment factor is 1.0
	return &models.RegionalAdjustment{
		Region:           region,
		AdjustmentFactor: 1.00,
		Source:           "homedepot",
	}, nil
}

type MockLowesProvider struct{}

func (p *MockLowesProvider) GetName() string {
	return "lowes"
}

func (p *MockLowesProvider) GetMaterials(ctx context.Context, region string) ([]models.MaterialCost, error) {
	// Mock implementation - in production, this would call Lowes API
	return []models.MaterialCost{
		{
			Name:        "Interior Door - Lowes",
			Description: strPtr("6-panel interior door - Lowes"),
			Category:    "door",
			Unit:        "each",
			BasePrice:   475.00,
			Source:      "lowes",
			SourceID:    strPtr("LOW-DOOR-001"),
			Region:      &region,
		},
		{
			Name:        "Standard Window - Lowes",
			Description: strPtr("Double-hung vinyl window - Lowes"),
			Category:    "window",
			Unit:        "each",
			BasePrice:   895.00,
			Source:      "lowes",
			SourceID:    strPtr("LOW-WIN-001"),
			Region:      &region,
		},
	}, nil
}

func (p *MockLowesProvider) GetLaborRates(ctx context.Context, region string) ([]models.LaborRate, error) {
	// Lowes doesn't typically provide labor rates, return empty
	return []models.LaborRate{}, nil
}

func (p *MockLowesProvider) GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error) {
	// Lowes pricing is already regional, so adjustment factor is 1.0
	return &models.RegionalAdjustment{
		Region:           region,
		AdjustmentFactor: 1.00,
		Source:           "lowes",
	}, nil
}

func strPtr(s string) *string {
	return &s
}
