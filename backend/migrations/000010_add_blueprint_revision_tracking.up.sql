-- Add version tracking to blueprints table
ALTER TABLE blueprints
ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS parent_blueprint_id UUID,
ADD COLUMN IF NOT EXISTS is_latest BOOLEAN NOT NULL DEFAULT true;

-- Create blueprint_revisions table to store version history and comparison metadata
CREATE TABLE IF NOT EXISTS blueprint_revisions (
    id UUID PRIMARY KEY,
    blueprint_id UUID NOT NULL,
    version INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    s3_key VARCHAR(500) NOT NULL,
    file_size BIGINT,
    mime_type VARCHAR(100),
    analysis_data JSONB,
    changes_summary JSONB, -- Stores detected changes from previous version
    created_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blueprint_revisions_blueprint FOREIGN KEY (blueprint_id) REFERENCES blueprints(id) ON DELETE CASCADE,
    CONSTRAINT fk_blueprint_revisions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT unique_blueprint_version UNIQUE (blueprint_id, version)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_blueprints_parent_id ON blueprints(parent_blueprint_id);
CREATE INDEX IF NOT EXISTS idx_blueprints_version ON blueprints(version);
CREATE INDEX IF NOT EXISTS idx_blueprints_is_latest ON blueprints(is_latest);
CREATE INDEX IF NOT EXISTS idx_blueprint_revisions_blueprint_id ON blueprint_revisions(blueprint_id);
CREATE INDEX IF NOT EXISTS idx_blueprint_revisions_version ON blueprint_revisions(version);
CREATE INDEX IF NOT EXISTS idx_blueprint_revisions_created_at ON blueprint_revisions(created_at DESC);

-- Add foreign key constraint for parent blueprint
ALTER TABLE blueprints
ADD CONSTRAINT fk_blueprints_parent FOREIGN KEY (parent_blueprint_id) REFERENCES blueprints(id) ON DELETE SET NULL;
