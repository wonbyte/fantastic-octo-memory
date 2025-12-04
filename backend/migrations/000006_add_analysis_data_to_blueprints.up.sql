-- Add analysis_data JSONB field to blueprints table
ALTER TABLE blueprints ADD COLUMN analysis_data JSONB;

-- Create GIN index for JSONB queries
CREATE INDEX IF NOT EXISTS idx_blueprints_analysis_data ON blueprints USING GIN (analysis_data);
