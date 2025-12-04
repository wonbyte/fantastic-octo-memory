-- Drop jobs table
DROP INDEX IF EXISTS idx_jobs_blueprint_id;
DROP INDEX IF EXISTS idx_jobs_status;
DROP INDEX IF EXISTS idx_jobs_created_at;
DROP TABLE IF EXISTS jobs;
