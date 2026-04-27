-- Add rating and review count to listings
ALTER TABLE listings ADD COLUMN rating REAL DEFAULT 0.0;
-- STATEMENT
ALTER TABLE listings ADD COLUMN review_count INTEGER DEFAULT 0;
-- STATEMENT
ALTER TABLE listings ADD COLUMN rating_updated_at DATETIME;
