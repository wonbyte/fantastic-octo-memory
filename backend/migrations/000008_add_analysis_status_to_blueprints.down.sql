-- Remove analysis_status field from blueprints table
DROP INDEX IF EXISTS idx_blueprints_analysis_status;
ALTER TABLE blueprints DROP COLUMN IF EXISTS analysis_status;
