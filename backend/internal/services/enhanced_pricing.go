package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
)

// EnhancedPricingService calculates costs using database-backed pricing with regional adjustments
type EnhancedPricingService struct {
	materialRepo         *repository.MaterialRepository
	laborRateRepo        *repository.LaborRateRepository
	regionalRepo         *repository.RegionalAdjustmentRepository
	companyOverrideRepo  *repository.CompanyPricingOverrideRepository
	defaultConfig        *models.PricingConfig
}

func NewEnhancedPricingService(
	materialRepo *repository.MaterialRepository,
	laborRateRepo *repository.LaborRateRepository,
	regionalRepo *repository.RegionalAdjustmentRepository,
	companyOverrideRepo *repository.CompanyPricingOverrideRepository,
) *EnhancedPricingService {
	return &EnhancedPricingService{
		materialRepo:        materialRepo,
		laborRateRepo:       laborRateRepo,
		regionalRepo:        regionalRepo,
		companyOverrideRepo: companyOverrideRepo,
		defaultConfig: &models.PricingConfig{
			MaterialPrices: map[string]float64{
				"drywall":  1.50,
				"lumber":   3.00,
				"paint":    25.00,
				"flooring": 8.50,
				"door":     450.00,
				"window":   850.00,
				"outlet":   125.00,
				"fixture":  200.00,
			},
			LaborRates: map[string]float64{
				"carpentry":  75.00,
				"electrical": 95.00,
				"plumbing":   85.00,
				"general":    65.00,
				"painting":   55.00,
				"framing":    70.00,
			},
			OverheadRate: 15.0,
			ProfitMargin: 20.0,
		},
	}
}

// GetPricingConfig retrieves pricing configuration with database prices, regional adjustments, and user overrides
func (s *EnhancedPricingService) GetPricingConfig(ctx context.Context, userID *uuid.UUID, region *string) (*models.PricingConfig, error) {
	config := &models.PricingConfig{
		MaterialPrices: make(map[string]float64),
		LaborRates:     make(map[string]float64),
		OverheadRate:   s.defaultConfig.OverheadRate,
		ProfitMargin:   s.defaultConfig.ProfitMargin,
	}

	// Get regional adjustment factor
	regionalFactor := 1.0
	if region != nil && s.regionalRepo != nil {
		adjustment, err := s.regionalRepo.GetByRegion(ctx, *region)
		if err == nil && adjustment != nil {
			regionalFactor = adjustment.AdjustmentFactor
		} else {
			slog.Warn("Regional adjustment not found, using default", "region", *region)
		}
	}

	// Load materials from database
	if s.materialRepo != nil {
		materials, err := s.materialRepo.GetAll(ctx, nil, region)
		if err != nil {
			slog.Error("Failed to load materials from database", "error", err)
			// Fall back to default prices
			config.MaterialPrices = s.defaultConfig.MaterialPrices
		} else {
			// Build material price map with regional adjustment
			for _, m := range materials {
				config.MaterialPrices[m.Category] = m.BasePrice * regionalFactor
			}
		}
	} else {
		// No repository, use defaults
		config.MaterialPrices = s.defaultConfig.MaterialPrices
	}

	// Load labor rates from database
	if s.laborRateRepo != nil {
		laborRates, err := s.laborRateRepo.GetAll(ctx, nil, region)
		if err != nil {
			slog.Error("Failed to load labor rates from database", "error", err)
			// Fall back to default rates
			config.LaborRates = s.defaultConfig.LaborRates
		} else {
			// Build labor rate map with regional adjustment
			for _, lr := range laborRates {
				config.LaborRates[lr.Trade] = lr.HourlyRate * regionalFactor
			}
		}
	} else {
		// No repository, use defaults
		config.LaborRates = s.defaultConfig.LaborRates
	}

	// Apply company-specific overrides if userID is provided
	if userID != nil && s.companyOverrideRepo != nil {
		overrides, err := s.companyOverrideRepo.GetByUserID(ctx, *userID)
		if err != nil {
			slog.Warn("Failed to load company overrides", "user_id", userID, "error", err)
		} else {
			for _, override := range overrides {
				switch override.OverrideType {
				case "material":
					if override.IsPercentage {
						// Apply percentage adjustment
						if basePrice, exists := config.MaterialPrices[override.ItemKey]; exists {
							config.MaterialPrices[override.ItemKey] = basePrice * (1 + override.OverrideValue/100)
						}
					} else {
						// Direct override
						config.MaterialPrices[override.ItemKey] = override.OverrideValue
					}
				case "labor":
					if override.IsPercentage {
						// Apply percentage adjustment
						if baseRate, exists := config.LaborRates[override.ItemKey]; exists {
							config.LaborRates[override.ItemKey] = baseRate * (1 + override.OverrideValue/100)
						}
					} else {
						// Direct override
						config.LaborRates[override.ItemKey] = override.OverrideValue
					}
				case "overhead":
					if override.IsPercentage {
						config.OverheadRate = override.OverrideValue
					}
				case "profit_margin":
					if override.IsPercentage {
						config.ProfitMargin = override.OverrideValue
					}
				}
			}
		}
	}

	// Ensure we have all required prices (fall back to defaults if missing)
	for key, price := range s.defaultConfig.MaterialPrices {
		if _, exists := config.MaterialPrices[key]; !exists {
			config.MaterialPrices[key] = price * regionalFactor
		}
	}
	for key, rate := range s.defaultConfig.LaborRates {
		if _, exists := config.LaborRates[key]; !exists {
			config.LaborRates[key] = rate * regionalFactor
		}
	}

	return config, nil
}

