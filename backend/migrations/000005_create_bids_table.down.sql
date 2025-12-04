-- Drop bids table
DROP INDEX IF EXISTS idx_bids_project_id;
DROP INDEX IF EXISTS idx_bids_job_id;
DROP INDEX IF EXISTS idx_bids_status;
DROP TABLE IF EXISTS bids;
