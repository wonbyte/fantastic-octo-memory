-- Remove analysis_data field from blueprints table
DROP INDEX IF EXISTS idx_blueprints_analysis_data;
ALTER TABLE blueprints DROP COLUMN IF EXISTS analysis_data;
