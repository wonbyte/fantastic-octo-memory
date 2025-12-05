-- Materials table - stores material pricing data from various sources
CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL, -- e.g., lumber, drywall, paint, etc.
    unit VARCHAR(50) NOT NULL, -- e.g., sq ft, board foot, gallon, each
    base_price DECIMAL(10, 2) NOT NULL,
    source VARCHAR(50) NOT NULL, -- e.g., rsmeans, homedepot, lowes, custom
    source_id VARCHAR(255), -- ID from external source
    region VARCHAR(100), -- Region this price applies to
    last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_materials_category ON materials(category);
CREATE INDEX idx_materials_source ON materials(source);
CREATE INDEX idx_materials_name ON materials(name);
CREATE INDEX idx_materials_region ON materials(region);

-- Labor rates table - stores labor rates by trade and region
CREATE TABLE IF NOT EXISTS labor_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade VARCHAR(100) NOT NULL, -- e.g., carpentry, electrical, plumbing, etc.
    description TEXT,
    hourly_rate DECIMAL(10, 2) NOT NULL,
    source VARCHAR(50) NOT NULL, -- e.g., rsmeans, custom
    source_id VARCHAR(255), -- ID from external source
    region VARCHAR(100), -- Region this rate applies to
    last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_labor_rates_trade ON labor_rates(trade);
CREATE INDEX idx_labor_rates_source ON labor_rates(source);
CREATE INDEX idx_labor_rates_region ON labor_rates(region);

-- Regional adjustment factors - multipliers for different regions
CREATE TABLE IF NOT EXISTS regional_adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    region VARCHAR(100) NOT NULL UNIQUE,
    state_code VARCHAR(2), -- US state code
    city VARCHAR(100),
    adjustment_factor DECIMAL(5, 4) NOT NULL DEFAULT 1.0000, -- e.g., 1.2 for 20% higher costs
    cost_of_living_index INTEGER, -- Optional COL index
    source VARCHAR(50) NOT NULL, -- e.g., rsmeans, bls, custom
    last_updated TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_regional_adjustments_region ON regional_adjustments(region);
CREATE INDEX idx_regional_adjustments_state ON regional_adjustments(state_code);

-- Company pricing overrides - allow companies to override base pricing
CREATE TABLE IF NOT EXISTS company_pricing_overrides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    override_type VARCHAR(50) NOT NULL, -- e.g., material, labor, overhead, profit_margin
    item_key VARCHAR(255) NOT NULL, -- e.g., material name or trade name
    override_value DECIMAL(10, 2) NOT NULL, -- The overridden value
    is_percentage BOOLEAN DEFAULT FALSE, -- Whether value is a percentage adjustment
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_company_pricing_overrides_user ON company_pricing_overrides(user_id);
CREATE INDEX idx_company_pricing_overrides_type ON company_pricing_overrides(override_type);
CREATE INDEX idx_company_pricing_overrides_item ON company_pricing_overrides(item_key);

-- Seed some default data for immediate use
INSERT INTO materials (name, description, category, unit, base_price, source, region) VALUES
    ('Drywall 1/2"', 'Standard 1/2" drywall sheet', 'drywall', 'sq ft', 1.50, 'custom', 'national'),
    ('Lumber 2x4', 'Standard 2x4 lumber', 'lumber', 'board foot', 3.00, 'custom', 'national'),
    ('Interior Paint', 'Premium interior latex paint', 'paint', 'gallon', 25.00, 'custom', 'national'),
    ('Vinyl Flooring', 'Standard vinyl flooring', 'flooring', 'sq ft', 8.50, 'custom', 'national'),
    ('Interior Door', 'Standard interior door', 'door', 'each', 450.00, 'custom', 'national'),
    ('Standard Window', 'Standard double-hung window', 'window', 'each', 850.00, 'custom', 'national'),
    ('Electrical Outlet', 'Standard electrical outlet', 'outlet', 'each', 125.00, 'custom', 'national'),
    ('Light Fixture', 'Standard light fixture', 'fixture', 'each', 200.00, 'custom', 'national');

INSERT INTO labor_rates (trade, description, hourly_rate, source, region) VALUES
    ('carpentry', 'General carpentry work', 75.00, 'custom', 'national'),
    ('electrical', 'Licensed electrical work', 95.00, 'custom', 'national'),
    ('plumbing', 'Licensed plumbing work', 85.00, 'custom', 'national'),
    ('general', 'General labor', 65.00, 'custom', 'national'),
    ('painting', 'Professional painting', 55.00, 'custom', 'national'),
    ('framing', 'Framing and structural work', 70.00, 'custom', 'national');

INSERT INTO regional_adjustments (region, state_code, adjustment_factor, source) VALUES
    ('national', NULL, 1.0000, 'custom'),
    ('california', 'CA', 1.2500, 'custom'),
    ('new_york', 'NY', 1.3000, 'custom'),
    ('texas', 'TX', 0.9500, 'custom'),
    ('florida', 'FL', 0.9800, 'custom'),
    ('illinois', 'IL', 1.1000, 'custom');