// GeneratePricingSummary calculates costs from takeoff data with database-backed pricing
func (s *EnhancedPricingService) GeneratePricingSummary(
	ctx context.Context,
	takeoffSummary *models.TakeoffSummary,
	analysisResult *models.AnalysisResult,
	userID *uuid.UUID,
	region *string,
) (*models.PricingSummary, error) {
	// Get pricing configuration with database prices, regional adjustments, and user overrides
	config, err := s.GetPricingConfig(ctx, userID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing config: %w", err)
	}

	var lineItems []models.LineItem
	var materialCost, laborCost float64
	costsByTrade := make(map[string]float64)

	// Calculate costs from rooms (framing, drywall, flooring)
	if takeoffSummary != nil && takeoffSummary.TotalArea > 0 {
		// Framing and drywall
		framingItem := models.LineItem{
			Description: "Framing and drywall installation",
			Trade:       "framing",
			Quantity:    takeoffSummary.TotalArea,
			Unit:        "sq ft",
			UnitCost:    5.50,
			Total:       math.Round(takeoffSummary.TotalArea * 5.50 * 100) / 100,
		}
		lineItems = append(lineItems, framingItem)
		materialCost += framingItem.Total * 0.4
		laborCost += framingItem.Total * 0.6
		costsByTrade["framing"] += framingItem.Total

		// Flooring
		flooringItem := models.LineItem{
			Description: "Flooring installation",
			Trade:       "general",
			Quantity:    takeoffSummary.TotalArea,
			Unit:        "sq ft",
			UnitCost:    config.MaterialPrices["flooring"],
			Total:       math.Round(takeoffSummary.TotalArea * config.MaterialPrices["flooring"] * 100) / 100,
		}
		lineItems = append(lineItems, flooringItem)
		materialCost += flooringItem.Total * 0.7
		laborCost += flooringItem.Total * 0.3
		costsByTrade["general"] += flooringItem.Total

		// Paint
		paintItem := models.LineItem{
			Description: "Paint and finishing",
			Trade:       "painting",
			Quantity:    takeoffSummary.TotalArea,
			Unit:        "sq ft",
			UnitCost:    3.50,
			Total:       math.Round(takeoffSummary.TotalArea * 3.50 * 100) / 100,
		}
		lineItems = append(lineItems, paintItem)
		materialCost += paintItem.Total * 0.3
		laborCost += paintItem.Total * 0.7
		costsByTrade["painting"] += paintItem.Total
	}

	// Calculate costs from openings (doors and windows)
	if analysisResult != nil {
		doorCount := 0
		windowCount := 0

		for _, opening := range analysisResult.Openings {
			if opening.OpeningType == "door" {
				doorCount += opening.Count
			} else if opening.OpeningType == "window" {
				windowCount += opening.Count
			}
		}

		if doorCount > 0 {
			doorItem := models.LineItem{
				Description: "Interior door installation",
				Trade:       "carpentry",
				Quantity:    float64(doorCount),
				Unit:        "each",
				UnitCost:    config.MaterialPrices["door"],
				Total:       math.Round(float64(doorCount) * config.MaterialPrices["door"] * 100) / 100,
			}
			lineItems = append(lineItems, doorItem)
			materialCost += doorItem.Total * 0.75
			laborCost += doorItem.Total * 0.25
			costsByTrade["carpentry"] += doorItem.Total
		}

		if windowCount > 0 {
			windowItem := models.LineItem{
				Description: "Window installation",
				Trade:       "carpentry",
				Quantity:    float64(windowCount),
				Unit:        "each",
				UnitCost:    config.MaterialPrices["window"],
				Total:       math.Round(float64(windowCount) * config.MaterialPrices["window"] * 100) / 100,
			}
			lineItems = append(lineItems, windowItem)
			materialCost += windowItem.Total * 0.80
			laborCost += windowItem.Total * 0.20
			costsByTrade["carpentry"] += windowItem.Total
		}

		// Calculate costs from fixtures
		fixtureCount := 0
		for _, fixture := range analysisResult.Fixtures {
			fixtureCount += fixture.Count
		}

		if fixtureCount > 0 {
			fixtureItem := models.LineItem{
				Description: "Electrical fixtures and outlets",
				Trade:       "electrical",
				Quantity:    float64(fixtureCount),
				Unit:        "each",
				UnitCost:    config.MaterialPrices["outlet"],
				Total:       math.Round(float64(fixtureCount) * config.MaterialPrices["outlet"] * 100) / 100,
			}
			lineItems = append(lineItems, fixtureItem)
			materialCost += fixtureItem.Total * 0.60
			laborCost += fixtureItem.Total * 0.40
			costsByTrade["electrical"] += fixtureItem.Total
		}
	}

	// Add labor line items by trade
	for trade, cost := range costsByTrade {
		if cost > 0 {
			rate, ok := config.LaborRates[trade]
			if !ok {
				rate = config.LaborRates["general"]
			}
			hours := math.Round((cost * LaborHoursEstimationFactor) / rate)
			if hours > 0 {
				laborItem := models.LineItem{
					Description: fmt.Sprintf("Labor - %s", trade),
					Trade:       trade,
					Quantity:    hours,
					Unit:        "hours",
					UnitCost:    rate,
					Total:       math.Round(hours * rate * 100) / 100,
				}
				lineItems = append(lineItems, laborItem)
				laborCost += laborItem.Total
			}
		}
	}

	// Round costs
	materialCost = math.Round(materialCost * 100) / 100
	laborCost = math.Round(laborCost * 100) / 100
	subtotal := math.Round((materialCost + laborCost) * 100) / 100

	// Calculate overhead and markup
	overheadAmount := math.Round(subtotal * (config.OverheadRate / 100) * 100) / 100
	markupAmount := math.Round((subtotal + overheadAmount) * (config.ProfitMargin / 100) * 100) / 100
	totalPrice := math.Round((subtotal + overheadAmount + markupAmount) * 100) / 100

	return &models.PricingSummary{
		LineItems:      lineItems,
		LaborCost:      laborCost,
		MaterialCost:   materialCost,
		Subtotal:       subtotal,
		OverheadAmount: overheadAmount,
		MarkupAmount:   markupAmount,
		TotalPrice:     totalPrice,
		CostsByTrade:   costsByTrade,
	}, nil
}

