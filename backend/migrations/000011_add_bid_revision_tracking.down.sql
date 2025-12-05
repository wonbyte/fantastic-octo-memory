-- Drop indexes
DROP INDEX IF EXISTS idx_bid_revisions_created_at;
DROP INDEX IF EXISTS idx_bid_revisions_version;
DROP INDEX IF EXISTS idx_bid_revisions_bid_id;
DROP INDEX IF EXISTS idx_bids_is_latest;
DROP INDEX IF EXISTS idx_bids_version;
DROP INDEX IF EXISTS idx_bids_parent_id;

-- Drop bid_revisions table
DROP TABLE IF EXISTS bid_revisions;

-- Remove version tracking columns from bids
ALTER TABLE bids
DROP CONSTRAINT IF EXISTS fk_bids_parent,
DROP COLUMN IF EXISTS is_latest,
DROP COLUMN IF EXISTS parent_bid_id,
DROP COLUMN IF EXISTS version;
