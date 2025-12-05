package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/repository"
)

// CachedCostIntegrationService wraps CostIntegrationService with Redis caching
type CachedCostIntegrationService struct {
	*CostIntegrationService
	cache *RedisClient
	// Cache TTL settings
	materialsCacheTTL     time.Duration
	laborRatesCacheTTL    time.Duration
	regionalAdjustmentTTL time.Duration
}

// NewCachedCostIntegrationService creates a new cached cost integration service
func NewCachedCostIntegrationService(
	materialRepo *repository.MaterialRepository,
	laborRateRepo *repository.LaborRateRepository,
	regionalRepo *repository.RegionalAdjustmentRepository,
	cache *RedisClient,
) *CachedCostIntegrationService {
	baseService := NewCostIntegrationService(materialRepo, laborRateRepo, regionalRepo)
	
	return &CachedCostIntegrationService{
		CostIntegrationService: baseService,
		cache:                  cache,
		materialsCacheTTL:      24 * time.Hour, // Materials cached for 24 hours
		laborRatesCacheTTL:     24 * time.Hour, // Labor rates cached for 24 hours
		regionalAdjustmentTTL:  7 * 24 * time.Hour, // Regional adjustments cached for 7 days
	}
}

// GetMaterials retrieves materials with caching
func (s *CachedCostIntegrationService) GetMaterials(ctx context.Context, category, region *string) ([]models.MaterialCost, error) {
	// Build cache key
	cacheKey := s.buildMaterialsCacheKey(category, region)
	
	// Try to get from cache if available
	if s.cache != nil && s.cache.IsAvailable() {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var materials []models.MaterialCost
			if err := json.Unmarshal([]byte(cached), &materials); err == nil {
				slog.Debug("Materials cache hit", "key", cacheKey)
				return materials, nil
			}
		}
	}
	
	// Cache miss - get from database
	materials, err := s.materialRepo.GetAll(ctx, category, region)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	if s.cache != nil && s.cache.IsAvailable() {
		if data, err := json.Marshal(materials); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, s.materialsCacheTTL); err != nil {
				slog.Warn("Failed to cache materials", "error", err)
			}
		} else {
			slog.Warn("Failed to marshal materials for caching", "error", err)
		}
	}
	
	return materials, nil
}

// GetLaborRates retrieves labor rates with caching
func (s *CachedCostIntegrationService) GetLaborRates(ctx context.Context, trade, region *string) ([]models.LaborRate, error) {
	// Build cache key
	cacheKey := s.buildLaborRatesCacheKey(trade, region)
	
	// Try to get from cache if available
	if s.cache != nil && s.cache.IsAvailable() {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var rates []models.LaborRate
			if err := json.Unmarshal([]byte(cached), &rates); err == nil {
				slog.Debug("Labor rates cache hit", "key", cacheKey)
				return rates, nil
			}
		}
	}
	
	// Cache miss - get from database
	rates, err := s.laborRateRepo.GetAll(ctx, trade, region)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	if s.cache != nil && s.cache.IsAvailable() {
		if data, err := json.Marshal(rates); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, s.laborRatesCacheTTL); err != nil {
				slog.Warn("Failed to cache labor rates", "error", err)
			}
		} else {
			slog.Warn("Failed to marshal labor rates for caching", "error", err)
		}
	}
	
	return rates, nil
}

// GetRegionalAdjustment retrieves regional adjustment with caching
func (s *CachedCostIntegrationService) GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error) {
	// Build cache key
	cacheKey := s.buildRegionalAdjustmentCacheKey(region)
	
	// Try to get from cache if available
	if s.cache != nil && s.cache.IsAvailable() {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var adjustment models.RegionalAdjustment
			if err := json.Unmarshal([]byte(cached), &adjustment); err == nil {
				slog.Debug("Regional adjustment cache hit", "key", cacheKey)
				return &adjustment, nil
			}
		}
	}
	
	// Cache miss - get from database
	adjustment, err := s.regionalRepo.GetByRegion(ctx, region)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	if s.cache != nil && s.cache.IsAvailable() {
		if data, err := json.Marshal(adjustment); err == nil {
			if err := s.cache.Set(ctx, cacheKey, data, s.regionalAdjustmentTTL); err != nil {
				slog.Warn("Failed to cache regional adjustment", "error", err)
			}
		} else {
			slog.Warn("Failed to marshal regional adjustment for caching", "error", err)
		}
	}
	
	return adjustment, nil
}

