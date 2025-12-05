-- Add version tracking to bids table
ALTER TABLE bids
ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS parent_bid_id UUID,
ADD COLUMN IF NOT EXISTS is_latest BOOLEAN NOT NULL DEFAULT true;

-- Create bid_revisions table to store version history and comparison metadata
CREATE TABLE IF NOT EXISTS bid_revisions (
    id UUID PRIMARY KEY,
    bid_id UUID NOT NULL,
    version INTEGER NOT NULL,
    name VARCHAR(255),
    total_cost DECIMAL(15, 2),
    labor_cost DECIMAL(15, 2),
    material_cost DECIMAL(15, 2),
    markup_percentage DECIMAL(5, 2),
    final_price DECIMAL(15, 2),
    status VARCHAR(50) NOT NULL,
    bid_data JSONB,
    changes_summary JSONB, -- Stores detected changes from previous version
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_bid_revisions_bid FOREIGN KEY (bid_id) REFERENCES bids(id) ON DELETE CASCADE,
    CONSTRAINT fk_bid_revisions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT unique_bid_version UNIQUE (bid_id, version)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_bids_parent_id ON bids(parent_bid_id);
CREATE INDEX IF NOT EXISTS idx_bids_version ON bids(version);
CREATE INDEX IF NOT EXISTS idx_bids_is_latest ON bids(is_latest);
CREATE INDEX IF NOT EXISTS idx_bid_revisions_bid_id ON bid_revisions(bid_id);
CREATE INDEX IF NOT EXISTS idx_bid_revisions_version ON bid_revisions(version);
CREATE INDEX IF NOT EXISTS idx_bid_revisions_created_at ON bid_revisions(created_at DESC);

-- Add foreign key constraint for parent bid
ALTER TABLE bids
ADD CONSTRAINT fk_bids_parent FOREIGN KEY (parent_bid_id) REFERENCES bids(id) ON DELETE SET NULL;
