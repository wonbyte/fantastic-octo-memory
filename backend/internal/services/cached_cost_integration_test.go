package services

import (
	"context"
	"testing"
)

func TestCachedCostIntegrationService_WithoutCache(t *testing.T) {
	// Test that cached service works even without Redis
	// Using nil cache should fall back to database
	
	service := NewCachedCostIntegrationService(
		nil, // materialRepo
		nil, // laborRateRepo
		nil, // regionalRepo
		nil, // cache (Redis client)
	)
	
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
	
	if service.cache != nil {
		t.Error("Expected nil cache")
	}
	
	// Test that cache key builders work
	category := "lumber"
	region := "california"
	
	key := service.buildMaterialsCacheKey(&category, &region)
	if key != "cost:materials:category:lumber:region:california" {
		t.Errorf("Unexpected cache key: %s", key)
	}
	
	trade := "carpentry"
	key = service.buildLaborRatesCacheKey(&trade, &region)
	if key != "cost:labor_rates:trade:carpentry:region:california" {
		t.Errorf("Unexpected cache key: %s", key)
	}
	
	key = service.buildRegionalAdjustmentCacheKey(region)
	if key != "cost:regional_adjustment:region:california" {
		t.Errorf("Unexpected cache key: %s", key)
	}
}

func TestCachedCostIntegrationService_CacheKeyGeneration(t *testing.T) {
	service := NewCachedCostIntegrationService(nil, nil, nil, nil)
	
	tests := []struct {
		name     string
		category *string
		region   *string
		expected string
	}{
		{
			name:     "no filters",
			category: nil,
			region:   nil,
			expected: "cost:materials",
		},
		{
			name:     "category only",
			category: strPtr("drywall"),
			region:   nil,
			expected: "cost:materials:category:drywall",
		},
		{
			name:     "region only",
			category: nil,
			region:   strPtr("texas"),
			expected: "cost:materials:region:texas",
		},
		{
			name:     "both filters",
			category: strPtr("paint"),
			region:   strPtr("florida"),
			expected: "cost:materials:category:paint:region:florida",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := service.buildMaterialsCacheKey(tt.category, tt.region)
			if key != tt.expected {
				t.Errorf("Expected key %s, got %s", tt.expected, key)
			}
		})
	}
}

func TestCachedCostIntegrationService_InvalidateMethods(t *testing.T) {
	// Test invalidation methods with nil cache (should not panic)
	service := NewCachedCostIntegrationService(nil, nil, nil, nil)
	ctx := context.Background()
	
	// These should all complete without error (graceful degradation)
	service.invalidateMaterialsCache(ctx)
	service.invalidateLaborRatesCache(ctx)
	service.invalidateRegionalAdjustmentCache(ctx, "test")
	
	err := service.InvalidateAllCache(ctx)
	if err != nil {
		t.Errorf("InvalidateAllCache should not fail with nil cache: %v", err)
	}
}
