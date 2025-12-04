-- Create bids table
CREATE TABLE IF NOT EXISTS bids (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    job_id UUID,
    name VARCHAR(255),
    total_cost DECIMAL(15, 2),
    labor_cost DECIMAL(15, 2),
    material_cost DECIMAL(15, 2),
    markup_percentage DECIMAL(5, 2),
    final_price DECIMAL(15, 2),
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    bid_data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_bids_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_bids_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE SET NULL
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_bids_project_id ON bids(project_id);
CREATE INDEX IF NOT EXISTS idx_bids_job_id ON bids(job_id);
CREATE INDEX IF NOT EXISTS idx_bids_status ON bids(status);
