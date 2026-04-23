-- Add latitude and longitude columns to listings table
ALTER TABLE listings ADD COLUMN latitude REAL DEFAULT 0.0;
-- STATEMENT
ALTER TABLE listings ADD COLUMN longitude REAL DEFAULT 0.0;
-- STATEMENT
-- Create an index for spatial queries (bounding box search)
CREATE INDEX idx_listings_coords ON listings(latitude, longitude) WHERE is_active = 1 AND status = 'Approved';
