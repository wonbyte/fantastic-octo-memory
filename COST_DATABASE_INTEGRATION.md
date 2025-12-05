# Cost Database Integration

This document describes the cost database integration feature that enables real-time pricing from industry sources like RSMeans, Home Depot, and Lowes.

## Overview

The cost database integration provides:
- **Material Pricing**: Real-time material costs from multiple providers
- **Labor Rates**: Industry-standard labor rates by trade and region
- **Regional Adjustments**: Cost-of-living adjustments for different regions
- **Company Overrides**: User/company-specific pricing customization

## Database Schema

### Materials Table
Stores material pricing data from various sources.

```sql
CREATE TABLE materials (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL,
    source VARCHAR(50) NOT NULL,
    source_id VARCHAR(255),
    region VARCHAR(100),
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Labor Rates Table
Stores labor rates by trade and region.

```sql
CREATE TABLE labor_rates (
    id UUID PRIMARY KEY,
    trade VARCHAR(100) NOT NULL,
    description TEXT,
    hourly_rate DECIMAL(10, 2) NOT NULL,
    source VARCHAR(50) NOT NULL,
    source_id VARCHAR(255),
    region VARCHAR(100),
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Regional Adjustments Table
Stores regional cost multipliers.

```sql
CREATE TABLE regional_adjustments (
    id UUID PRIMARY KEY,
    region VARCHAR(100) NOT NULL UNIQUE,
    state_code VARCHAR(2),
    city VARCHAR(100),
    adjustment_factor DECIMAL(5, 4) NOT NULL DEFAULT 1.0000,
    cost_of_living_index INTEGER,
    source VARCHAR(50) NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Company Pricing Overrides Table
Allows companies to override default pricing.

```sql
CREATE TABLE company_pricing_overrides (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    override_type VARCHAR(50) NOT NULL,
    item_key VARCHAR(255) NOT NULL,
    override_value DECIMAL(10, 2) NOT NULL,
    is_percentage BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

## API Endpoints

### Get Materials
```http
GET /api/materials?category=lumber&region=california
```

Returns materials with optional filtering by category and region.

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "Lumber 2x4 8'",
    "description": "Standard 2x4 lumber",
    "category": "lumber",
    "unit": "each",
    "base_price": 7.50,
    "source": "homedepot",
    "region": "california",
    "last_updated": "2024-01-01T00:00:00Z"
  }
]
```

### Get Labor Rates
```http
GET /api/labor-rates?trade=carpentry&region=california
```

Returns labor rates with optional filtering by trade and region.

**Response:**
```json
[
  {
    "id": "uuid",
    "trade": "carpentry",
    "description": "Skilled carpentry work",
    "hourly_rate": 78.00,
    "source": "rsmeans",
    "region": "california",
    "last_updated": "2024-01-01T00:00:00Z"
  }
]
```

### Get Regional Adjustments
```http
GET /api/regional-adjustments
```

Returns all regional cost adjustment factors.

**Response:**
```json
[
  {
    "id": "uuid",
    "region": "california",
    "state_code": "CA",
    "adjustment_factor": 1.2500,
    "source": "rsmeans"
  }
]
```

### Get Company Pricing Overrides
```http
GET /api/company/pricing-overrides
```

Returns all pricing overrides for the authenticated user.

**Response:**
```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "override_type": "material",
    "item_key": "lumber",
    "override_value": 8.00,
    "is_percentage": false,
    "notes": "Custom lumber pricing"
  }
]
```

### Create Pricing Override
```http
POST /api/company/pricing-overrides
Content-Type: application/json

{
  "override_type": "material",
  "item_key": "lumber",
  "override_value": 8.00,
  "is_percentage": false,
  "notes": "Custom lumber pricing"
}
```

### Update Pricing Override
```http
PUT /api/company/pricing-overrides/:id
Content-Type: application/json

{
  "override_value": 8.50,
  "is_percentage": false,
  "notes": "Updated lumber pricing"
}
```

### Delete Pricing Override
```http
DELETE /api/company/pricing-overrides/:id
```

### Sync Cost Data (Admin Only)
```http
POST /api/admin/sync-cost-data
Content-Type: application/json

{
  "provider": "rsmeans",
  "region": "california"
}
```

Syncs cost data from external providers. Available providers:
- `rsmeans` - RSMeans construction cost data
- `homedepot` - Home Depot material pricing
- `lowes` - Lowes material pricing
- `all` - Sync from all providers

## Cost Provider Integration

### Provider Interface
All cost providers implement the `CostProvider` interface:

```go
type CostProvider interface {
    GetMaterials(ctx context.Context, region string) ([]models.MaterialCost, error)
    GetLaborRates(ctx context.Context, region string) ([]models.LaborRate, error)
    GetRegionalAdjustment(ctx context.Context, region string) (*models.RegionalAdjustment, error)
    GetName() string
}
```

### Mock Providers
Currently, the system includes mock implementations for:
- **RSMeans**: Comprehensive construction cost data (materials + labor)
- **Home Depot**: Material pricing only
- **Lowes**: Material pricing only

### Adding Real Provider Implementations
To add a real provider implementation:

1. Create a new provider struct implementing `CostProvider`
2. Implement API calls to the external service
3. Handle authentication (API keys)
4. Transform external data to internal models
5. Register the provider in `NewCostIntegrationService`

Example:
```go
type RealRSMeansProvider struct {
    apiKey    string
    baseURL   string
    httpClient *http.Client
}

