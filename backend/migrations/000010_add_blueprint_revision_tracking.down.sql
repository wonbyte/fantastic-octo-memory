-- Drop indexes
DROP INDEX IF EXISTS idx_blueprint_revisions_created_at;
DROP INDEX IF EXISTS idx_blueprint_revisions_version;
DROP INDEX IF EXISTS idx_blueprint_revisions_blueprint_id;
DROP INDEX IF EXISTS idx_blueprints_is_latest;
DROP INDEX IF EXISTS idx_blueprints_version;
DROP INDEX IF EXISTS idx_blueprints_parent_id;

-- Drop blueprint_revisions table
DROP TABLE IF EXISTS blueprint_revisions;

-- Remove version tracking columns from blueprints
ALTER TABLE blueprints
DROP CONSTRAINT IF EXISTS fk_blueprints_parent,
DROP COLUMN IF EXISTS is_latest,
DROP COLUMN IF EXISTS parent_blueprint_id,
DROP COLUMN IF EXISTS version;
