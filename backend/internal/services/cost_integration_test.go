package services

import (
	"context"
	"testing"
)

func TestCostIntegrationService_MockProviders(t *testing.T) {
	// Test that mock providers are registered
	service := &CostIntegrationService{
		providers: make(map[string]CostProvider),
	}

	// Register mock providers
	service.RegisterProvider(&MockRSMeansProvider{})
	service.RegisterProvider(&MockHomeDepotProvider{})
	service.RegisterProvider(&MockLowesProvider{})

	// Check providers are registered
	if len(service.providers) != 3 {
		t.Errorf("Expected 3 providers, got %d", len(service.providers))
	}

	providers := []string{"rsmeans", "homedepot", "lowes"}
	for _, name := range providers {
		if _, ok := service.providers[name]; !ok {
			t.Errorf("Provider %s not registered", name)
		}
	}
}

func TestMockRSMeansProvider_GetMaterials(t *testing.T) {
	provider := &MockRSMeansProvider{}
	ctx := context.Background()

	materials, err := provider.GetMaterials(ctx, "national")
	if err != nil {
		t.Fatalf("GetMaterials failed: %v", err)
	}

	if len(materials) == 0 {
		t.Error("Expected materials to be returned")
	}

	// Check that materials have required fields
	for _, m := range materials {
		if m.Name == "" {
			t.Error("Material name is empty")
		}
		if m.Category == "" {
			t.Error("Material category is empty")
		}
		if m.Unit == "" {
			t.Error("Material unit is empty")
		}
		if m.BasePrice <= 0 {
			t.Error("Material base price is not positive")
		}
		if m.Source != "rsmeans" {
			t.Errorf("Expected source 'rsmeans', got '%s'", m.Source)
		}
	}
}

func TestMockRSMeansProvider_GetLaborRates(t *testing.T) {
	provider := &MockRSMeansProvider{}
	ctx := context.Background()

	rates, err := provider.GetLaborRates(ctx, "national")
	if err != nil {
		t.Fatalf("GetLaborRates failed: %v", err)
	}

	if len(rates) == 0 {
		t.Error("Expected labor rates to be returned")
	}

	// Check that rates have required fields
	for _, r := range rates {
		if r.Trade == "" {
			t.Error("Labor rate trade is empty")
		}
		if r.HourlyRate <= 0 {
			t.Error("Labor rate hourly rate is not positive")
		}
		if r.Source != "rsmeans" {
			t.Errorf("Expected source 'rsmeans', got '%s'", r.Source)
		}
	}
}

func TestMockRSMeansProvider_GetRegionalAdjustment(t *testing.T) {
	provider := &MockRSMeansProvider{}
	ctx := context.Background()

	tests := []struct {
		region         string
		expectedFactor float64
	}{
		{"national", 1.00},
		{"california", 1.25},
		{"new_york", 1.30},
		{"texas", 0.95},
		{"unknown_region", 1.00}, // Should default to 1.0
	}

	for _, tt := range tests {
		t.Run(tt.region, func(t *testing.T) {
			adjustment, err := provider.GetRegionalAdjustment(ctx, tt.region)
			if err != nil {
				t.Fatalf("GetRegionalAdjustment failed: %v", err)
			}

			if adjustment.AdjustmentFactor != tt.expectedFactor {
				t.Errorf("Expected factor %f, got %f", tt.expectedFactor, adjustment.AdjustmentFactor)
			}

			if adjustment.Source != "rsmeans" {
				t.Errorf("Expected source 'rsmeans', got '%s'", adjustment.Source)
			}
		})
	}
}

func TestMockHomeDepotProvider_GetMaterials(t *testing.T) {
	provider := &MockHomeDepotProvider{}
	ctx := context.Background()

	materials, err := provider.GetMaterials(ctx, "national")
	if err != nil {
		t.Fatalf("GetMaterials failed: %v", err)
	}

	if len(materials) == 0 {
		t.Error("Expected materials to be returned")
	}

	// Check that all materials have Home Depot source
	for _, m := range materials {
		if m.Source != "homedepot" {
			t.Errorf("Expected source 'homedepot', got '%s'", m.Source)
		}
	}
}

func TestMockLowesProvider_GetMaterials(t *testing.T) {
	provider := &MockLowesProvider{}
	ctx := context.Background()

	materials, err := provider.GetMaterials(ctx, "national")
	if err != nil {
		t.Fatalf("GetMaterials failed: %v", err)
	}

	if len(materials) == 0 {
		t.Error("Expected materials to be returned")
	}

	// Check that all materials have Lowes source
	for _, m := range materials {
		if m.Source != "lowes" {
			t.Errorf("Expected source 'lowes', got '%s'", m.Source)
		}
	}
}

func TestMockHomeDepotProvider_GetLaborRates(t *testing.T) {
	provider := &MockHomeDepotProvider{}
	ctx := context.Background()

	rates, err := provider.GetLaborRates(ctx, "national")
	if err != nil {
		t.Fatalf("GetLaborRates failed: %v", err)
	}

	// Home Depot doesn't provide labor rates, so should return empty
	if len(rates) != 0 {
		t.Errorf("Expected empty labor rates from Home Depot, got %d", len(rates))
	}
}

func TestMockLowesProvider_GetLaborRates(t *testing.T) {
	provider := &MockLowesProvider{}
	ctx := context.Background()

	rates, err := provider.GetLaborRates(ctx, "national")
	if err != nil {
		t.Fatalf("GetLaborRates failed: %v", err)
	}

	// Lowes doesn't provide labor rates, so should return empty
	if len(rates) != 0 {
		t.Errorf("Expected empty labor rates from Lowes, got %d", len(rates))
	}
}
