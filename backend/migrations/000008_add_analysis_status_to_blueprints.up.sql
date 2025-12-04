-- Add analysis_status VARCHAR field to blueprints table
ALTER TABLE blueprints ADD COLUMN analysis_status VARCHAR(50) NOT NULL DEFAULT 'not_started';

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_blueprints_analysis_status ON blueprints(analysis_status);