// SyncMaterials syncs materials and invalidates cache
func (s *CachedCostIntegrationService) SyncMaterials(ctx context.Context, providerName, region string) error {
	// Call base implementation
	if err := s.CostIntegrationService.SyncMaterials(ctx, providerName, region); err != nil {
		return err
	}
	
	// Invalidate materials cache
	s.invalidateMaterialsCache(ctx)
	
	return nil
}

// SyncLaborRates syncs labor rates and invalidates cache
func (s *CachedCostIntegrationService) SyncLaborRates(ctx context.Context, providerName, region string) error {
	// Call base implementation
	if err := s.CostIntegrationService.SyncLaborRates(ctx, providerName, region); err != nil {
		return err
	}
	
	// Invalidate labor rates cache
	s.invalidateLaborRatesCache(ctx)
	
	return nil
}

// SyncRegionalAdjustment syncs regional adjustment and invalidates cache
func (s *CachedCostIntegrationService) SyncRegionalAdjustment(ctx context.Context, providerName, region string) error {
	// Call base implementation
	if err := s.CostIntegrationService.SyncRegionalAdjustment(ctx, providerName, region); err != nil {
		return err
	}
	
	// Invalidate regional adjustment cache
	s.invalidateRegionalAdjustmentCache(ctx, region)
	
	return nil
}

// Cache key builders
func (s *CachedCostIntegrationService) buildMaterialsCacheKey(category, region *string) string {
	key := "cost:materials"
	if category != nil {
		key += fmt.Sprintf(":category:%s", *category)
	}
	if region != nil {
		key += fmt.Sprintf(":region:%s", *region)
	}
	return key
}

func (s *CachedCostIntegrationService) buildLaborRatesCacheKey(trade, region *string) string {
	key := "cost:labor_rates"
	if trade != nil {
		key += fmt.Sprintf(":trade:%s", *trade)
	}
	if region != nil {
		key += fmt.Sprintf(":region:%s", *region)
	}
	return key
}

func (s *CachedCostIntegrationService) buildRegionalAdjustmentCacheKey(region string) string {
	return fmt.Sprintf("cost:regional_adjustment:region:%s", region)
}

// Cache invalidation methods
func (s *CachedCostIntegrationService) invalidateMaterialsCache(ctx context.Context) {
	if s.cache != nil && s.cache.IsAvailable() {
		if err := s.cache.DeletePattern(ctx, "cost:materials*"); err != nil {
			slog.Warn("Failed to invalidate materials cache", "error", err)
		} else {
			slog.Info("Materials cache invalidated")
		}
	}
}

func (s *CachedCostIntegrationService) invalidateLaborRatesCache(ctx context.Context) {
	if s.cache != nil && s.cache.IsAvailable() {
		if err := s.cache.DeletePattern(ctx, "cost:labor_rates*"); err != nil {
			slog.Warn("Failed to invalidate labor rates cache", "error", err)
		} else {
			slog.Info("Labor rates cache invalidated")
		}
	}
}

func (s *CachedCostIntegrationService) invalidateRegionalAdjustmentCache(ctx context.Context, region string) {
	if s.cache != nil && s.cache.IsAvailable() {
		cacheKey := s.buildRegionalAdjustmentCacheKey(region)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			slog.Warn("Failed to invalidate regional adjustment cache", "error", err, "region", region)
		} else {
			slog.Info("Regional adjustment cache invalidated", "region", region)
		}
	}
}

// InvalidateAllCache clears all cost-related caches
func (s *CachedCostIntegrationService) InvalidateAllCache(ctx context.Context) error {
	if s.cache == nil || !s.cache.IsAvailable() {
		return nil
	}
	
	if err := s.cache.DeletePattern(ctx, "cost:*"); err != nil {
		return fmt.Errorf("failed to invalidate all cost caches: %w", err)
	}
	
	slog.Info("All cost caches invalidated")
	return nil
}
