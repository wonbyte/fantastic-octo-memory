-- Remove PDF URL field from bids table
DROP INDEX IF EXISTS idx_bids_pdf_url;
ALTER TABLE bids DROP COLUMN IF EXISTS pdf_s3_key;
ALTER TABLE bids DROP COLUMN IF EXISTS pdf_url;
