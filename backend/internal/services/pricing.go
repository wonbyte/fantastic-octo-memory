package services

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

// PricingService calculates costs and generates pricing summaries
type PricingService struct {
	defaultConfig *models.PricingConfig
}

func NewPricingService() *PricingService {
	return &PricingService{
		defaultConfig: &models.PricingConfig{
			MaterialPrices: map[string]float64{
				"drywall":     1.50,  // per sq ft
				"lumber":      3.00,  // per board foot
				"paint":       25.00, // per gallon
				"flooring":    8.50,  // per sq ft
				"door":        450.00, // per unit
				"window":      850.00, // per unit
				"outlet":      125.00, // per unit
				"fixture":     200.00, // per unit
			},
			LaborRates: map[string]float64{
				"carpentry":   75.00,  // per hour
				"electrical":  95.00,  // per hour
				"plumbing":    85.00,  // per hour
				"general":     65.00,  // per hour
				"painting":    55.00,  // per hour
				"framing":     70.00,  // per hour
			},
			OverheadRate: 15.0, // 15% overhead
			ProfitMargin: 20.0, // 20% profit margin
		},
	}
}

// GeneratePricingSummary calculates costs from takeoff data
func (s *PricingService) GeneratePricingSummary(
	takeoffSummary *models.TakeoffSummary,
	analysisResult *models.AnalysisResult,
	config *models.PricingConfig,
) (*models.PricingSummary, error) {
	if config == nil {
		config = s.defaultConfig
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
		materialCost += framingItem.Total * 0.4 // 40% material
		laborCost += framingItem.Total * 0.6    // 60% labor
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
		materialCost += flooringItem.Total * 0.7 // 70% material
		laborCost += flooringItem.Total * 0.3    // 30% labor
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
		materialCost += paintItem.Total * 0.3 // 30% material
		laborCost += paintItem.Total * 0.7    // 70% labor
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
			materialCost += doorItem.Total * 0.75 // 75% material
			laborCost += doorItem.Total * 0.25    // 25% labor
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
			materialCost += windowItem.Total * 0.80 // 80% material
			laborCost += windowItem.Total * 0.20    // 20% labor
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
			materialCost += fixtureItem.Total * 0.60 // 60% material
			laborCost += fixtureItem.Total * 0.40    // 40% labor
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
			hours := math.Round((cost * 0.5) / rate) // Estimate hours based on cost
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

// GetDefaultPricingConfig returns the default pricing configuration
func (s *PricingService) GetDefaultPricingConfig() *models.PricingConfig {
	return s.defaultConfig
}

// ParseTakeoffData parses takeoff data from JSON string
func (s *PricingService) ParseTakeoffData(jsonData string) (*models.TakeoffSummary, *models.AnalysisResult, error) {
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