func (p *RealRSMeansProvider) GetMaterials(ctx context.Context, region string) ([]models.MaterialCost, error) {
    // Make API call to RSMeans
    // Parse response
    // Return materials
}
```

## Enhanced Pricing Service

The `EnhancedPricingService` uses database-backed pricing with regional adjustments and company overrides:

```go
service := services.NewEnhancedPricingService(
    materialRepo,
    laborRateRepo,
    regionalRepo,
    companyOverrideRepo,
)

// Get pricing config for a user and region
config, err := service.GetPricingConfig(ctx, &userID, &region)

// Generate pricing summary
summary, err := service.GeneratePricingSummary(
    ctx,
    takeoffSummary,
    analysisResult,
    &userID,
    &region,
)
```

The pricing calculation follows this order:
1. Load base prices from database
2. Apply regional adjustment factor
3. Apply company-specific overrides
4. Fall back to default prices if not found

## Regional Adjustments

Regional adjustments are cost-of-living multipliers applied to base prices:

| Region | State | Factor | Description |
|--------|-------|--------|-------------|
| national | - | 1.00 | Baseline |
| california | CA | 1.25 | 25% higher |
| new_york | NY | 1.30 | 30% higher |
| texas | TX | 0.95 | 5% lower |
| florida | FL | 0.98 | 2% lower |

## Company Pricing Overrides

Companies can override default pricing in several ways:

### Override Types
- `material`: Override material price
- `labor`: Override labor rate
- `overhead`: Override overhead percentage
- `profit_margin`: Override profit margin percentage

### Value Types
- **Absolute**: Direct price replacement
- **Percentage**: Percentage adjustment to base price

Example overrides:
```json
{
  "override_type": "material",
  "item_key": "drywall",
  "override_value": 1.75,
  "is_percentage": false
}
```

```json
{
  "override_type": "labor",
  "item_key": "carpentry",
  "override_value": 10.0,
  "is_percentage": true
}
```

## Seeded Data

The migration includes seeded data for immediate use:

### Materials
- Drywall 1/2"
- Lumber 2x4
- Interior Paint
- Vinyl Flooring
- Interior Door
- Standard Window
- Electrical Outlet
- Light Fixture

### Labor Rates
- Carpentry: $75/hr
- Electrical: $95/hr
- Plumbing: $85/hr
- General: $65/hr
- Painting: $55/hr
- Framing: $70/hr

### Regional Adjustments
- National: 1.00x
- California: 1.25x
- New York: 1.30x
- Texas: 0.95x
- Florida: 0.98x
- Illinois: 1.10x

## Migration

To apply the migration:

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

To rollback:

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Caching Implementation

### Redis Caching

The cost database integration includes Redis caching to improve performance and reduce database load. The `CachedCostIntegrationService` wraps the base service with caching capabilities.

#### Cache Configuration

- **Materials**: Cached for 24 hours
- **Labor Rates**: Cached for 24 hours
- **Regional Adjustments**: Cached for 7 days

#### Cache Keys

The system uses structured cache keys for easy invalidation:

```
cost:materials[:category:<name>][:region:<region>]
cost:labor_rates[:trade:<trade>][:region:<region>]
cost:regional_adjustment:region:<region>
```

#### Cache Invalidation

Cache is automatically invalidated when:
- External data is synced via `/api/admin/sync-cost-data`
- Materials are updated in the database
- Labor rates are updated in the database
- Regional adjustments are updated

Manual cache invalidation can be triggered by syncing data through the admin endpoint.

#### Graceful Degradation

If Redis is unavailable:
- The system logs a warning and continues without caching
- All requests fall back to direct database queries
- No functionality is lost

#### Monitoring Cache Performance

Monitor these metrics:
- Cache hit/miss ratio for cost data
- Redis connection status
- Average response times with and without cache
- Cache memory usage

## Testing

Run the tests:

```bash
# Unit tests
go test ./internal/services/... -v

# Repository tests (requires database)
go test ./internal/repository/... -v

# Caching tests
go test ./internal/services/ -v -run "TestCached|TestRedis"
```

## Future Enhancements

1. **Real API Integration**: Replace mock providers with real API implementations
2. ~~**Caching**: Add Redis caching for frequently accessed pricing data~~ ✅ **COMPLETED**
3. **Historical Pricing**: Track price changes over time
4. **Bulk Import**: Support CSV/Excel import for custom pricing
5. **Price Alerts**: Notify users of significant price changes
6. **API Rate Limiting**: Implement rate limiting for external API calls
7. **Audit Logging**: Track all pricing changes and overrides
8. **Price Comparison**: Compare prices across providers
9. **Custom Materials**: Allow users to add custom materials
10. **Export**: Export pricing data to various formats
11. **Admin UI**: Build admin interface for rate table management
12. **Cache Warmup**: Pre-populate cache on application startup
13. **Cache Metrics**: Add detailed cache performance monitoring

## Security Considerations

1. **API Keys**: Store provider API keys securely in environment variables
2. **Authentication**: All endpoints require authentication
3. **Authorization**: Only users can access their own overrides
4. **Admin Routes**: `/api/admin/*` routes should have admin-only middleware
5. **Input Validation**: Validate all user inputs
6. **SQL Injection**: Use parameterized queries (already implemented)
7. **Rate Limiting**: Implement rate limiting on API endpoints

## Monitoring

Monitor these metrics:
- API call latency to external providers
- Cache hit/miss ratio
- Number of pricing overrides per company
- Frequency of price updates
- Database query performance

## Support

For questions or issues, please open a GitHub issue.
