-- Add Ada-centric sensory signals to listings
ALTER TABLE listings ADD COLUMN heat_level INTEGER DEFAULT 0;
ALTER TABLE listings ADD COLUMN regional_specialty TEXT;
ALTER TABLE listings ADD COLUMN top_dish TEXT;

-- Update FTS table to include new textual signals for searchability
-- Note: SQLite FTS5 doesn't support easy ALTER TABLE for adding columns to existing index directly.
-- However, we can rebuild the FTS or simply leave them out of FTS if they aren't primary search targets.
-- For Ada, regional_specialty and top_dish ARE search targets.
-- We'll add them to the listings_fts in a future migration if needed, or update the trigger.

-- For now, just the schema columns.
