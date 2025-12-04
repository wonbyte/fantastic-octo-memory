-- Add PDF URL field to bids table
ALTER TABLE bids ADD COLUMN IF NOT EXISTS pdf_url TEXT;
ALTER TABLE bids ADD COLUMN IF NOT EXISTS pdf_s3_key TEXT;

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_bids_pdf_url ON bids(pdf_url) WHERE pdf_url IS NOT NULL;
