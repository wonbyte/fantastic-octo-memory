package services

import (
	"context"
	"testing"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

// TestEnhancedPricingService_DefaultConfiguration tests that the enhanced pricing service
// can be created with a default configuration
func TestEnhancedPricingService_DefaultConfiguration(t *testing.T) {
	// Create service with nil repositories (will use defaults)
	service := NewEnhancedPricingService(nil, nil, nil, nil)
	
	if service == nil {
		t.Fatal("Expected service to be created")
	}
	
	if service.defaultConfig == nil {
		t.Fatal("Expected default config to be set")
	}
	
	// Verify default config has required prices
	requiredMaterials := []string{"drywall", "lumber", "paint", "flooring", "door", "window", "outlet", "fixture"}
	for _, material := range requiredMaterials {
		if _, ok := service.defaultConfig.MaterialPrices[material]; !ok {
			t.Errorf("Missing default price for material: %s", material)
		}
	}
	
	// Verify default config has required labor rates
	requiredTrades := []string{"carpentry", "electrical", "plumbing", "general", "painting", "framing"}
	for _, trade := range requiredTrades {
		if _, ok := service.defaultConfig.LaborRates[trade]; !ok {
			t.Errorf("Missing default labor rate for trade: %s", trade)
		}
	}
	
	// Verify overhead and profit margin are set
	if service.defaultConfig.OverheadRate == 0 {
		t.Error("Overhead rate should be set")
	}
	if service.defaultConfig.ProfitMargin == 0 {
		t.Error("Profit margin should be set")
	}
}

// TestEnhancedPricingService_ParseTakeoffData tests the takeoff data parsing
func TestEnhancedPricingService_ParseTakeoffData(t *testing.T) {
	service := NewEnhancedPricingService(nil, nil, nil, nil)
	
	// Test with valid JSON
	validJSON := `{
		"blueprint_id": "test-id",
		"status": "completed",
		"rooms": [
			{"name": "Living Room", "dimensions": "20x15", "area": 300.0},
			{"name": "Bedroom", "dimensions": "15x12", "area": 180.0}
		],
		"openings": [
			{"opening_type": "door", "count": 3, "size": "36x80"},
			{"opening_type": "window", "count": 5, "size": "36x48"}
		],
		"fixtures": [
			{"fixture_type": "outlet", "category": "electrical", "count": 10},
			{"fixture_type": "switch", "category": "electrical", "count": 5}
		],
		"measurements": [],
		"materials": [],
		"confidence_score": 0.95,
		"processing_time_ms": 1000
	}`
	
	takeoff, analysis, err := service.ParseTakeoffData(validJSON)
	if err != nil {
		t.Fatalf("ParseTakeoffData failed: %v", err)
	}
	
	if takeoff == nil {
		t.Fatal("Expected takeoff to be returned")
	}
	if analysis == nil {
		t.Fatal("Expected analysis to be returned")
	}
	
	// Verify takeoff summary calculations
	expectedArea := 300.0 + 180.0
	if takeoff.TotalArea != expectedArea {
		t.Errorf("Expected total area %f, got %f", expectedArea, takeoff.TotalArea)
	}
	
	if takeoff.RoomCount != 2 {
		t.Errorf("Expected room count 2, got %d", takeoff.RoomCount)
	}
	
	if takeoff.OpeningCounts["door"] != 3 {
		t.Errorf("Expected 3 doors, got %d", takeoff.OpeningCounts["door"])
	}
	
	if takeoff.OpeningCounts["window"] != 5 {
		t.Errorf("Expected 5 windows, got %d", takeoff.OpeningCounts["window"])
	}
	
	if takeoff.FixtureCounts["electrical"] != 15 {
		t.Errorf("Expected 15 electrical fixtures, got %d", takeoff.FixtureCounts["electrical"])
	}
}

// TestEnhancedPricingService_GetDefaultPricingConfig tests the default config getter
func TestEnhancedPricingService_GetDefaultPricingConfig(t *testing.T) {
	service := NewEnhancedPricingService(nil, nil, nil, nil)
	
	config := service.GetDefaultPricingConfig()
	if config == nil {
		t.Fatal("Expected config to be returned")
	}
	
	// Verify it's the same as the internal default config
	if len(config.MaterialPrices) != len(service.defaultConfig.MaterialPrices) {
		t.Error("Config material prices don't match")
	}
	
	if len(config.LaborRates) != len(service.defaultConfig.LaborRates) {
		t.Error("Config labor rates don't match")
	}
}

// TestEnhancedPricingService_GeneratePricingSummary_WithDefaults tests pricing calculation
// with default configuration (no database)
func TestEnhancedPricingService_GeneratePricingSummary_WithDefaults(t *testing.T) {
	service := NewEnhancedPricingService(nil, nil, nil, nil)
	ctx := context.Background()
	
	// Create test data
	takeoff := &models.TakeoffSummary{
		TotalArea:     500.0,
		TotalPerimeter: 100.0,
		RoomCount:     2,
		OpeningCounts: map[string]int{},
		FixtureCounts: map[string]int{},
	}
	
	analysis := &models.AnalysisResult{
		BlueprintID: "test-id",
		Status:      "completed",
		Rooms: []models.Room{
			{Name: "Room 1", Dimensions: "20x15", Area: 300.0},
			{Name: "Room 2", Dimensions: "15x12", Area: 200.0},
		},
		Openings: []models.Opening{
			{OpeningType: "door", Count: 2, Size: "36x80"},
			{OpeningType: "window", Count: 3, Size: "36x48"},
		},
		Fixtures: []models.Fixture{
			{FixtureType: "outlet", Category: "electrical", Count: 10},
		},
		Measurements:     []models.Measurement{},
		Materials:        []models.Material{},
		ConfidenceScore:  0.95,
		ProcessingTimeMs: 1000,
	}
	
	// Generate pricing summary with nil user and region (will use defaults)
	summary, err := service.GeneratePricingSummary(ctx, takeoff, analysis, nil, nil)
	if err != nil {
		t.Fatalf("GeneratePricingSummary failed: %v", err)
	}
	
	if summary == nil {
		t.Fatal("Expected summary to be returned")
	}
	
	// Verify basic calculations
	if len(summary.LineItems) == 0 {
		t.Error("Expected line items to be generated")
	}
	
	if summary.MaterialCost <= 0 {
		t.Error("Material cost should be positive")
	}
	
	if summary.LaborCost <= 0 {
		t.Error("Labor cost should be positive")
	}
	
	if summary.Subtotal <= 0 {
		t.Error("Subtotal should be positive")
	}
	
	if summary.TotalPrice <= summary.Subtotal {
		t.Error("Total price should be greater than subtotal (includes overhead and markup)")
	}
	
	// Verify cost breakdown by trade
	if len(summary.CostsByTrade) == 0 {
		t.Error("Expected costs by trade to be calculated")
	}
	
	// Verify overhead and markup are applied
	if summary.OverheadAmount <= 0 {
		t.Error("Overhead amount should be positive")
	}
	
	if summary.MarkupAmount <= 0 {
		t.Error("Markup amount should be positive")
	}
}