// GetDefaultPricingConfig returns the default pricing configuration (for backward compatibility)
func (s *EnhancedPricingService) GetDefaultPricingConfig() *models.PricingConfig {
	return s.defaultConfig
}

// ParseTakeoffData parses takeoff data from JSON string (for backward compatibility)
func (s *EnhancedPricingService) ParseTakeoffData(jsonData string) (*models.TakeoffSummary, *models.AnalysisResult, error) {
	var analysis models.AnalysisResult
	if err := json.Unmarshal([]byte(jsonData), &analysis); err != nil {
		return nil, nil, fmt.Errorf("failed to parse takeoff data: %w", err)
	}

	// Calculate takeoff summary from analysis
	takeoff := &models.TakeoffSummary{
		OpeningCounts: make(map[string]int),
		FixtureCounts: make(map[string]int),
	}

	for _, room := range analysis.Rooms {
		takeoff.TotalArea += room.Area
		takeoff.RoomCount++
		takeoff.RoomBreakdown = append(takeoff.RoomBreakdown, models.RoomSummary{
			Name:       room.Name,
			RoomType:   room.RoomType,
			Area:       room.Area,
			Dimensions: room.Dimensions,
		})
	}

	for _, opening := range analysis.Openings {
		takeoff.OpeningCounts[opening.OpeningType] += opening.Count
		takeoff.OpeningBreakdown = append(takeoff.OpeningBreakdown, models.OpeningSummary{
			OpeningType: opening.OpeningType,
			Count:       opening.Count,
			Size:        opening.Size,
		})
	}

	for _, fixture := range analysis.Fixtures {
		takeoff.FixtureCounts[fixture.Category] += fixture.Count
		takeoff.FixtureBreakdown = append(takeoff.FixtureBreakdown, models.FixtureSummary{
			FixtureType: fixture.FixtureType,
			Category:    fixture.Category,
			Count:       fixture.Count,
		})
	}

	return takeoff, &analysis, nil
}
